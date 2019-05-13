package helper

import (
	"encoding/json"
	"fmt"
	"github.com/marmotcai/uploadagent/logger"
	"log"
	"os"
	"path"
	"time"
)

func GetEigenvalue(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()// + time.Now().Unix()
}

func GetMediainfo(filepath string) (string, error) {

	var opts []string
	opts = append(opts, filepath)
	opts = append(opts, "--Inform=file://mediaformat")
	ms, err := Exec("mediainfo", opts...)
	if err != nil {
		return "", err
	}

	return ms, nil
}

func GetJsondata(key, localfile, url, jsonfile string, level int) ([]byte, error) {
	if (!IsMediafile(localfile)) {
		return nil, nil
	}

	ms, err := GetMediainfo(localfile)
	//先创建一个目标类型的实例对象，用于存放解码后的值
	var inter interface{}
	err = json.Unmarshal([]byte(ms), &inter)
	if err != nil {
		return nil, fmt.Errorf("error in translating,", err.Error())
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
						name, err := GetNameFromPath(localfile, level)
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

	if (len(jsonfile) >= 0) {
		filePtr, err := os.Create(jsonfile)
		if err!=nil{
			logger.Info("write file media json file error : %v", err)
		}
		defer filePtr.Close()

		enc := json.NewEncoder(filePtr)
		err = enc.Encode(inter)
		if (err != nil) {
			logger.Info("write file media json file error : %v", err)
		}
	}

	return data, nil
}