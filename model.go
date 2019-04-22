package main

import (
	"fmt"
	"github.com/marmotcai/uploadagent/archive"
	"github.com/marmotcai/uploadagent/compressor"
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/database"
	"github.com/marmotcai/uploadagent/encryptor"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/marmotcai/uploadagent/storage"
	"github.com/marmotcai/uploadagent/sysapi"
	"os"
)

// Model class
type Model struct {
	Config config.ModelConfig
}


func loadMediaifo(key string, value interface{}) (error) {
	v, ok := value.(map[string]interface{})
	if (!ok) {
		return fmt.Errorf("json Unmarshal error")
	}

	for key, value := range v {
		fmt.Printf("%v", key)
		fmt.Printf("%v", value)
	}

	switch key {
	case "general": {
	}
	case "video": {

	}
	case "audio": {

	}
	}
	return nil
}

func PostMMS(model config.ModelConfig, data []byte) (error) {
	if (data != nil) {
		sysapi.Run(model, data)
	}
	return nil
}

// Perform model
func (ctx Model) perform() (error) {

	logger.Info("work dir: ", ctx.Config.DumpPath)
	defer ctx.cleanup()

	if len(ctx.Config.Databases) != 0 {
		err := database.Run(ctx.Config)
		if err != nil {
			return err
		}
	}

	var archivePath string

	if (ctx.Config.PackWith.Type != "none") {
		if ctx.Config.Archive != nil {
			err := archive.Run(ctx.Config)
			if err != nil {
				return err
			}
		}

		archivePath, err := compressor.Run(ctx.Config)
		if err != nil {
			return err
		}

		archivePath, err = encryptor.Run(archivePath, ctx.Config)
		if err != nil {
			return err
		}
	}

	err := storage.Run(ctx.Config, archivePath, PostMMS)
	if err != nil {
		return err
	}

	return nil
}

func (ctx Model) check() (error) {
	if len(ctx.Config.Databases) != 0 {
		err := database.Run(ctx.Config)
		if err != nil {
			return err
		}
	}

	err := storage.Check(ctx.Config, "", nil)
	if err != nil {
		return err
	}

	return nil
}

// Cleanup model temp files
func (ctx Model) cleanup() {
	logger.Info("Cleanup temp dir:" + config.TempPath + "...\n")
	err := os.RemoveAll(config.TempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+config.TempPath+" error:", err)
	}
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")
}
