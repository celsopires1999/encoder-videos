package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(ctx context.Context, objectPath string, client *storage.Client) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})
	return err
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	if err := vu.loadPaths(); err != nil {
		return err
	}

	ctx, uploadClient, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(ctx, in, returnChannel, uploadClient)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	countDoneWorker := 0
	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
		if countDoneWorker == len(vu.Paths) {
			close(in)
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(ctx context.Context, in chan int, returnChannel chan string, uploadClient *storage.Client) {
	for x := range in {
		if err := vu.UploadObject(ctx, vu.Paths[x], uploadClient); err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("Error uploading file: %v. Error: %v", vu.Paths[x], err)
			returnChannel <- err.Error()
		}
		returnChannel <- ""
	}
	returnChannel <- "upload completed"
}

func getClientUpload() (context.Context, *storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return ctx, client, nil
}
