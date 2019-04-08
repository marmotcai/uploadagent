package helper

import (
	"log"
	"os"
	"time"
)

func GetFileModTime(path string) int64 {
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

	return fi.ModTime().Unix() + time.Now().Unix()
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
