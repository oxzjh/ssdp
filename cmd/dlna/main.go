package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/oxzjh/ssdp"
	"github.com/oxzjh/ssdp/dlna"
)

func main() {
	var (
		seconds int64
	)
	flag.Int64Var(&seconds, "s", 3, "Seconds")
	flag.Parse()

	fmt.Println("设备扫描中...")
	devices, err := dlna.FindDevices(time.Second * time.Duration(seconds))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(devices) == 0 {
		fmt.Println("未发现设备")
		return
	}

	var selectedDevice *ssdp.Device
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请选择设备：")
	for {
		for i, d := range devices {
			fmt.Printf("\t%d - %s\n", i+1, d.FriendlyName)
		}
		if !scanner.Scan() {
			fmt.Println("读取输入失败")
			return
		}
		id, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if id < 1 || id > len(devices) {
			fmt.Println("设备不存在，请重新选择：")
			continue
		}
		selectedDevice = devices[id-1]
		break
	}
	var file string
	fmt.Print("请输入视频文件路径：")
	for {
		if !scanner.Scan() {
			fmt.Println("读取输入失败")
			return
		}
		file = scanner.Text()
		if fi, err := os.Stat(file); err == nil && !fi.IsDir() {
			break
		}
		fmt.Print("视频不存在，请重现输入：")
	}

	server, err := dlna.NewServer(selectedDevice.LocalIP, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	d := dlna.New(selectedDevice)
	if err = d.SetAVTransportURI(server.GetURL()); err == nil {
		err = d.Play("1")
	}
	if err != nil {
		server.Close()
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("视频播放中...")
	fmt.Println(server.Serve())
}
