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
	flag.StringVar(&config.ContentOptions, "co", "all", "Only for icourse : Specify the download content {all,most,videoPPT,assignments,testPaper,shareResource}\nOnly for chinesemooc : Specify the download content {all, video , PPT}")
	//华文慕课的下载选项，只有三个:全部下载，只下载视频 以及 只下载课件

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
	}
	return true
}

func main() {
	//此处参考了annie的代码
	//fmt.Println(parser.GetStartURLs("http://www.chinesemooc.org/mooc/4880","pku_auth=161evS%2BQJtmq%2FGJRyU%2BFhfaNLyG88SrUPqUX5a0eOUW49JVtBaPxY7lt1vp2MvvcC9UaH8qYx3%2B0cSja0MeVNCmDSWRQ; pku_loginuser=univeroner%40gmail.com; pku_reward_log=daylogin%2C1173273; Hm_lvt_ff4f6e9862a4e0e16fd1f5a7f6f8953b=1569321857,1569494843,1569494850,1569759380; PHPSESSID=p72d5gqftbmp65mmr2n9ghrah5; pku__refer=%252Fmooc%252F4880; Hm_lpvt_ff4f6e9862a4e0e16fd1f5a7f6f8953b=1569761588"))

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

