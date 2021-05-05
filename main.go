package main

import (
	"bilibili_download/config"
	"bilibili_download/download"
	"flag"
	"fmt"
)

var (
	cfg  = flag.String("c", "cfg.yml", "The path of configuration file")
)

func init() {
	flag.Parse()
	config.ParseConfig(*cfg)
}
func main() {
	fmt.Printf("\033[2J")
	config:=config.Config()
	if len(config.RoomInfo)==0 {
		fmt.Println("没有任何任务,退出")
		return
	}
	for _, info := range config.RoomInfo {
		i:=download.MakeNewManager(info)
		//i.Run()
		go func(c *download.DownloadManager) {
			c.Run()
		}(i)
	}
	<- make(chan string)
}
