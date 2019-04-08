package main

import (
	"encoding/json"
	"fmt"
	"github.com/marmotcai/uploadagent/archive"
	"github.com/marmotcai/uploadagent/compressor"
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/database"
	"github.com/marmotcai/uploadagent/encryptor"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/marmotcai/uploadagent/storage"
	"github.com/marmotcai/uploadagent/sysapi"
	"os"
	"path"
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

func PostMMS(model config.ModelConfig, key, localfile, url string) (error) {
	if (!helper.IsMediafile(localfile)) {
		return nil
	}

	if (model.Api == nil) {
		return nil
	}

	ms, err := helper.GetMediainfo(localfile)
	//先创建一个目标类型的实例对象，用于存放解码后的值
	var inter interface{}
	err = json.Unmarshal([]byte(ms), &inter)
	if err != nil {
		return fmt.Errorf("error in translating,", err.Error())
	}
	//要访问解码后的数据结构，需要先判断目标结构是否为预期的数据类型
	mi, ok := inter.(map[string]interface{})
	//然后通过for循环一一访问解码后的目标数据
	if ok {
		for k, v := range mi {
			switch v.(type) {
			case interface{}:
				{
					v, _ := v.(map[string]interface{})

					if (k == "general") {
						name, err := helper.GetNameFromPath(localfile)
						if (err != nil) {
							name = path.Base(localfile)
						}

						v["path"] = localfile
						v["url"] = url
						v["name"] = name
						v["filekey"] = key
					}
				}

			default:
				fmt.Println("illegle type")
			}
		}
	}
	data, err := json.Marshal(inter)
	
 	str := string(data[:])
	fmt.Println(str)
	
	sysapi.Run(model, data)

	return nil
}

// Perform model
func (ctx Model) perform() {
	logger.Info("======== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")
	defer ctx.cleanup()

	if len(ctx.Config.Databases) != 0 {
		err := database.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	var archivePath string

	if (ctx.Config.PackWith.Type != "none") {
		if ctx.Config.Archive != nil {
			err := archive.Run(ctx.Config)
			if err != nil {
				logger.Error(err)
				return
			}
		}

		archivePath, err := compressor.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}

		archivePath, err = encryptor.Run(archivePath, ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	err := storage.Run(ctx.Config, archivePath, PostMMS)
	if err != nil {
		logger.Error(err)
		return
	}

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
