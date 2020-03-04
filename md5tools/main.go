package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

/**
Calculate md5 of a file
*/
func calcMd5(filePath string) string {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	fileSize := info.Size()

	blocks := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		buf := make([]byte, blockSize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	md5string := hex.EncodeToString(hash.Sum(nil))

	return md5string
}

const (
	fileChunk = 8192 // 8KB

	Usage       = "MD5 tools"
	Name        = "MD5Tools"
	Version     = "1.0"
	Description = "Display the file/dir md5 value"
)

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

const MediaSuffix = ".go|.ts|.wmv|.asf|.mp3|.wav|.mpg|.mpeg|.avi|.mkv|.mov|.mp4|.rm|.rmvb|.m2ts|.vob|.mxf"

func IsMediafile(filename string) bool {
	suffix := strings.ToLower(path.Ext(filename))
	if strings.Contains(MediaSuffix, suffix) {
		return true
	}

	return false
}

func IsTempfile(filepath string) bool {
	filename := path.Base(filepath)
	first := filename[0:1]

	if first == "." {
		return true
	}
	return false
}

func GetFilelist(dirPth, suffix string) ([]string, error) {
	var files []string
	var suffixs []string
	if len(suffix) > 0 {
		suffixs = strings.Split(suffix, "|")
	}

	err := filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi == nil {
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}

		if IsTempfile(filename) {
			return nil
		}

		have := false
		for j := 0; j < len(suffixs); j++ {
			s := strings.ToUpper(suffixs[j])
			if strings.HasSuffix(strings.ToUpper(path.Ext(fi.Name())), s) {
				// if strings.Contains(strings.ToUpper(filename), s) {
				have = true
				break
			}
		}
		if (len(suffixs) <= 0) || (have) {
			files = append(files, filename)
		}
		return nil
	})

	return files, err
}

func getsrcfiles(src, suffix string) ([]string, error) {
	filelist := []string{}
	if IsDir(src) {
		filelist, _ = GetFilelist(src, suffix)
	} else {
		filelist = append(filelist, src)
	}
	return filelist, nil
}

func IsExistsPath(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func generate(src, suffix string) (int, error) {
	fmt.Println("generate files md5 from : ", src)

	filelist, _ := getsrcfiles(src, suffix)
	count := 0
	for _, file := range filelist {
		md5string := calcMd5(file)
		fmt.Println(md5string + "  " + file)

		md5file := file + ".md5"
		if IsExistsPath(md5file) {

			fmt.Println("md5 file is exists : " + file)
		} else {
			err := ioutil.WriteFile(file+".md5", []byte(md5string), 0644)
			if err != nil {
				panic(err)
			}

			count++
		}
	}
	return count, nil
}

func check(src, suffix string) (int, []string, error) {
	count := 0
	var err_files []string
	files, _ := getsrcfiles(src, suffix)
	for _, file := range files {
		if IsExistsPath(file + ".md5") {
			md5string_old := ""
			if contents, err := ioutil.ReadFile(file + ".md5"); err == nil {
				//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
				md5string_old = strings.Replace(string(contents), "\n", "", 1)
				// fmt.Println(md5string_old)
			}
			md5string := calcMd5(file)
			if md5string_old != md5string {
				err_files = append(err_files, file)
				fmt.Println("check md5 error : " + file + "( " + md5string_old + " => " + md5string + ")")
			} else {
				count++

				fmt.Println("check md5 ok : " + file)
			}
		} else {
			fmt.Println("md5 file is not exists : " + file)
		}
	}
	return count, err_files, nil
}

func clean(src string) (int, error) {
	count := 0
	filelist, _ := getsrcfiles(src, ".md5")
	for _, file := range filelist {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("del file：" + file + " fail.")
		} else {
			count++
			fmt.Println("del file：" + file + " succeed.")
		}
	}
	return count, nil
}

func Print_Color(text string) error {
	fmt.Println(text)
	// fmt.Printf("\n %c[1;30;33m%s%c[0m\n\n", 0x1B, text, 0x1B)
	return nil
}

func main() {
	src := "./"
	for {
		fmt.Println("-------------------")
		fmt.Println("MD5生成和验证工具")
		fmt.Println("生成 MD5文件请按: g")
		fmt.Println("校验 MD5文件请按: c")
		fmt.Println("清除 MD5文件请按: d")
		fmt.Println("按其它键退出")
		fmt.Println("-------------------")

		fmt.Print("-> ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		fmt.Println("你输入的是：", input.Text())
		text := input.Text()
		if strings.Compare("g", text) == 0 {
			fmt.Println("会花费一定时间，请稍等...")
			count, _ := generate(src, MediaSuffix)
			Print_Color("** 生成md5已完成, 共生成： " + strconv.Itoa(count) + " 个文件 **")

		} else if strings.Compare("c", text) == 0 {
			fmt.Println("会花费一定时间，请稍等...")
			count, err_files, _ := check(src, MediaSuffix)
			Print_Color("** 校验md5已完成, 校验正确： " + strconv.Itoa(count) + " 个文件， 校验错误： " + strconv.Itoa(len(err_files)) + " 个文件 **")
			if len(err_files) > 0 {
				fmt.Println("校验出现错误的文件如下:")
				for _, err_file := range err_files {
					fmt.Println(err_file)
				}
			}
		} else if strings.Compare("d", text) == 0 {
			count, _ := clean(src)
			Print_Color("** 清除md5已完成, 共清除 " + strconv.Itoa(count) + " 个文件 **")
		} else {
			os.Exit(0)
		}
	}
	/*
		var src string

		app := cli.NewApp()
		app.Name = Name
		app.Usage = Usage
		app.Version = Version
		app.Description = Description

		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name: "file, s", Usage: "file", Destination: &src,
			},
		}
		app.Commands = []cli.Command {
			{
				Name:    "generate",
				Aliases: []string{"generate"},
				Usage:   "generate file md5",
				Action: func(c *cli.Context) error {
					return generate(src)
				},
			},

			{
				Name:    "check",
				Aliases: []string{"check"},
				Usage:   "check file md5",
				Action: func(c *cli.Context) error {
					return check(src)
				},
			},

			{
				Name:    "clean",
				Aliases: []string{"clean"},
				Usage:   "clean md5 file",
				Action: func(c *cli.Context) error {
					return clean(src)
				},
			},
		}

		app.Run(os.Args)
	*/
}
