package database

import (
	"database/sql"
	"fmt"
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/spf13/viper"
	"path"
)

// Base database
type Base struct {
	model    config.ModelConfig
	dbConfig config.SubConfig
	viper    *viper.Viper
	name     string
	dumpPath string
}

// Context database interface
type DBContext interface {
	Perform() error
	GetDBObj() *sql.DB
}

func newBase(model config.ModelConfig, dbConfig config.SubConfig) (base Base) {
	base = Base{
		model:    model,
		dbConfig: dbConfig,
		viper:    dbConfig.Viper,
		name:     dbConfig.Name,
	}
	base.dumpPath = path.Join(model.DumpPath, dbConfig.Type, base.name)
	helper.MkdirP(base.dumpPath)
	return
}

func GetDBModel(model config.ModelConfig, dbConfig config.SubConfig) (DBContext) {
	if (len(dbConfig.Type) == 0) {
		return nil
	}

	base := newBase(model, dbConfig)
	var ctx DBContext
	switch dbConfig.Type {
	case "mysql":
		ctx = &MySQL{Base: base}
	case "redis":
		ctx = &Redis{Base: base}
	case "postgresql":
		ctx = &PostgreSQL{Base: base}
	case "mongodb":
		ctx = &MongoDB{Base: base}
	default:
		logger.Warn(fmt.Errorf("model: %s databases.%s config `type: %s`, but is not implement", model.Name, dbConfig.Name, dbConfig.Type))
		return nil
	}

	logger.Info("=> database |", dbConfig.Type, ":", base.name)

	return ctx
}


// New - initialize Database
func runModel(model config.ModelConfig, dbConfig config.SubConfig) (err error) {

	ctx := GetDBModel(model, dbConfig)
	// perform
	err = ctx.Perform()
	if err != nil {
		return err
	}
	logger.Info("")

	return
}

// Run databases
func Run(model config.ModelConfig) error {
	if len(model.Databases) == 0 {
		return nil
	}

	logger.Info("------------- Databases -------------")
	for _, dbCfg := range model.Databases {
		err := runModel(model, dbCfg)
		if err != nil {
			return err
		}
	}
	logger.Info("------------- Databases -------------")

	return nil
}
