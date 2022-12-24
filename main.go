package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type key struct {
	TestPhrase string `json:"test_phrase"`
	Key        string `json:"key"`
}

func main() {
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var in key

	err := readJSON(w, r, &in)
	if err != nil {
		return
	}

	if in.TestPhrase != os.Getenv("CHALLENGE") {
		log.Fatalln("failed access attempt", r.RemoteAddr)
	}
	log.Println("access authorized for", r.RemoteAddr)
	ctx := context.Background()
	err = uploadToGcs(ctx, in.Key)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("key reaved"))
	if err != nil {
		log.Println("failed to write the response", err)
	}
}

func NewStorageClient(ctx context.Context) (*storage.Client, error) {
	// ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return &storage.Client{}, fmt.Errorf("failed to create storage instance")
	}
	return client, nil
}

func readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	log.Println("reading input")
	// 666K
	maxBytes := 681984
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		log.Println("readJSON encountered a fatal error", err)
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("error parsing json")
	}
	return nil
}

func uploadToGcs(ctx context.Context, obj string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	store := os.Getenv("STORAGE")
	log.Println("bucket:", store)
	bucket, err := NewStorageClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	b := bucket.Bucket(store).Object("key.txt").NewWriter(ctx)
	b.ContentType = "text/plain"
	b.Metadata = map[string]string{
		"x-goog-meta-app": "application-tag",
		"x-goog-meta-bar": "bar",
	}
	if _, err := b.Write([]byte(obj + "\n")); err != nil {
		return fmt.Errorf("coudlnt write to bucket: %v", err)
	}
	if err := b.Close(); err != nil {
		return fmt.Errorf("couldnt close bucket (or save file): %v", err)
	}
	return nil
}
