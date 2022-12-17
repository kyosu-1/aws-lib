package s3

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	localPath string = "master_data/images"
	bucket    string = "tokkun-test-masterdata"
	prefix    string = ""
)

func UploadFolder(localPath string, bucket string, prefix string) error {
	// caluculate time
	start := time.Now()
	defer func() {
		fmt.Printf("time: %f[s] \r ", time.Since(start).Seconds())
	}()

	walker := make(fileWalk)
	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localPath, walker.Walk); err != nil {
			log.Fatalln("Walk failed:", err)
		}
		close(walker)
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("error:", err)
	}

	// For each file found walking, upload it to Amazon S3
	uploader := manager.NewUploader(s3.NewFromConfig(cfg))

	var wg sync.WaitGroup

	for path := range walker {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			rel, err := filepath.Rel(localPath, path)
			if err != nil {
				log.Fatalln("Unable to get relative path:", path, err)
			}
			file, err := os.Open(path)
			if err != nil {
				log.Println("Failed opening file", path, err)
				return
			}
			defer file.Close()
			_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    aws.String(filepath.Join(prefix, rel)),
				Body:   file,
			})
			if err != nil {
				log.Fatalln("Failed to upload", path, err)
			}
		}(path)
	}

	wg.Wait()
	return nil
}

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}
