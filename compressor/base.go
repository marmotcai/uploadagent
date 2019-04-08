package compressor

import (
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

// Base compressor
type Base struct {
	name  string
	model config.ModelConfig
	viper *viper.Viper
}

// Context compressor
type Context interface {
	perform() (archivePath string, err error)
}

func (ctx *Base) archiveFilePath(ext string) string {
	return path.Join(os.TempDir(), "UploadAgent", time.Now().Format("2006.01.02.15.04.05")+ext)
}

func newBase(model config.ModelConfig) (base Base) {
	base = Base{
		name:  model.Name,
		model: model,
		viper: model.PackWith.Viper,
	}
	return
}

// Run compressor
func Run(model config.ModelConfig) (archivePath string, err error) {
	base := newBase(model)

	var ctx Context
	switch model.PackWith.Type {
	case "tgz":
		ctx = &Tgz{Base: base}
	default:
		ctx = &Tgz{}
	}

	logger.Info("------------ Compressor -------------")
	logger.Info("=> Compress | " + model.PackWith.Type)

	// set workdir
	os.Chdir(path.Join(model.DumpPath, "../"))
	archivePath, err = ctx.perform()
	if err != nil {
		return
	}
	logger.Info("->", archivePath)
	logger.Info("------------ Compressor -------------\n")

	return
}
