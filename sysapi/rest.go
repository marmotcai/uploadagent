package sysapi

import (
	"bytes"
	"fmt"
	"github.com/marmotcai/uploadagent/logger"
	"net/http"
)

type Rest struct {
	Base
}

func (Rest) Post(url string, data []byte) error {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	logger.Info("post rest api success: ", resp.Status)
	/*
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
	*/
	return nil
}
