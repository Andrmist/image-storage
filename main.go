package main

import (
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cr := cron.New()
	// cr.AddFunc("0 0 0 * * *", DeleteWorker)
	cr.Start()
	HttpWorker()
	log.Println("http worker started")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	log.Println("Shutdown...")
	cr.Stop()
	os.Exit(0)
}
