package main

import (
	"flag"
	"go-downloader/downloader"
	"log"
)

var (
	url        string
	path       string
	sliceCount int
)

func init() {
	flag.StringVar(&url, "url", "", "待下载文件的URL")
	flag.StringVar(&path, "o", "", "文件保存位置")
	flag.IntVar(&sliceCount, "n", 8, "下载切片数量")
}

func main() {
	flag.Parse()
	if url == "" {
		log.Fatal("未传入url")
	}

	d, err := downloader.NewDownloader(url, path, sliceCount)
	if err != nil {
		log.Fatal(err)
	}

	err = d.Download()
	if err != nil {
		log.Fatal(err)
	}
}
