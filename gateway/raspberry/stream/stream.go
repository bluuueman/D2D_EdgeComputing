package stream

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var frames <-chan []uint8
var stopStream context.CancelFunc
var camera *device.Device

func send(frame []uint8, url string) {
	client := http.Client{}
	bodyBuf := &bytes.Buffer{}
	bodyWrite := multipart.NewWriter(bodyBuf)
	// file 为key
	fileWrite, err := bodyWrite.CreateFormFile("img", "table.jpg")
	_, err = fileWrite.Write(frame)
	if err != nil {
		fmt.Println("Write frame failed")
	}
	bodyWrite.Close() //要关闭，会将w.w.boundary刷写到w.writer中
	// 创建请求
	contentType := bodyWrite.FormDataContentType()
	req, err := http.NewRequest(http.MethodPost, url, bodyBuf)
	if err != nil {
		fmt.Println("err 2")
	}
	// 设置头
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err 3")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err 4")
	}
	fmt.Println(string(b))
	defer resp.Body.Close()
}

func initDevice() {
	devName := "/dev/video0"
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.Parse()

	// open device
	device, err := device.Open(
		devName,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMPEG, Width: 1920, Height: 1080}),
	)
	camera = device
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
}

func Streamer(url string, interval time.Duration, run *bool) {
	for frame := range frames {
		if *run == false {
			break
		}
		go send(frame, url)
		time.Sleep(interval)
	}
	fmt.Println("exit")
}

func InitStream() {
	initDevice()
	ctx, stop := context.WithCancel(context.TODO())
	if err := camera.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}
	frames = camera.GetOutput()
	stopStream = stop
}

func StopStream() {
	stopStream()
	defer camera.Close()
}
