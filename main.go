package main

import (
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	_ "github.com/spf13/viper"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func Perform(configfile, modelName string) {
	config.LoadConfig(configfile)

	for _, modelConfig := range config.Models {
		if (modelConfig.Name == modelName) || (len(modelName) == 0) {
			model := Model{
				Config: modelConfig,
			}

			logger.Info("======== " + modelConfig.Name + " ========")
			logger.Info("begin perform... ")

			err := model.perform()
			if (err != nil) {
				logger.Error(err)
			}

			logger.Info("end perform. ")
			logger.Info("======== " + modelConfig.Name + " ========")

			return
		}
	}
}

func Check(configfile, modelName string) {
	config.LoadConfig(configfile)

	for _, modelConfig := range config.Models {
		if (modelConfig.Name == modelName) || (len(modelName) == 0) {
			model := Model{
				Config: modelConfig,
			}

			logger.Info("======== " + modelConfig.Name + " ========")
			logger.Info("begin check... ")

			err := model.check()
			if (err != nil) {
				logger.Error(err)
			}

			logger.Info("end check. ")
			logger.Info("======== " + modelConfig.Name + " ========")

			return
		}
	}
}

//option demo : -logspath "./logs" -suffixw "rm/rmvb/mxf" -suffixb "exe/txt" config
//db demo : -dt "mysql" -dh "db.cloudgather.cn" -dp "33306" -dd "mms_test" -du "root" -dw "cg123456" config
//s3 demo : -st "s3" -surl "http://192.168.2.9:3090" -suser "4V1cweFJGTlhjM2hOUkVGM1RVUm9RV0l5U25GYVYwNHdURmhLTTA5WE9WcE5Wa1U5PNJI:4WVRGQ01GWklVak50VWpGamMyWmFZV014Y0ZWbFFUMDk=qEyE" -spath "/" -sregion "my-region" -sbucket "input" -ssorcepath_style "true" -keyformat "%CLASS_LAST0%/%HASH_TOP0%/%HASHFULL%" -l "./" exec
//scp demo : -st "scp" -surl "192.168.2.72:22" -suser "root:cg112233" -spath "/root/temp" config
//ftp demo : -st "ftp" -surl "192.168.2.9:21" -suser "caijun:aa112233" -spath "cloudgather/source/raw/senyu/series" -keyformat "%HASHFULL%" -oismove "false" -l "/Users/andrewcai/9/raw/guizhou/SY-01/电视剧/" coonfig
//local demo : -st "local" -spath "/Users/andrewcai/9/raw/" -keyformat "%CLASS_LAST2%/%HASH_TOP2%/%HASHFULL%" -oismove "false" -l "/Users/andrewcai/9/raw/chengdu/190124 YK 媒资/电影" config
//api demo : -at "rest" -aurl "http://192.168.2.7/restApi/movie/add" config
//check demo : check
func main() {
	app := cli.NewApp()

	app.Name = helper.App_name
	app.Version = helper.Version

	app.Usage = helper.Usage

	db := config.NewModel_DB()
	api := config.NewModel_API()

	store := config.NewModel_Store()

	option := config.NewModel_Option()

	var configfile string
	var localpath string
	var modelname string


	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8000,
			Usage: "listening port",
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "config file path",
			Destination: &configfile,
		},
		cli.StringFlag{
			Name:        "model, m",
			Usage:       "Model name that you want execute",
			Destination: &modelname,
		},

		//store
		cli.StringFlag{
			Name: "store-type, st", Usage: "store-type", Destination: &store.Type,
		},
		cli.StringFlag{
			Name: "store-url, surl", Usage: "store-url", Destination: &store.Url,
		},
		cli.StringFlag{
			Name: "store-user, suser", Usage: "store-user", Destination: &store.User,
		},
		cli.StringFlag{
			Name: "store-path, spath", Usage: "store-path", Destination: &store.Path,
		},
		cli.StringFlag{
			Name: "store-region, sregion", Usage: "store-region", Destination: &store.Region,
		},
		cli.StringFlag{
			Name: "store-Bucket, sbucket", Usage: "store-bucket", Destination: &store.Bucket,
		},
		cli.StringFlag{
			Name: "store-forcepath_style, ssorcepath_style", Usage: "store-forcepath_style", Destination: &store.ForcePath_style,
		},
		cli.StringFlag{
			Name: "store-timeout, stimeout", Usage: "store-timeout", Destination: &store.Timeout,
		},

		//option
		cli.StringFlag{
			Name: "option-ismove, oismove", Usage: "option ismove", Destination: &store.IsMove,
		},

		//naming rule
		cli.StringFlag{
			Name: "filekey-format, keyformat", Usage: "filekey format", Destination: &store.FileKeyFormat,
		},
		cli.StringFlag{
			Name: "naming-rule, namingrule", Usage: "naming rule", Destination: &store.NamingRule,
		},

		//db
		cli.StringFlag{
			Name: "db-type, dt",	Usage: "DB Type",	Destination: &db.Type,
		},
		cli.StringFlag{
			Name: "db-host, dh", Usage: "DB Host", Destination: &db.Host,
		},
		cli.StringFlag{
			Name: "db-port, dp", Usage: "DB Port", Destination: &db.Port,
		},
		cli.StringFlag{
			Name: "db-database, dd", Usage: "DB Database", Destination: &db.Database,
		},
		cli.StringFlag{
			Name:"db-username, du",	Usage: "DB Username", Destination: &db.Username,
		},
		cli.StringFlag{
			Name: "db-password, dw", Usage: "DB Password", Destination: &db.Password,
		},

		//api
		cli.StringFlag{
			Name: "api-type, at", Usage: "api type", Destination: &api.Type,
		},
		cli.StringFlag{
			Name: "api-url, aurl", Usage: "api url", Destination: &api.Url,
		},

		//source path
		cli.StringFlag{
			Name: "localpath, l", Usage: "local source path", Destination: &localpath,
		},

		//logs file path
		cli.StringFlag{
			Name: "Logs-Filepath, logspath", Usage: "logs file path", Destination: &option.LogsFilepath,
		},


		//suffix white
		cli.StringFlag{
			Name: "suffix-white, suffixw", Usage: "suffix-white", Destination: &option.Suffixwhite,
		},

		//suffix black
		cli.StringFlag{
			Name: "suffix-black, suffixb", Usage: "suffix-black", Destination: &option.Suffixblack,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"config"},
			Usage:   "write config file",
			Action: func(c *cli.Context) error {
				logger.Info("config model ...")
				config.WriteConfig(store, db, api, option, localpath, helper.GetDefaultConfigPath())
				return nil
			},
		},
		{
			Name:    "exec",
			Aliases: []string{"exec"},
			Usage:   "write config file & upload",
			Action: func(c *cli.Context) error {
				logger.Info("config & exec model...")
				config.WriteConfig(store, db, api, option, localpath, helper.GetDefaultConfigPath())

				Perform(helper.GetDefaultConfigPath(), modelname)

				return nil
			},
		},
		{
			Name:    "check",
			Aliases: []string{"check"},
			Usage:   "check data",
			Action: func(c *cli.Context) error {
				logger.Info("check data model...")
				config.WriteConfig(store, db, api, option, localpath, helper.GetDefaultConfigPath())

				Check(helper.GetDefaultConfigPath(), modelname)

				return nil
			},
		},
		{
			Name:    "load",
			Aliases: []string{"load"},
			Usage:   "load config file & exec",
			Action: func(c *cli.Context) error {
				if configfile != "" {
					logger.Info("load config file %s", configfile)
					config.LoadConfig(configfile)

					Perform(configfile, modelname)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)

}
