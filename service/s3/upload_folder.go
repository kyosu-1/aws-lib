package s3

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// 並行処理でS3にフォルダの中身をアップロードする
func UploadFolder(localPath string, bucket string, prefix string) error {

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
		return err
	}

	// For each file found walking, upload it to Amazon S3
	uploader := manager.NewUploader(s3.NewFromConfig(cfg))

	var g errgroup.Group

	for path := range walker {
		p := path
		g.Go(func() error {
			rel, err := filepath.Rel(localPath, p)
			if err != nil {
				return err
			}
			file, err := os.Open(p)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    aws.String(filepath.Join(prefix, rel)),
				Body:   file,
			})
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

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
