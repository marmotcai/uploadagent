package storage

import (
	"fmt"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"os"
	"path"
	"strings"

	"github.com/secsy/goftp"
	"time"
)

// FTP storage
//
// type: ftp
// path: /backups
// host: ftp.your-host.com
// port: 21
// timeout: 30
// username:
// password:
type FTP struct {
	Base
	path     string
	host     string
	port     string
	username string
	password string

	client *goftp.Client
}

func (ctx *FTP) MkdirP(remotepath string) (error) {
	params := strings.Split(remotepath, "/")
	dir := ""
	for j := 0; j < len(params); j++ {
		if (len(params[j]) <= 0) {
			continue
		}

		dir = path.Join(dir, params[j])
		_, err := ctx.client.Stat(dir)
		if (err != nil) {
			if _, err := ctx.client.Mkdir(dir); err != nil {
				continue
				// return err
			}
		}
	}
	return nil
}

func (ctx *FTP) uploadfile(fileKey, filepath, remotepath string) (string, error) {

	logger.Info("-> Uploading...")

	err := ctx.MkdirP(remotepath)
	if err != nil {
		return "", err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	remotefilepath := path.Join(remotepath, fileKey)

	err = ctx.client.Store(remotefilepath, file)
	if err != nil {
		fmt.Printf("upload error (%s)\n", err)
		return "", err
	}

	logger.Info("Store successed")
	return remotefilepath, nil
}

func (ctx *FTP) open() (err error) {
	ctx.viper.SetDefault("port", "21")
	ctx.viper.SetDefault("timeout", 300)

	host := strings.Split(ctx.viper.GetString("url"), ":")
	ctx.host = host[0]
	if (len(host) >= 2) {
		ctx.port = host[1]
	}

	ctx.path = ctx.viper.GetString("path")

	user := strings.Split(ctx.viper.GetString("user"), ":")
	if (len(user) >= 2) {
		ctx.username = user[0]
		ctx.password = user[1]
	}

	ftpConfig := goftp.Config{
		User:     ctx.username,
		Password: ctx.password,
		Timeout:  ctx.viper.GetDuration("timeout") * time.Second,
	}
	ctx.client, err = goftp.DialConfig(ftpConfig, ctx.host + ":" + ctx.port)
	if err != nil {
		return err
	}
/*
	fileinfo, err := ctx.client.ReadDir(ctx.path)
	if (fileinfo == nil) {
		path, err := ctx.client.Mkdir(ctx.path)
		if err != nil {
			fmt.Printf("%s already exists \n", path)
		}
	}
*/
	return nil
}

func (ctx *FTP) close() {
	ctx.client.Close()
}

func (ctx *FTP) upload(fileKey string) (error) {
	destpath, filekey := helper.GetFileKey(ctx.archivePath, ctx.model.StoreWith.Viper.GetString("FileKeyFormat"))
	_, err := ctx.uploadfile(filekey, ctx.archivePath, destpath)
	return err
}

func (ctx *FTP) delete(fileKey string) (err error) {
	remotePath := path.Join(ctx.path, fileKey)
	err = ctx.client.Delete(remotePath)
	return
}
