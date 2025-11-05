package downloader

import (
	"fmt"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

type ProgressBar struct {
	startTime      time.Time
	n              int
	total          int
	avgSpeed       float64
	speed          float64
	lastTime       time.Time
	speedInterval  time.Duration
	intervalStartN int
	bar            *uiprogress.Bar
}

func NewProgressBar(total int) *ProgressBar {
	pb := new(ProgressBar)
	pb.speedInterval = 2 * time.Second
	pb.total = total
	pb.startTime = time.Now()
	pb.bar = uiprogress.AddBar(total)

	// 添加时间
	pb.bar.PrependFunc(func(b *uiprogress.Bar) string {
		t := time.Now().Sub(pb.startTime)
		return strutil.PadLeft(strutil.PrettyTime(t), 5, ' ')
	})

	// 添加下载进度
	pb.bar.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf(
			"%s/%s %s",
			pb.formatSize(float64(pb.n)),
			pb.formatSize(float64(pb.total)),
			pb.formatSpeed(pb.speed),
		)
	})

	return pb
}

func (pb *ProgressBar) formatSize(size float64) string {
	if size < 1024 {
		return fmt.Sprintf("%7fB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%7.2fKB", size/float64(1024))
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%7.2fMB", size/float64(1024*1024))
	} else if size < 1024*1024*1024*1024 {
		return fmt.Sprintf("%7.2fGB", size/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%7.2fTB", size/float64(1024*1024*1024*1024))
	}
}

func (pb *ProgressBar) formatSpeed(speed float64) string {
	return pb.formatSize(speed) + "/s"
}

func (pb *ProgressBar) Set(n int) error {
	interval := time.Now().Sub(pb.lastTime)
	if interval > pb.speedInterval {
		pb.speed = float64(n-pb.intervalStartN) / interval.Seconds()
		pb.intervalStartN = n
		pb.lastTime = time.Now()
	}

	pb.avgSpeed = float64(n) / time.Now().Sub(pb.startTime).Seconds()
	pb.n = n

	return pb.bar.Set(n)
}
