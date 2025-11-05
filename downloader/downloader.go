package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gosuri/uiprogress"
)

type Downloader struct {
	url        string
	path       string
	sliceCount int
}

func NewDownloader(url string, path string, threads int) (*Downloader, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	if path == "" {
		s := strings.Split(url, "/")
		path = s[len(s)-1]
	}

	return &Downloader{url, path, threads}, nil
}

func (d *Downloader) dir() string {
	return filepath.Dir(d.path)
}

func (d *Downloader) name() string {
	return filepath.Base(d.path)
}

func (d *Downloader) slicePath(idx int) string {
	return filepath.Join(
		d.dir(),
		fmt.Sprintf("%s-%d.tmp", d.name(), idx),
	)
}

func (d *Downloader) downloadSlice(idx, start, end int, update func(int)) error {
	// 判断是否已经下载完成或部分下载完成
	downloaded := 0
	continueDownload := false
	path := d.slicePath(idx)
	if stat, err := os.Stat(path); err == nil {
		downloaded += int(stat.Size())
		update(downloaded)
		if start > end {
			return nil
		}
		continueDownload = true
	}

	// 建立连接
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start+downloaded, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// 创建或打开文件
	var temp *os.File
	if continueDownload {
		temp, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	} else {
		temp, err = os.Create(path)
	}

	if err != nil {
		return err
	}
	defer temp.Close()

	// 下载

	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		downloaded += n
		update(downloaded)
		_, err = temp.Write(buf[:n])
		if err != nil {
			return err
		}
	}
}

func (d *Downloader) Download() error {
	resp, err := http.Head(d.url)
	if err != nil {
		return err
	}

	// 获取总长度
	sizeStr := resp.Header.Get("Content-Length")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return err
	}

	// 分片
	var ranges [][2]int
	if resp.Header.Get("Accept-Ranges") == "bytes" {
		sliceSize := size / d.sliceCount
		for i := 0; i < size; i += sliceSize {
			ranges = append(ranges, [2]int{i, min(size, i+sliceSize) - 1})
		}
	} else {
		ranges = append(ranges, [2]int{0, size})
	}

	uiprogress.Start()

	// 下载分片
	var wg sync.WaitGroup
	errors := make([]error, len(ranges))
	progresses := make([]int, len(ranges))
	for i, r := range ranges {
		b := NewProgressBar(r[1] - r[0] + 1)
		wg.Go(func() {
			update := func(p int) {
				progresses[i] = p
				_ = b.Set(p)
			}
			errors[i] = d.downloadSlice(i, r[0], r[1], update)
		})
	}

	wg.Wait()

	// 合并
	file, err := os.Create(d.path)
	if err != nil {
		return err
	}

	for i := range ranges {
		sliceFile, err := os.Open(d.slicePath(i))

		if err != nil {
			return err
		}

		_, err = io.Copy(file, sliceFile)
		if err != nil {
			return err
		}

		sliceFile.Close()
	}

	// 删除临时文件
	for i := range ranges {
		err := os.Remove(d.slicePath(i))
		if err != nil {
			return err
		}
	}

	return nil
}
