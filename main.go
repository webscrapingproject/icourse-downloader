package main

import (
	"icourse/config"
	"icourse/parser"
	"icourse/utils"
	"flag"
	"fmt"
	"os"
	"strings"
)



//初始化相关参数
func init() {
	flag.BoolVar(&config.Version, "v", false, "Show version")
	//all为全部下载，most为视频课件以及试卷，也为下载默认选项，videoPPT仅下载视频和课件，exams为仅下载试卷，resources仅下载其它资源
	flag.StringVar(&config.ContentOptions, "co", "all", "Only for icourse : Specify the download content {all,most,videoPPT,assignments,testPaper,shareResource}\nOnly for chinesemooc : Specify the download content {all, video , PPT}\nOnly for Icourse163 : Specify the download content {all, video , PPT , RichText}\nOnly for Datacamp : all the content links will be extracted")
	//华文慕课的下载选项，只有三个:全部下载（默认），只下载视频 以及 只下载课件
	//中国大学mooc的下载选项：全部下载（默认），只下载视频，只下载课件以及只下载富文本文件
	//设置下载路径
	flag.StringVar(&config.OutputPath, "o", "", "Specify the output path")
	//设置cookie
	flag.StringVar(&config.Cookie, "c", "", "Cookie or the path of Cookie file")
}

func download(url string) bool {
	domain := utils.Domain(url)
	switch domain {
	case "icourses":
		parser.DownloadIcourse(url,config.ContentOptions)
	case "chinesemooc":
		//fmt.Println("SUCCESS")
		parser.DownloadChinesemooc(url,config.ContentOptions,config.Cookie)
	case "icourse163":
		parser.DownloadIcourse163(url,config.ContentOptions,config.Cookie)
	case "datacamp":
		parser.DownloadDatacamp(url,config.ContentOptions)
	default:
		fmt.Println("The website is not supported now ")
		return false
	}
	return true
}

func main() {
	//此处参考了annie的代码
	//fmt.Println(parser.GetDCStartURLs("https://campus.datacamp.com/courses/data-visualization-with-ggplot2-2"))

	flag.Parse()
	args := flag.Args()
	if config.Version {
		utils.PrintVersion()
		return
	}
	if len(args) < 1 {
		fmt.Println("Too few arguments")
		fmt.Println("Usage: icourse [args] URLs...")
		flag.PrintDefaults()
		return
	}
	if config.Cookie != ""{
		utils.ReadCookieFromFile(config.Cookie)
		//fmt.Println(config.Cookie)
	}
	var isErr bool
	//可以下载多个url
	for _, videoURL := range args {
		if err := download(strings.TrimSpace(videoURL)); err {
			isErr = true
		}
	}
	if isErr {
		os.Exit(1)
	}

}

