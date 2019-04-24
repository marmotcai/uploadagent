package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/marmotcai/uploadagent/config"
	"github.com/marmotcai/uploadagent/database"
	"github.com/marmotcai/uploadagent/helper"
	"github.com/marmotcai/uploadagent/logger"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
	"strconv"
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
type StorageContext interface {
	open() error
	close()
	uploadfile(fileKey, filepath, remotepath string) (string, error)
	check(fileKey string) (error)
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

func uploadinfofile(ctx StorageContext, model config.ModelConfig, filekey, filepath, remotepath, remoteurl string, complete UploadComplete) (error) {
	jsonfilekey := filekey + ".info"
	jsonfile := config.TempPath + "/" + jsonfilekey
	data, err := helper.GetJsondata(filekey, filepath, remoteurl, jsonfile)
	if (err != nil) {
		logger.Info("get mediainfo json error : %s", err)
	}
	if (helper.IsExistsPath(jsonfile)) {
		_, err := ctx.uploadfile(jsonfilekey, jsonfile, remotepath)
		if (err != nil) {
			logger.Info("upload %s error : $s\n", jsonfile, err)
		}
	}
	if (complete != nil) {
		err = complete(model, data)
	}
	if (err != nil) {
		return fmt.Errorf("complete post %s error : $s\n", filekey, err)
	}
	return nil
}

func uploadfile(ctx StorageContext, model config.ModelConfig, filekey string, filepath string, remotepath string, complete UploadComplete) (error) {
	remoteurl, err := ctx.uploadfile(filekey, filepath, remotepath)
	if (err != nil) {
		return fmt.Errorf("upload %s error : $s\n", filekey, err)
	}
	err = uploadinfofile(ctx, model, filekey, filepath, remotepath, remoteurl, complete)

	return err
}

func uploaddir(ctx StorageContext, model config.ModelConfig, dir string, complete UploadComplete) (error) {
	logger.Info("upload ", dir)

	upload_error := 0
	upload_success := 0

	suffix_white := model.OptionWith.GetString("suffix_white")
	files, _ := helper.GetFilelist(dir, suffix_white)
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
			upload_error ++
			continue;
			// return fmt.Errorf("uploadfile %s error : $s\n", filekey, err)
		}
		upload_success ++
	}

	logger.Info("upload success: " + strconv.Itoa(upload_success))
	logger.Info("upload error: " + strconv.Itoa(upload_error))

	return nil
}

func checkfile(ctx StorageContext, model config.ModelConfig, filekey string, filepath string, remotepath string, complete UploadComplete) (error) {
	err := ctx.check(filekey)
	if (err != nil) {
		return fmt.Errorf("checkfile %s error : $s\n", filekey, err)
	}

	return nil
}

func Query(dbobj sql.DB, sql string) (*sql.Rows, error) {
	rows, err := dbobj.Query(sql)
	if (err != nil) {
		return nil, err
	}
	return rows, nil
}

func checkinfofile(dbobj *sql.DB, infofilepath string) (error) {

	f, err := ioutil.ReadFile(infofilepath)
	if err != nil {
		return fmt.Errorf("Load .info failed (%s)", err.Error())
	}

	var inter interface{}
	err = json.Unmarshal(f, &inter)
	if err != nil {
		return fmt.Errorf("json Unmarshal error (%s)", err.Error())
	}

	//要访问解码后的数据结构，需要先判断目标结构是否为预期的数据类型
	mi, ok := inter.(map[string]interface{})
	//然后通过for循环一一访问解码后的目标数据
	if !ok {
		return fmt.Errorf("mi can not get")
	}

	v := mi["general"].(map[string]interface{})
	filekey := v["filekey"]

	rows, err := dbobj.Query(fmt.Sprintf("SELECT id, code FROM cg_movie WHERE code = \"%s\"", filekey))
	defer rows.Close()

	if (err == nil) {
		for rows.Next() {
			var id string
			var code string
			err = rows.Scan(&id, &code)
			if (err == nil) {
				return nil
			}
		}
	}

	return fmt.Errorf("data not found in db")
}


func checkdir(ctx StorageContext, model config.ModelConfig, dir string, complete UploadComplete) (error) {
	logger.Info("check ", dir)

	var dbctx database.DBContext
	if len(model.Databases) != 0 {
		for _, dbCfg := range model.Databases {
			dbctx = database.GetDBModel(model, dbCfg)
			if dbctx != nil {
				err := dbctx.Perform()
				if (err != nil) {
					logger.Error("db open failed (%s)", err)
				}
				break
			}
		}
	}
	if (dbctx == nil) {
		return fmt.Errorf("db is not config, check failed")
	}

	suffix_white := model.OptionWith.GetString("suffix_white")
	infodblost_count := 0
	infofilelost_count := 0
	dbobj := dbctx.GetDBObj()
	files, _ := helper.GetFilelist(dir, suffix_white)
	for _, filepath := range files {
		logger.Info("--- check start ---")
		logger.Info(filepath)

		if (helper.IsTempfile(filepath)) {
			err := ctx.open()
			if (err != nil) {
				logger.Info("storage open failed (%s)", err)
			}
			continue
		}

		suffix := strings.ToLower(path.Ext(filepath))

		if 	(suffix != ".info") {
			infofilename := filepath + ".info"
			if (helper.IsExistsPath(infofilename)) {
				//如果找到已经保存的文件则查询是否有对应的info文件，同时校验info文件是否上报到了数据库

				err := checkinfofile(dbobj, infofilename)
				if (err == nil) {
					logger.Info("check info file ok: ", filepath)
				} else {
					logger.Error("check info file failed: ", err)
					infodblost_count ++
				}

			} else {
				if (helper.IsMediafile(filepath)) {
					logger.Error("info file not found! the media file: ", filepath)
					infofilelost_count ++
/*
					prefix_path := model.OptionWith.GetString("prefix_path")

					uploadinfofile(ctx, model, path.Base(filepath), filepath, path.Dir(filepath),
						strings.Replace(filepath, prefix_path, "", -1), complete)*/
				} else {
					logger.Error("info file not found! is not media file: ", filepath)
				}
			}
		} else {
			//如果是info文件则反查对应的数据文件是否存在，一般来说不会出现不存在的情况
			mediafilepath := strings.Replace(filepath, suffix, "", -1)
			if (!helper.IsExistsPath(mediafilepath)) {
				logger.Error("media file not found! the info file: ", filepath)
			}
		}
		logger.Info("--- check end ---\n")
	}

	logger.Info("info db lost count: " + strconv.Itoa(infodblost_count))
	logger.Info("info file lost count: " + strconv.Itoa(infofilelost_count))

	return nil
}

// Run storage
func Run(model config.ModelConfig, archivePath string, complete UploadComplete) (err error) {
	logger.Info("------------- Storage (Run) --------------")

	base := newBase(model, archivePath)
	var ctx StorageContext
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

	logger.Info("------------- Storage (Run) --------------")
	return nil
}

func Check(model config.ModelConfig, archivePath string, complete UploadComplete) (err error) {
	logger.Info("------------- Storage (Check) --------------")

	base := newBase(model, archivePath)
	var ctx StorageContext
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
		logger.Error(err)
		return err
	}
	defer ctx.close()

	err = checkdir(ctx, model, model.StoreWith.Viper.GetString("path"), complete)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("------------- Storage (Check) --------------")
	return nil
}