package storage

import (
	"fmt"
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/spf13/viper"
	"path"
	"strings"
)

// Base storage
type Base struct {
	model       config.ModelConfig
	archivePath string
	viper       *viper.Viper
	keep        int
}

// Context storage interface
type Context interface {
	open() error
	close()
	uploadfile(fileKey, filepath, remotepath string) (string, error)
	delete(fileKey string) error
}

type UploadComplete func(model config.ModelConfig, data []byte) (error)

func newBase(model config.ModelConfig, archivePath string) (base Base) {
	base = Base{
		model:       model,
		archivePath: archivePath,
		viper:       model.StoreWith.Viper,
	}

	if base.viper != nil {
		base.keep = base.viper.GetInt("keep")
	}

	return
}

func getremotepath(defaultpath, destpath string) (string) {
	remotepath := defaultpath

	if (len(destpath) > 0) {
		remotepath = path.Join(remotepath, destpath)
	}
	return remotepath
}

func uploadfile(ctx Context, model config.ModelConfig, filekey string, filepath string, remotepath string, complete UploadComplete) (error) {
	remoteurl, err := ctx.uploadfile(filekey, filepath, remotepath)
	if (err != nil) {
		return fmt.Errorf("upload %s error : $s\n", filekey, err)
	}

	jsonfilekey := filekey + ".info"
	jsonfile := config.TempPath + "/" + jsonfilekey
	data, err := helper.GetJsondata(filekey, filepath, remoteurl, jsonfile)
	if (err != nil) {
		logger.Info("get mediainfo json error : %s", err)
	}
	if (helper.PathExists(jsonfile)) {
		_, err := ctx.uploadfile(jsonfilekey, jsonfile, remotepath)
		if (err != nil) {
			logger.Info("upload %s error : $s\n", jsonfile, err)
		}
	}
	err = complete(model, data)
	if (err != nil) {
		return fmt.Errorf("complete post %s error : $s\n", filekey, err)
	}

	return nil
}

func uploaddir(ctx Context, model config.ModelConfig, path string, complete UploadComplete) (error) {
	logger.Info("upload ", path)

	files, _ := helper.GetFilelist(path, "")
	for _, filepath := range files {
		logger.Info("-- ", filepath)
		if (helper.IsTempfile(filepath)) {
			continue
		}

		destpath, filekey := helper.GetFileKey(filepath, model.StoreWith.Viper.GetString("FileKeyFormat"))
		remotepath := getremotepath(model.StoreWith.Viper.GetString("path"), destpath)
		logger.Info("get file key :", filekey, remotepath)

		err := uploadfile(ctx, model, filekey, filepath, remotepath, complete)
		if (err != nil) {
			return fmt.Errorf("uploadfile %s error : $s\n", filekey, err)
		}
	}
	return nil
}

// Run storage
func Run(model config.ModelConfig, archivePath string, complete UploadComplete) (err error) {
	logger.Info("------------- Storage --------------")

	base := newBase(model, archivePath)
	var ctx Context
	switch model.StoreWith.Type {
	case "local":
		ctx = &Local{Base: base}
	case "ftp":
		ctx = &FTP{Base: base}
	case "scp":
		ctx = &SCP{Base: base}
	case "s3":
		ctx = &S3{Base: base}
	case "oss":
		ctx = &OSS{Base: base}
	default:
		return fmt.Errorf("[%s] storage type has not implement", model.StoreWith.Type)
	}

	logger.Info("=> Storage | " + model.StoreWith.Type)
	err = ctx.open()
	if err != nil {
		return err
	}
	defer ctx.close()

	if (archivePath == "") {
		includes := strings.Split(model.Archive.GetString("includes"), "|")
		includes = cleanPaths(includes)

		excludes := model.Archive.GetStringSlice("excludes")
		excludes = cleanPaths(excludes)

		for _, include := range includes {
			uploaddir(ctx, model, include, complete)
		}

	} else {
		destpath, filekey := helper.GetFileKey(archivePath, model.StoreWith.Viper.GetString("FileKeyFormat"))
		remotepath := getremotepath(model.StoreWith.Viper.GetString("path"), destpath)
		fmt.Printf("Get file key : (%s to %s)", filekey, remotepath)

		err := uploadfile(ctx, model, filekey, archivePath, remotepath, complete)
		if (err != nil) {
			return fmt.Errorf("uploadfile %s error : $s\n", filekey, err)
		}

		cycler := Cycler{}
		cycler.run(model.Name, filekey, base.keep, ctx.delete)
	}

	logger.Info("------------- Storage --------------\n")
	return nil
}
