package utils

import (
	"fmt"
	"icourse/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

//定义结构体,数据包括文件的下载路径以及保存路径
type File struct {
	FileURL string
	FilePATH string
}

//打印版本
// PrintVersion print version information
func PrintVersion() {

	fmt.Printf(
		"\n%s: version %s\n",
		"icourse",
		config.VERSION,
	)
}
//处理错误情况
func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
//判断文件是否存在
func FileExists(filePath string) bool{
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else{
		return false
	}
}

//去除字符串里的空格换行等
func Format(str string)string{
	re := regexp.MustCompile("[\r\n\t]")
	res := re.ReplaceAllString(str, "")
	return res
}

// 符合某一个正则表达式
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// 返回所有的符合结果
func MatchAll(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	value := re.FindAllStringSubmatch(text, -1)
	return value
}

//根据url,构造get请求
func HttpGet(s string) string {
	res, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	return string(body)
}