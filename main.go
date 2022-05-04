package main

import (
	"CeylonPlatform/middleware/initialization"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 监听系统信号，初始化主程序环境
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Default().SetFlags(log.Ldate | log.Ltime)
	log.Default().SetPrefix("[main]")

	// 初始化运行实体
	err := initialization.InitEntities("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	// 运行服务
	err = initialization.StartUp()
	if err != nil {
		initialization.Logger.Fatal(err.Error())
	}

	// 监听退出信号
	go initialization.ListenSignal(sigs)
}
