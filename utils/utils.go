package utils

import (
	"fmt"
	"icourse/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)
//将字符串写入某个文件
func WriteFile(fileName string,content []string){
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	for _, v := range content {
		_, err:=fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}


//unicode 转换为utf8
func Unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}

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
	re := regexp.MustCompile("[\r\n\t ]")
	res := re.ReplaceAllString(str, "")
	return res
}

//替换
func Replace(str string,pattern string,new string)string{
	re := regexp.MustCompile(pattern)
	res := re.ReplaceAllString(str, new)
	return res
}


// 符合某一个正则表达式
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
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

//获取url链接的域名
// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return "Universal"
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

//根据url以及cookie构造get请求
func HttpGetCookie(s string,cookie string) string {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", s, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqest.Header.Set("Cookie",cookie)
	res, _ := client.Do(reqest)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	return string(body)
}

//根据url以及cookie构造post请求
func HttpPostCookie(s string,cookie string,postForm url.Values) string {
	client := &http.Client{}
	req, err := http.NewRequest("POST", s, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie",cookie)
	req.Body = ioutil.NopCloser(strings.NewReader(postForm .Encode()))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	return string(body)
}

//输入post参数以及网址，返回post结果
func HttpPostForm(postUrl string,postForm url.Values)string{
	client := &http.Client{}
	req, err := http.NewRequest("POST",postUrl,nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Body = ioutil.NopCloser(strings.NewReader(postForm .Encode()))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	return string(body)
}


//从文件中读取cookie
func ReadCookieFromFile(filePath string){
	//如果cookie是一个文件并且存在
	if _, fileErr := os.Stat(config.Cookie); fileErr == nil {
		// Cookie is a file
		data, err := ioutil.ReadFile(config.Cookie)
		Check(err)
		config.Cookie = string(data)
	}
}

//输入视频文件列表，返回ffmpeg下载命令的列表
func M3u8DownloadList(files []File,courseName string) {
	var  downloadList []string
	for _,file:=range(files){
		//检查目录是否存在
		filePath := filepath.Join(config.OutputPath,courseName,file.FilePATH)
		if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
			//建立目录
			_ = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		}
		command := fmt.Sprintf("ffmpeg -i '%s' -codec copy '%s'",file.FileURL,file.FilePATH)
		downloadList = append(downloadList,command)
	}
	WriteFile(filepath.Join(config.OutputPath,courseName,"downloadlist.txt"),downloadList)
}