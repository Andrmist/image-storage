package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func HttpWorker() {
	HostName := os.Getenv("HOSTNAME")
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

		path := filepath.Join("photos", fmt.Sprintf("%v%v", id, filepath.Ext(handler.Filename)))
		dst, err := os.Create(path)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v/%v", HostName, path)))
	})
	go func() {
		if err := http.ListenAndServe("0.0.0.0:8081", r); err != nil {
			log.Fatal(errors.Wrap(err, "failed to create http server"))
		}
	}()
}
