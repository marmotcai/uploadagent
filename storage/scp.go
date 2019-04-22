package storage

import (
	"fmt"
	"github.com/marmotcai/uploadagent/helper"
	"golang.org/x/crypto/ssh"
	"os"
	"path"
	"strings"
	"time"
	// "crypto/tls"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/marmotcai/uploadagent/logger"
)

// SCP storage
//
// type: scp
// host: 192.168.1.2
// port: 22
// username: root
// password:
// timeout: 300
// private_key: ~/.ssh/id_rsa
type SCP struct {
	Base
	path       string
	host       string
	port       string
	privateKey string
	username   string
	password   string
	client     scp.Client
}

func (ctx *SCP) check(fileKey string) (error) {
	panic("implement me")
}

func (ctx *SCP) uploadfile(fileKey, filepath, remotepath string) (string, error) {
	err := ctx.client.Connect()
	if err != nil {
		return "", err
	}
	defer ctx.client.Session.Close()

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	remotePath := path.Join(ctx.path, fileKey)
	logger.Info("-> scp", remotePath)
	err = ctx.client.CopyFromFile(*file, remotePath, "0655")
	if err != nil {
		fmt.Printf("upload error (%s)\n", err)
		return "", err
	}
	logger.Info("Store successed")

	return remotePath, nil
}

func (ctx *SCP) open() (err error) {

	ctx.viper.SetDefault("timeout", 300)
	ctx.viper.SetDefault("private_key", "~/.ssh/id_rsa")

	ctx.port = "22"
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

	ctx.privateKey = helper.ExplandHome(ctx.viper.GetString("private_key"))
	var clientConfig ssh.ClientConfig
	logger.Info("PrivateKey", ctx.privateKey)
	clientConfig, err = auth.PrivateKey(
		ctx.username,
		ctx.privateKey,
		ssh.InsecureIgnoreHostKey(),
	)
	if err != nil {
		logger.Warn(err)
		logger.Info("PrivateKey fail, Try User@Host with Password")
		clientConfig = ssh.ClientConfig{
			User:            ctx.username,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}
	clientConfig.Timeout = ctx.viper.GetDuration("timeout") * time.Second
	if len(ctx.password) > 0 {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(ctx.password))
	}

	ctx.client = scp.NewClient(ctx.host+":"+ctx.port, &clientConfig)

	err = ctx.client.Connect()
	if err != nil {
		return err
	}
	defer ctx.client.Session.Close()
	ctx.client.Session.Run("mkdir -p " + ctx.path)
	return
}

func (ctx *SCP) close() {}

func (ctx *SCP) upload(fileKey string) (err error) {
	destpath, filekey := helper.GetFileKey(ctx.archivePath, ctx.model.StoreWith.Viper.GetString("FileKeyFormat"))

	ctx.uploadfile(filekey, ctx.archivePath, destpath)

	return nil
}

func (ctx *SCP) delete(fileKey string) (err error) {
	err = ctx.client.Connect()
	if err != nil {
		return
	}
	defer ctx.client.Session.Close()

	remotePath := path.Join(ctx.path, fileKey)
	logger.Info("-> remove", remotePath)
	err = ctx.client.Session.Run("rm " + remotePath)
	return
}
