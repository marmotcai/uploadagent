package config

import (
	"fmt"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var (
	// Exist Is config file exist
	Exist bool
	// Models configs
	Models []ModelConfig
	// IsTest env
	IsTest bool
	// HomeDir of user
	HomeDir    string
	TempPath   string
	ConfigFile string
)

// ModelConfig for special case
type ModelConfig struct {
	Name        string
	DumpPath    string
	PackWith    SubConfig
	EncryptWith SubConfig
	StoreWith   SubConfig
	Archive     *viper.Viper
	Databases   []SubConfig
	Storages    []SubConfig
	Api		    []SubConfig
	Viper       *viper.Viper
}

// SubConfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

type Model_Store struct {
	Type              	string `yaml:"type"`
	Path				string `yaml:"path"`
	Url					string `yaml:"url"`
	User				string `yaml:"user"`
	Region            	string `yaml:"region"`
	Bucket            	string `yaml:"bucket"`
	ForcePath_style		string `yaml:"forcepath_style"`
	Timeout				string `yaml:"timeout"`
	IsMove				string `yaml:"ismove"`
	FileKeyFormat		string `yaml:"filekeyformat"`
	NamingRule			string `yaml:"namingrule"`
}

type Model_DB struct{
	Type				string `yaml:"type"`
	Host				string `yaml:"host"`
	Port				string `yaml:"port"`
	Database			string `yaml:"database"`
	Username			string `yaml:"username"`
	Password			string `yaml:"password"`
}

type Model_API struct{
	Type				string `yaml:"type"`
	Url					string `yaml:"Url"`
}

type Model struct {
	Pack_with struct {
		Type              	string `yaml:"type"`
	}

	Store_with			Model_Store `yaml:"store_with"`

	Archive struct {
		Pack				string `yaml:"pack"`
		Includes 			string `yaml:""`
	}

	Databases struct {
		DB_with			Model_DB `yaml:"db"`
	}

	API struct {
		API_with		Model_API `yaml:"rest"`
	}
}

type DefaultConfig struct {
	Models struct {
		M Model `yaml:"default"`
	}
}

func NewModel_Store() *Model_Store {
	return &Model_Store{}
}

func NewModel_DB() *Model_DB {
	return &Model_DB{}
}

func NewModel_API() *Model_API {
	return &Model_API{}
}

func WriteConfig(store *Model_Store, db *Model_DB, api *Model_API, local, filepath string) error {
	TempPath = path.Join(os.TempDir(), helper.App_name)
	if (!helper.PathExists(TempPath)) {
		helper.MkdirP(TempPath)
	}
	dc := new(DefaultConfig)

	f, err := ioutil.ReadFile(filepath)
	if (f != nil) {
		err = yaml.Unmarshal(f, dc)
	}

	dc.Models.M.Pack_with.Type = "none"
	dc.Models.M.Archive.Pack = "false"

	if (local != "") {
		dc.Models.M.Archive.Includes = local
	}

	//Store
	switch store.Type {
	case "s3":
		{
			dc.Models.M.Store_with.Type = store.Type
			if (store.Url != "") {
				dc.Models.M.Store_with.Url = store.Url
			}
			if (store.Region != "") {
				dc.Models.M.Store_with.Region = store.Region
			}
			if (store.User != "") {
				dc.Models.M.Store_with.User = store.User
			}
			if (store.Bucket != "") {
				dc.Models.M.Store_with.Bucket = store.Bucket
			}
			if (store.Path != "") {
				dc.Models.M.Store_with.Path = store.Path
			}
			if (store.ForcePath_style != "") {
				dc.Models.M.Store_with.ForcePath_style = store.ForcePath_style
			}
			if (store.FileKeyFormat != "") {
				dc.Models.M.Store_with.FileKeyFormat = store.FileKeyFormat
			}
			if (store.NamingRule != "") {
				dc.Models.M.Store_with.NamingRule = store.NamingRule
			}
		}
	case "scp", "ftp":
		{
			dc.Models.M.Store_with.Type = store.Type
			if (store.Url != "") {
				dc.Models.M.Store_with.Url = store.Url
			}
			if (store.User != "") {
				dc.Models.M.Store_with.User = store.User
			}
			if (store.Path != "") {
				dc.Models.M.Store_with.Path = store.Path
			}
			if (store.Timeout != "") {
				dc.Models.M.Store_with.Timeout = store.Timeout
			}
			if (store.FileKeyFormat != "") {
				dc.Models.M.Store_with.FileKeyFormat = store.FileKeyFormat
			}
			if (store.NamingRule != "") {
				dc.Models.M.Store_with.NamingRule = store.NamingRule
			}
		}
	case "local":
		{
			dc.Models.M.Store_with.Type = store.Type
			dc.Models.M.Store_with.Url = ""
			dc.Models.M.Store_with.User = ""

			if (store.Path != "") {
				dc.Models.M.Store_with.Path = store.Path
			}
			if (store.IsMove != "") {
				dc.Models.M.Store_with.IsMove = store.IsMove
			}
			if (store.FileKeyFormat != "") {
				dc.Models.M.Store_with.FileKeyFormat = store.FileKeyFormat
			}
			if (store.NamingRule != "") {
				dc.Models.M.Store_with.NamingRule = store.NamingRule
			}
		}

	default:
	}

	logger.Info("load store config for " + dc.Models.M.Store_with.Type)

	switch db.Type {
	case "mysql", "postgresql", "redis":
		{
			//db
			dc.Models.M.Databases.DB_with.Type = db.Type

			if (db.Host != "") {
				dc.Models.M.Databases.DB_with.Host = db.Host
			}

			if (db.Port != "") {
				dc.Models.M.Databases.DB_with.Port = db.Port
			}
			if (db.Database != "") {
				dc.Models.M.Databases.DB_with.Database = db.Database
			}
			if (db.Username != "") {
				dc.Models.M.Databases.DB_with.Username = db.Username
			}
			if (db.Password != "") {
				dc.Models.M.Databases.DB_with.Password = db.Password
			}
		}
	default:
	}

	logger.Info("load db config for " + dc.Models.M.Databases.DB_with.Type)

	switch api.Type {
	case "rest": {
			dc.Models.M.API.API_with.Type = api.Type
			if (api.Url != "") {
				dc.Models.M.API.API_with.Url = api.Url
			}
		}
	default:
	}

	logger.Info("load api config for " + dc.Models.M.API.API_with.Type)

	d, err := yaml.Marshal(dc)

	err = ioutil.WriteFile(filepath, d, 0644)

	return err
}

// loadConfig from:
// - ./UploadAgent.yml
// - ~/.github.com/marmotcai/uploadagent/UploadAgent.yml
// - /etc/github.com/marmotcai/uploadagent/UploadAgent.yml
func init2() {
	// os.Setenv("UA_CONFIGFILE", "/Users/andrewcai/Desktop/SynologyDrive/MySpace/go/src/github.com/marmotcai/uploadagent/script/ua.yaml")
	ConfigFile := os.Getenv("UA_CONFIGFILE")

	if ConfigFile != "" {
		fmt.Printf("Configfile is %s\n", ConfigFile)
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.SetConfigType("yaml")

		IsTest = os.Getenv("GO_ENV") == "test"
		HomeDir = os.Getenv("HOME")
		TempPath = path.Join(os.TempDir(), helper.App_name)

		if IsTest {
			viper.SetConfigName(helper.App_name + "_test")
			HomeDir = "../"
		} else {
			viper.SetConfigName(helper.App_name)
		}

		// ./ua.yml
		viper.AddConfigPath(".")
		if IsTest {
			viper.AddConfigPath("../")
		} else {
			// ~/.ua/ua.yml
			viper.AddConfigPath(fmt.Sprintf("%s/.%s", HomeDir, helper.App_name)) // call multiple times to add many search paths
			// /etc/ua/ua.yml
			viper.AddConfigPath(fmt.Sprintf("/etc/%s/", helper.App_name)) // path to look for the config file in
		}
	}
	LoadConfig("")

	return
}

func LoadConfig(configfile string) {
	if configfile != "" {
		viper.SetConfigFile(configfile)
	}
	// viper.Reset()
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Load UploadAgent config faild", err)
		return
	}

	Exist = true
	Models = []ModelConfig{}
	for key := range viper.GetStringMap("models") {
		Models = append(Models, loadModel(key))
	}
}

func loadModel(key string) (model ModelConfig) {
	model.Name = key
	model.DumpPath = path.Join(TempPath, fmt.Sprintf("%d", time.Now().UnixNano()), key)
	model.Viper = viper.Sub("models." + key)

	model.PackWith = SubConfig{
		Type:  model.Viper.GetString("pack_with.type"),
		Viper: model.Viper.Sub("pack_with"),
	}

	model.EncryptWith = SubConfig{
		Type:  model.Viper.GetString("encrypt_with.type"),
		Viper: model.Viper.Sub("encrypt_with"),
	}

	model.StoreWith = SubConfig{
		Type:  model.Viper.GetString("store_with.type"),
		Viper: model.Viper.Sub("store_with"),
	}

	model.Archive = model.Viper.Sub("archive")

	loadDatabasesConfig(&model)
	loadAPIConfig(&model)
	loadStoragesConfig(&model)

	return
}

func loadDatabasesConfig(model *ModelConfig) {
	subViper := model.Viper.Sub("databases")
	for key := range model.Viper.GetStringMap("databases") {
		dbViper := subViper.Sub(key)
		model.Databases = append(model.Databases, SubConfig{
			Name:  key,
			Type:  dbViper.GetString("type"),
			Viper: dbViper,
		})
	}
}

func loadAPIConfig(model *ModelConfig) {
	subViper := model.Viper.Sub("api")
	for key := range model.Viper.GetStringMap("api") {
		apiViper := subViper.Sub(key)
		model.Api = append(model.Api, SubConfig{
			Name:  key,
			Type:  apiViper.GetString("type"),
			Viper: apiViper,
		})
	}
}

func loadStoragesConfig(model *ModelConfig) {
	subViper := model.Viper.Sub("storages")
	for key := range model.Viper.GetStringMap("storages") {
		dbViper := subViper.Sub(key)
		model.Storages = append(model.Storages, SubConfig{
			Name:  key,
			Type:  dbViper.GetString("type"),
			Viper: dbViper,
		})
	}
}

// GetModelByName get model by name
func GetModelByName(name string) (model *ModelConfig) {
	for _, m := range Models {
		if m.Name == name {
			model = &m
			return
		}
	}
	return
}

// GetDatabaseByName get database config by name
func (model *ModelConfig) GetDatabaseByName(name string) (subConfig *SubConfig) {
	for _, m := range model.Databases {
		if m.Name == name {
			subConfig = &m
			return
		}
	}
	return
}
