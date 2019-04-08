package database

import (
	"fmt"
	"github.com/marmotcai/uploadagent/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Monkey struct {
	Base
}

func (ctx Monkey) perform() error {
	if ctx.model.Name != "TestMonkey" {
		return fmt.Errorf("Error")
	}
	if ctx.dbConfig.Name != "mysql1" {
		return fmt.Errorf("Error")
	}
	return nil
}

func TestBaseInterface(t *testing.T) {
	base := Base{
		model: config.ModelConfig{
			Name: "TestMonkey",
		},
		dbConfig: config.SubConfig{
			Name: "mysql1",
		},
	}
	ctx := Monkey{Base: base}
	err := ctx.perform()
	assert.Nil(t, err)
}

func TestBase_newBase(t *testing.T) {
	model := config.ModelConfig{
		DumpPath: "/tmp/github.com/marmotcai/uploadagent/test",
	}
	dbConfig := config.SubConfig{
		Type: "mysql",
		Name: "mysql-master",
	}
	base := newBase(model, dbConfig)

	assert.Equal(t, base.model, model)
	assert.Equal(t, base.dbConfig, dbConfig)
	assert.Equal(t, base.viper, dbConfig.Viper)
	assert.Equal(t, base.name, "mysql-master")
	assert.Equal(t, base.dumpPath, "/tmp/github.com/marmotcai/uploadagent/test/mysql/mysql-master")
}
