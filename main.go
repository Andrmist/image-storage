package main

import (
	// "github.com/robfig/cron"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// cr := cron.New()
	// cr.AddFunc("0 0 0 * * *", DeleteWorker)
	// cr.Start()
	ctx, cancel := context.WithCancel(context.Background())
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioAccessSecret := os.Getenv("MINIO_ACCESS_SECRET")
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioAccessSecret, ""),
		Secure: false,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to minio: %s", err))
	}

	HttpWorker(ctx, minioClient)
	log.Println("http worker started")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	log.Println("Shutdown...")
	cancel()
	// cr.Stop()
	os.Exit(0)
}
