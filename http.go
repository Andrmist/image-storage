package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

func HttpWorker(ctx context.Context, minioClient *minio.Client) {
	minioEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioPrefix := os.Getenv("MINIO_PREFIX_PATH")
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(120 * time.Second))

	r.Post("/photo", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		id := uuid.New()

		path := filepath.Join(minioPrefix, fmt.Sprintf("%v%v", id, filepath.Ext(handler.Filename)))
		if _, err = minioClient.PutObject(ctx, minioBucket, path, file, handler.Size, minio.PutObjectOptions{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("%s/%s", minioEndpoint, filepath.Join(minioBucket, path))))
	})
	r.Get("/url", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		log.Println(url)
		res, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := uuid.New()

		path := filepath.Join(minioPrefix, fmt.Sprintf("%v%v", id, filepath.Ext(url)))
		if _, err = minioClient.PutObject(ctx, minioBucket, path, res.Body, res.ContentLength, minio.PutObjectOptions{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("%s/%s", minioEndpoint, filepath.Join(minioBucket, path))))
	})
	go func() {
		if err := http.ListenAndServe("0.0.0.0:8081", r); err != nil {
			log.Fatal(errors.Wrap(err, "failed to create http server"))
		}
	}()
}
