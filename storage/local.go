package storage

import (
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"path"
)

// Local storage
//
// type: local
// path: /data/backups
type Local struct {
	Base
	isMove bool
	destPath string
	filekeyformat string
}

func (ctx *Local) uploadfile(fileKey, filepath, remotepath string) (string, error) {

	helper.MkdirP(remotepath)

	remotefilepath := path.Join(remotepath, fileKey)

	cmdstr := "cp"
	if (ctx.model.StoreWith.Viper.GetBool("IsMove")) {
		cmdstr = "mv"
	}
	_, err := helper.Exec(cmdstr, filepath, remotefilepath)

	if err != nil {
		return "", err
	}
	logger.Info("Store successed", ctx.destPath)

	if (cmdstr == "mv") {
		return remotefilepath, nil
	}

	url := remotefilepath
	// url = "file://" + url

	return url, nil
}

func (ctx *Local) open() (err error) {
	ctx.destPath = ctx.model.StoreWith.Viper.GetString("path")
	if (ctx.destPath != "" ) {
		helper.MkdirP(ctx.destPath)
	}
	ctx.isMove = false
	ctx.filekeyformat =  ctx.model.StoreWith.Viper.GetString("FileKeyFormat")
	return
}

func (ctx *Local) close() {}

func (ctx *Local) upload(fileKey string) (err error) {
	_, err = helper.Exec("cp", ctx.archivePath, ctx.destPath)
	if err != nil {
		return err
	}
	logger.Info("Store successed", ctx.destPath)
	return nil
}

func (ctx *Local) delete(fileKey string) (err error) {
	_, err = helper.Exec("rm", path.Join(ctx.destPath, fileKey))
	return
}
