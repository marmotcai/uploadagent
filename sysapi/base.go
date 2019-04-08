package sysapi

import (
	"fmt"
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/spf13/viper"
)

// Base api
type Base struct {
	model       		config.ModelConfig
	url					string
	viper       		*viper.Viper
}

type MI struct {
	Code 				string      	`json:"code"`
	Type 				string      	`json:"type"`
	Duration			int				`json:"duration"`
	SourceDrmType 		string			`json:"sourceDrmType"`
	DestDrmType 		string			`json:"destDrmType"`
	AudioType 			string			`json:"audioType"`
	ScreenFormat 		string			`json:"screenFormat"`
	ClosedCaptioning 	string			`json:"closedCaptioning"`
	MediaSpec 			string			`json:"mediaSpec"`
	BitrateType 		string			`json:"bitrateType"`
}

func NewMI() *MI {
	return &MI{}
}

type Context interface {
	Post(url string, data []byte) error
}

func newBase(model config.ModelConfig, apiConfig config.SubConfig) (base Base) {
	base = Base{
		model:    	model,
		viper:    	apiConfig.Viper,
		url:     	apiConfig.Name,
	}
	return
}

// Run api
func runModel(model config.ModelConfig, apiConfig config.SubConfig, data []byte) (err error) {
	logger.Info("------------- API --------------")

	base := newBase(model, apiConfig)
	var ctx Context
	switch apiConfig.Type {
	case "rest":
		ctx = &Rest{Base: base}
	default:
		return fmt.Errorf("[%s] api type has not implement", apiConfig.Type)
	}

	logger.Info("=> API | " + apiConfig.Type, ":")

	ctx.Post(apiConfig.Viper.GetString("Url"), data)

	return nil
}

// Run api
func Run(model config.ModelConfig, data []byte) error {
	if len(model.Api) == 0 {
		return nil
	}

	logger.Info("------------- API -------------")
	for _, apiCfg := range model.Api {
		err := runModel(model, apiCfg, data)
		if err != nil {
			return err
		}
	}
	logger.Info("------------- API -------------\n")

	return nil
}
