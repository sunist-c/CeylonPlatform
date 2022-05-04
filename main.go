package main

import (
	"CeylonPlatform/middleware/initialization"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Default().SetFlags(log.Ldate | log.Ltime)
	log.Default().SetPrefix("[main]")

	err := initialization.InitEntities("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	err = initialization.StartUp()
	if err != nil {
		log.Fatalln(err)
	}

	go initialization.ListenSignal(sigs)
}
