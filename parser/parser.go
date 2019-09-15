package parser

import (
	"github.com/PuerkitoBio/goquery"
	"icourse/config"
	"icourse/download"
	"icourse/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

//此处放置部分页面特有大的处理函数

//不同下载选项的具体实现
func DownloadAll(id string){
	DownloadMost(id)
	DownloadShareResource(id)
}
func DownloadMost(id string){
	DownloadVideoPPT(id)
	DownloadAssignments(id)
	DownloadTestPaper(id)
}
func DownloadVideoPPT(id string){
	s:=config.VideoPPT+id
	files:=extractURLs(s)
	download.DownloadFiles(files)
}
func DownloadAssignments(id string){
	s:=config.Assignments+id
	files:=extractURLs(s)
	download.DownloadFiles(files)
}
func DownloadTestPaper(id string){
	s:=config.TestPaper+id
	files:=extractURLs(s)
	download.DownloadFiles(files)
}
func DownloadShareResource(id string){
	s :=config.ShareResource+id
	files:=extractOthers(s)
	download.DownloadFiles(files)
}
//根据sectionID，构造post请求，parentPath为前面的路径名，返回文件的数组
func getVideo(id string,parentPath string) []utils.File {
	var files []utils.File
	res, err := http.PostForm("https://www.icourses.cn/web//sword/portal/getRess",url.Values{"sectionId": {id}})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	//匹配模式,json在go中真的不太好处理呀
	//匹配url
	fileUrl:=utils.MatchAll(string(body),`"fullResUrl":"(.*?)"`)
	//匹配文件名
	title:=utils.MatchAll(string(body),`"title":"(.*?)"`)
	//匹配文件类型
	fileType:=utils.MatchAll(string(body),`"resMediaType":"(.*?)"`)
	for i:=0;i<len(fileUrl);i++{
		//为了生成正确的文件名真的不容易,json中ppt类型的文件实际为pdf类型
		files=append(files,utils.File{fileUrl[i][1],filepath.Join(parentPath,title[i][1]+"."+strings.Replace(fileType[i][1],"ppt","pdf",-1))})
	}
	return files
}

//下载有目录结构的页面，返回文件的数组
func extractURLs(url string)  []utils.File {
	// 加载html
	data:=utils.HttpGet(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	//需要下载的文件集合
	var files []utils.File
	// 定位每一章节
	doc.Find("#chapters >.panel").Each(func(i int, s *goquery.Selection) {
		//提取章节的名称
		parentPATH := utils.Format(s.Find(".chapter-title-text").Text())
        //如果找不到媒体文件
		if s.Find("a[data-class=media]").Nodes==nil{
			//如果得不到文件的路径
			//fmt.Println("find zero media")
			s.Find("a[data-secid]").Each(func(j int, a *goquery.Selection) {

				secid, _ :=a.Attr("data-secid")
				//fmt.Println(secid)
				secondPATH:=utils.Format(a.Text())
				files=append(files,(getVideo(secid,filepath.Join(parentPATH,secondPATH)))...)
			})
		}else {
			//如果能直接得到文件的路径
			s.Find("a[data-class=media]").Each(func(j int, a *goquery.Selection) {
				// 分别记录文件名称、下载URL、文件类型以及构造的文件路径
				fileTitle, _ := a.Attr("data-title")
				fileURL, _ := a.Attr("data-url")
				fileType, _ := a.Attr("data-type")
				filePATH := filepath.Join(parentPATH, fileTitle+"."+fileType)
				files = append(files, utils.File{fileURL, filePATH})
			})
		}


	})
	return files
}

//直接提取页面文件，下载到“其他文件”文件夹中，返回文件的数组
func extractOthers(url string) []utils.File{
	data:=utils.HttpGet(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	//需要下载的文件集合
	var files []utils.File
	doc.Find("#other-sources > ul > li > a").Each(func(i int, s *goquery.Selection) {
		fileTitle, _ :=s.Attr("data-title")
		fileType, _ :=s.Attr("data-type")
		fileURL, _ :=s.Attr("data-url")
		files=append(files,utils.File{fileURL,filepath.Join("其它文件",fileTitle+"."+fileType)})
	})
	return files
}
