package helper

import (
	"flag"
	"fmt"
	"github.com/yanyiwu/gosimhash"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	Version = "1.0"

	Usage = "CG Media Upload Agent"

	App_name = "UploadAgent"
	App_config = "ua-default"
)

var (
	// IsGnuTar show tar type
	IsGnuTar = false
)

var top_n = flag.Int("top_n", 6, "")
var sher gosimhash.Simhasher

func init() {
	checkIsGnuTar()

	sher = gosimhash.New("./dict/jieba.dict.utf8",
						"./dict/hmm_model.utf8",
						"./dict/idf.utf8",
					"./dict/stop_words.utf8")
	// defer sher.Free()
}

func checkIsGnuTar() {
	out, _ := Exec("tar", "--version")
	IsGnuTar = strings.Contains(out, "GNU")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetHomePath() string {
	return os.Getenv("HOME")
}

func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func GetDefaultConfigPath() string {
	return fmt.Sprintf("%s/.%s.yml", GetCurrentPath(), App_config)
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func IsMediafile(filename string) bool {
	suffix := strings.ToLower(path.Ext(filename))

	if 	(suffix == ".ts") ||
		(suffix == ".mpg") ||
		(suffix == ".mpeg") ||
		(suffix == ".avi") ||
		(suffix == ".mkv") ||
		(suffix == ".mov") ||
		(suffix == ".mp4") {
		return true
	}
	return false
}

func IsTempfile(filename string) bool {
	suffix := path.Base(filename)[0: 1]

	if 	(suffix == ".") {
		return true
	}
	return false
}

func GetFilelist(dirPth, suffix string) ([]string, error) {
	var files []string
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err := filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi == nil {
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}

		if ((suffix == "") ||
			(strings.HasSuffix(strings.ToUpper(fi.Name()), suffix))) {

			files = append(files, filename)
		}
		return nil
	})

	return files, err
}

// Expected to be equal.
func equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
	}
}

// Expected to be unequal.
func unequal(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
	}
}

func SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}
	return sim
}

func GetNameFromPath(url string) (string, error) {
	parentdir := ""

	suffix := strings.ToLower(path.Ext(path.Base(url)))
	basefilename := ""

	if (len(suffix) <= 0) {
		parentdir = url
	} else {
		parentdir = path.Dir(url)
		basefilename = path.Base(url)
	}

	counter := 0
	err := filepath.Walk(parentdir, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi == nil {
			return err
		}

		depth := strings.Count(filename,"/") - strings.Count(parentdir,"/")
		if depth > 1 {
			return filepath.SkipDir
		}
		if (filename == parentdir) {
			return nil
		}

		if fi.IsDir() { // 如果有子目录则返回
			return fmt.Errorf("ERROR:have a subdire，non-deterministic")
		}
		if (fi.Name()[0: 1] == ".") {
			return nil
		}
		if (len(basefilename) <= 0) {
			basefilename = fi.Name()
			return  nil
		}

		var percent float64
		SimilarText(basefilename, fi.Name(), &percent)

		if (percent > 60) {
			counter ++
			if (counter > 3) {
				dirname := path.Base(parentdir)
				name := strings.TrimSuffix(basefilename, suffix)

				/*SimilarText(dirname, name, &percent)
				if (percent > 60) {*/
				if (strings.Contains(name, dirname)) {
					return fmt.Errorf("OK:" + name)
				}

				return fmt.Errorf("OK:" + dirname + "-" + name)
			}
		}

		return nil})

	if (err != nil) {
		msg := err.Error()
		msgs := strings.Split(msg, ":")
		for j := 0; j < len(msgs); j++ {
			switch strings.ToUpper(msgs[j]) {
			case "OK" :{
				return msgs[j + 1], nil //这是电视剧，取目录名作为名称
			}
			case "ERROR" :{
				return "", err
			}
			}
		}
	}
	name := strings.TrimSuffix(basefilename, suffix)
	return name, err //可能是电影
}

func GetFileKey(filepath, formatstr string) (string, string) {
	filekey := fmt.Sprintf("%x",
		sher.MakeSimhash(strconv.FormatInt(GetEigenvalue(filepath),10) + "_" + path.Base(filepath), *top_n))
	filekey = filekey + path.Ext(filepath)

	destpath := ""
	params := strings.Split(formatstr, "%")
	for j := 0; j < len(params); j++ {
		switch strings.ToUpper(params[j]) {
		case "CLASS_LAST0", "CLASS_LAST1", "CLASS_LAST2", "CLASS_LAST3", "CLASS_LAST4", "CLASS_LAST5": {
			param := params[j]
			lastdir := ""

			PathLevels, err := strconv.Atoi(param[len(param)-1:])
			if (PathLevels > 0) && (err == nil) {
				pdir := path.Dir(filepath)
				for a := 0; a < PathLevels; a++ {
					dirname := path.Base(pdir)
					pdir = path.Dir(pdir)

					lastdir = path.Join(dirname, lastdir)
				}
			}
			if (len(lastdir) > 0) {
				destpath = path.Join(destpath, lastdir)
			}
		}
		case "HASH_TOP0", "HASH_TOP1", "HASH_TOP2", "HASH_TOP3", "HASH_TOP4", "HASH_TOP5": {
			param := params[j]
			PathLevels, err := strconv.Atoi(param[len(param)-1:])
			if (PathLevels > 0) && (err == nil) {
				for a := 0; a < PathLevels; a++ {
					destpath = path.Join(destpath, filekey[a:a+1])
				}
			}
		}
		case "CURTIME": {
			currentTimeData:=time.Now().Format("2006-01-02-15")

			destpath = path.Join(destpath, currentTimeData)
		}
		case "HASHFULL": {
			// destpath = path.Join(destpath, filekey)
		}
		default: {
			if (len(params[j]) > 0) {
				destpath = path.Join(destpath, params[j])
			}
		}

		}
	}

	return destpath, filekey
}
