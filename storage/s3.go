package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/marmotcai/uploadagent/logger"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// S3 - Amazon S3 storage
//
// type: s3
// bucket: UploadAgent-test
// region: us-east-1
// path: backups
// access_key_id: your-access-key-id
// secret_access_key: your-secret-access-key
// max_retries: 5
// timeout: 300
type S3 struct {
	Base
	bucket string
	path   string
	client *s3manager.Uploader
}


func cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}

func (ctx *S3) open() (err error) {
	ctx.viper.SetDefault("region", "local")

	cfg := aws.NewConfig()
	endpoint := ctx.viper.GetString("url")
	if len(endpoint) > 0 {
		cfg.Endpoint = aws.String(endpoint)
	}

	user := strings.Split(ctx.viper.GetString("user"), ":")
	if (len(user) >= 2) {
		cfg.Credentials = credentials.NewStaticCredentials(
			user[0],
			user[1],
			ctx.viper.GetString("token"),
		)
	}

	cfg.Region = aws.String(ctx.viper.GetString("region"))
	cfg.MaxRetries = aws.Int(ctx.viper.GetInt("max_retries"))

	cfg.S3ForcePathStyle = aws.Bool(ctx.viper.GetBool("forcepath_style")) //aws.Bool(true)
	cfg.DisableSSL = aws.Bool(true)

	ctx.bucket = ctx.viper.GetString("bucket")
	ctx.path = ctx.viper.GetString("path")

	sess := session.Must(session.NewSession(cfg))
	ctx.client = s3manager.NewUploader(sess)
	if (ctx.client == nil) {
		logger.Info("NewUploader failed!\n")
	}

	svc := s3.New(sess)
	resp, _:= svc.ListBuckets(&s3.ListBucketsInput{})
	for _, bucket := range resp.Buckets {
		fmt.Println(*bucket.Name)
	}

	return
}

func (ctx *S3) close() {}

func (ctx *S3) uploadfile(fileKey, filepath, remotepath string) (string, error) {
	logger.Info("open file :", filepath)
	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filepath, err)
	}
	defer f.Close()

	bucketPath := path.Join(remotepath, fileKey)
	input := &s3manager.UploadInput{
		Bucket: aws.String(ctx.bucket),
		Key:    aws.String(bucketPath),
		Body:   f,
	}

	logger.Info("-> S3 Uploading :", *input.Bucket, *input.Key)
	result, err := ctx.client.Upload(input)
	if err != nil {
		logger.Info("failed to upload file, %v", err)
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	url := remotepath + fileKey
	// url = "s3://" + ctx.bucket + url
	logger.Info("=>", result.Location, url)

	return url, nil
}

func (ctx *S3) delete(fileKey string) (err error) {
	remotePath := path.Join(ctx.path, fileKey)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(ctx.bucket),
		Key:    aws.String(remotePath),
	}
	_, err = ctx.client.S3.DeleteObject(input)
	return
}
