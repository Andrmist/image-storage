package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

func DeleteWorker() {
	log.Printf("delete worker started at %v", time.Now())
	files, err := os.ReadDir("./photos")
	if err != nil {
		log.Print(errors.Wrap(err, "failed to get files in delete worker"))
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				log.Print(errors.Wrap(err, fmt.Sprintf("failed to get file info on %v in delete worker", file.Name())))
			}
			delay := time.Now().Sub(info.ModTime())
			if delay.Hours() > 24 {
				log.Print(delay)
				if err := os.Remove(file.Name()); err != nil {
					log.Print(errors.Wrap(err, fmt.Sprintf("failed to remove file %v in delete worker", file.Name())))
				}
			}
		}
	}
}
