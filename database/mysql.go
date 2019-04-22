package database

import (
	"database/sql"
	"fmt"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"path"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL database
//
// type: mysql
// host: 127.0.0.1
// port: 3306
// database:
// username: root
// password:
// additional_options:
type MySQL struct {
	Base
	host              string
	port              string
	database          string
	username          string
	password          string
	additionalOptions []string
}

func (ctx *MySQL) linkstr() string {
	viper := ctx.viper

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("username", "root")
	viper.SetDefault("port", 3306)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.database = viper.GetString("database")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", ctx.username, ctx.password, ctx.host, ctx.port, ctx.database)
}

func (ctx *MySQL) perform() (err error) {
	sqlstr := ctx.linkstr()

	viper := ctx.viper

	addOpts := viper.GetString("additional_options")
	if len(addOpts) > 0 {
		ctx.additionalOptions = strings.Split(addOpts, " ")
	}

	// mysqldump command
	if len(ctx.database) == 0 {
		return fmt.Errorf("mysql database config is required")
	}

	db, err := sql.Open("mysql", sqlstr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	// err = ctx.dump()
	return
}

func (ctx *MySQL) dumpArgs() []string {
	dumpArgs := []string{}
	if len(ctx.host) > 0 {
		dumpArgs = append(dumpArgs, "--host", ctx.host)
	}
	if len(ctx.port) > 0 {
		dumpArgs = append(dumpArgs, "--port", ctx.port)
	}
	if len(ctx.username) > 0 {
		dumpArgs = append(dumpArgs, "-u", ctx.username)
	}
	if len(ctx.password) > 0 {
		dumpArgs = append(dumpArgs, `-p`+ctx.password)
	}
	if len(ctx.additionalOptions) > 0 {
		dumpArgs = append(dumpArgs, ctx.additionalOptions...)
	}

	dumpArgs = append(dumpArgs, ctx.database)
	dumpFilePath := path.Join(ctx.dumpPath, ctx.database+".sql")
	dumpArgs = append(dumpArgs, "--result-file="+dumpFilePath)
	return dumpArgs
}

func (ctx *MySQL) dump() error {
	logger.Info("-> Dumping MySQL...")
	_, err := helper.Exec("mysqldump", ctx.dumpArgs()...)
	if err != nil {
		return fmt.Errorf("-> Dump error: %s", err)
	}
	logger.Info("dump path:", ctx.dumpPath)
	return nil
}
