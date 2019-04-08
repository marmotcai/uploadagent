package helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// MarshalJson 把对象以json格式放到response中
func MarshalJson(w http.ResponseWriter, v interface{}) (int, error) {
	data, err := json.Marshal(v)
	if (err != nil) {
		return -1, err
	}

	size, err := w.Write(data)

	return size, err
}

// UnMarshalJson 从request中取出对象
func UnMarshalJson(req *http.Request, v interface{}) (error) {
	result, err := ioutil.ReadAll(req.Body)
	if (err != nil) {
		return err
	}
	err = json.Unmarshal([]byte(bytes.NewBuffer(result).String()), v)
	return err
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
