package main

import (
	"flag"
	"fmt"
	"icourse/config"
	"icourse/parser"
	"icourse/utils"
	"os"
	"strings"
)



//初始化相关参数
func init() {
	flag.BoolVar(&config.Version, "v", false, "Show version")
	//all为全部下载，most为视频课件以及试卷，也为下载默认选项，videoPPT仅下载视频和课件，exams为仅下载试卷，resources仅下载其它资源
	flag.StringVar(&config.ContentOptions, "c", "most", "Specify the download content {all,most,videoPPT,assignments,testPaper,shareResource}")
	flag.StringVar(&config.OutputPath, "o", "", "Specify the output path")
	//flag.StringVar(&config.StartUrl, "F", "", "course URL")
}
func download(url string,options string) bool{
    id:=utils.MatchAll(url,`course_([0-9]*)`)
	if id != nil{
		//得到课程的id地址
		idNum:=id[0][1]
		//fmt.Println(idNum)
		switch options{
		case "all":
			parser.DownloadAll(idNum)
		case "most":
			parser.DownloadMost(idNum)
		case "videoPPT":
			parser.DownloadVideoPPT(idNum)
		case "assignments":
			parser.DownloadAssignments(idNum)
		case "testPaper":
			parser.DownloadTestPaper(idNum)
		case "shareResource":
			parser.DownloadShareResource(idNum)
			}
	} else{
		//网址不符合格式
		fmt.Printf("this website %s is not supported now",url)
		return true
	}
 return true
}

func main() {
	//此处参考了annie的代码

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
	var isErr bool
	//可以下载多个url
	for _, videoURL := range args {
		if err := download(strings.TrimSpace(videoURL),config.ContentOptions); err {
			isErr = true
		}
	}
	if isErr {
		os.Exit(1)
	}

}

