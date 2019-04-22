package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/marmotcai/uploadagent/logger"
	"path"
)

// OSS - Aliyun OSS storage
//
// type: oss
// bucket: UploadAgent-test
// endpoint: oss-cn-beijing.aliyuncs.com
// path: /
// access_key_id: your-access-key-id
// access_key_secret: your-access-key-secret
// max_retries: 5
// timeout: 300
type OSS struct {
	Base
	endpoint        string
	bucket          string
	accessKeyID     string
	accessKeySecret string
	path            string
	maxRetries      int
	timeout         int
	client          *oss.Bucket
}

func (ctx *OSS) check(fileKey string) (error) {
	panic("implement me")
}

func (ctx *OSS) uploadfile(fileKey, filepath, remotepath string) (string, error) {
	remotePath := path.Join(ctx.path, fileKey)

	logger.Info("-> Uploading OSS...")
	err := ctx.client.UploadFile(remotePath, ctx.archivePath, ossPartSize, oss.Routines(4))

	if err != nil {
		return "", err
	}
	logger.Info("Success")

	return fileKey, nil
}

var (
	// 4 Mb
	ossPartSize int64 = 4 * 1024 * 1024
)

func (ctx *OSS) open() (err error) {
	ctx.viper.SetDefault("endpoint", "oss-cn-beijing.aliyuncs.com")
	ctx.viper.SetDefault("max_retries", 3)
	ctx.viper.SetDefault("path", "/")
	ctx.viper.SetDefault("timeout", 300)

	ctx.endpoint = ctx.viper.GetString("endpoint")
	ctx.bucket = ctx.viper.GetString("bucket")
	ctx.accessKeyID = ctx.viper.GetString("access_key_id")
	ctx.accessKeySecret = ctx.viper.GetString("access_key_secret")
	ctx.path = ctx.viper.GetString("path")
	ctx.maxRetries = ctx.viper.GetInt("max_retries")
	ctx.timeout = ctx.viper.GetInt("timeout")

	logger.Info("endpoint:", ctx.endpoint)
	logger.Info("bucket:", ctx.bucket)

	ossClient, err := oss.New(ctx.endpoint, ctx.accessKeyID, ctx.accessKeySecret)
	if err != nil {
		return err
	}
	ossClient.Config.Timeout = uint(ctx.timeout)
	ossClient.Config.RetryTimes = uint(ctx.maxRetries)

	ctx.client, err = ossClient.Bucket(ctx.bucket)
	if err != nil {
		return err
	}

	return
}

func (ctx *OSS) close() {
}

func (ctx *OSS) upload(fileKey string) (err error) {
	remotePath := path.Join(ctx.path, fileKey)

	logger.Info("-> Uploading OSS...")
	err = ctx.client.UploadFile(remotePath, ctx.archivePath, ossPartSize, oss.Routines(4))

	if err != nil {
		return err
	}
	logger.Info("Success")

	return nil
}

func (ctx *OSS) delete(fileKey string) (err error) {
	remotePath := path.Join(ctx.path, fileKey)
	err = ctx.client.DeleteObject(remotePath)
	return
}
