package main

import (
	_ "CeylonPlatform/middleware/authentication"
	"CeylonPlatform/middleware/initialization"
	_ "CeylonPlatform/service"
	_ "CeylonPlatform/service/authorize"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 监听系统信号，初始化主程序环境
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.Default().SetPrefix("[main]")
	rand.Seed(time.Now().Unix())

	// 监听退出信号
	go initialization.ListenSignal(sigs)

	// 初始化运行实体
	err := initialization.InitEntities("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	// 准备运行服务
	err = initialization.StartUp()
	if err != nil {
		initialization.Logger.Fatal(err.Error())
	}

	// 开始服务
	initialization.Serve()
}
