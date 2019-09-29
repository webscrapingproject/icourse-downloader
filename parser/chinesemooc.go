package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"icourse/download"
	"icourse/utils"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// 华文慕课下载部分
func DownloadChinesemooc(url string,options string,cookie string){
	courseName := GetChinesemoocName(url)
	fmt.Println(courseName)
	switch options{
	case "all":
		DownloadCMAll(url,courseName,cookie)
	case "video":
		DownloadCMVideos(url,courseName,cookie)
	case "PPT":
		DownloadCMPPTs(url,courseName,cookie)
	}

}

//下载视频和课件
func DownloadCMAll(url string,courseName,cookie string){
	DownloadCMVideos(url,courseName,cookie)
	DownloadCMPPTs(url,courseName,cookie)
}

//下载所有视频
func DownloadCMVideos(url string,courseName,cookie string){
	videoURL,_ := GetStartURLs(url,cookie)
	//fmt.Println(videoURL)
	files := ExtractCMVideo(videoURL,cookie)
	//fmt.Println(files)
	//可以直接下载
	download.DownloadCookieFiles(files,courseName,cookie)
}
//下载所有课件
func DownloadCMPPTs(url string,courseName,cookie string){
	_,PPTURL := GetStartURLs(url,cookie)
	//fmt.Println(PPTURL)
	files := ExtractCMPPTs(PPTURL,cookie)
	//使用cookie进行下载
	download.DownloadCookieFiles(files,courseName,cookie)
}


//得到课程名称
func GetChinesemoocName(url string)string {
	data := utils.HttpGet(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	courseName := utils.Format(doc.Find("head > title").Text())
	return courseName
}

func getVideoURL(url string,cookies string) string{
	data:=utils.HttpGetCookie(url,cookies)
	//设置长度，精确匹配
	// 可以改进为自动匹配最高清晰度,，此处默认标清
	//fmt.Println(data)
	SDvideo := utils.MatchAll(string(data),`http.{100,200}SD.mp4`)
	//
	videoUrl := ""

	//如果能选择清晰度的话选择高清
	if SDvideo != nil{
		videoUrl = SDvideo[0][0]
	} else { //只有一种清晰度的情况
		videoUrl = utils.MatchAll(string(data),"http.{100,200}.mp4")[0][0]
	}
	//fmt.Println(videoUrl)
	re := regexp.MustCompile(`\\`)
	videoUrl = re.ReplaceAllString(videoUrl,"")
	return videoUrl
}

//返回待下载的视频文件列表
func ExtractCMVideo(url string,cookies string) []utils.File{
	//需要下载的文件集合
	var files []utils.File
	//把cookie加在文件里
	data:=utils.HttpGetCookie(url,cookies)
	//fmt.Println(data)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	doc.Find("#coursefile > div.file-lists > div").Each(func(i int, s *goquery.Selection) {
		chapterName := s.Find(".main-item >span").Text()
		//fmt.Println(chapterName)

		s.Find("div.item-detail > ul > li.light.clearfix").Each(func(i int, p *goquery.Selection) {
			//小节的标题
			contentName := p.Find(".course-name >span").Text()
			//fmt.Println(contentName)
			//每一小节可能有多个视频
			p.Find(".icon-spow-wrap > .video").Each(func(j int, q *goquery.Selection) {

				filename,_ := q.Attr("original-title")
				str,_ := q.Attr("href")
				filePath := filepath.Join(chapterName,contentName,strconv.Itoa(j+1)+filename+".mp4")
				//提取两个id
				ID := utils.MatchAll(str,`&id=([0-9]*)`)[0][1]
				eid := utils.MatchAll(str,`eid=([0-9]*)`)[0][1]
				URL := "http://www.chinesemooc.org/api/course_video_watch.php?course_id="+ID+"&eid="+eid
				videoURL := getVideoURL(URL,cookies)
				//fmt.Println(filePath)
				files = append(files, utils.File{videoURL,filePath})
			})

		})
	})
	return files
}

//返回待下载的课件文件列表
func ExtractCMPPTs(url string,cookies string) []utils.File{
	//需要下载的文件集合
	var files []utils.File
	//把cookie加在文件里
	data:=utils.HttpGetCookie(url,cookies)
	//fmt.Println(data)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	doc.Find("#coursefile > ul > li").Each(func(i int, s *goquery.Selection) {
		chapterName := utils.Format(s.Find("div.title.clearfix").Text())
		//fmt.Println(chapterName)
		//每个章节底下的小节
		s.Find(".download-list").Each(func(i int, p *goquery.Selection){
			//提取小节名
			contentName := p.Find(".download-list-tit").Text() + p.Find("ul > li.download-list-num > a").Text()
			//fmt.Println(contentName)
			filePath := filepath.Join(chapterName,contentName);
			//fmt.Println(filePath)
			str,_ := p.Find("span[onclick].download-load").Attr("onclick");
			//fmt.Println(str)
			//提取pdf下载地址
			URL := utils.MatchAll(str,`window.open\("(.*)",`)[0][1]
			//如果是相对路径
			if(utils.MatchAll(URL,`http`) == nil){
				URL = "http://www.chinesemooc.org/" +URL
			}
			//baseURL := "http://www.chinesemooc.org/"
			//fmt.Println(URL)
			files = append(files,utils.File{URL,filePath})
		})
	})
	//注意下载的时候需要cookie
	return files
}

//从起始地址得到视频和课件的下载页面
func GetStartURLs(url string, cookies string) (string,string) {
	data:=utils.HttpGetCookie(url,cookies)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	s := doc.Find("#top-select > div > div > button")
	str,_ := s.Attr("onclick")
	//fmt.Println(str)
	//匹配classesid
	//匹配courseid
	classID := utils.MatchAll(str,`([0-9]*),`)[0][1]
	courseID := utils.MatchAll(str,`,([0-9]*)`)[0][1]
	//视频页面url
	courseProgress := "http://www.chinesemooc.org/kvideo.php?do=course_progress&kvideoid=" + classID + "&classesid=" + courseID
	//fmt.Println(courseProgress)
	//课件页面url
	courseCware :="http://www.chinesemooc.org/kvideo.php?do=course_cware_list&kvideoid=" + classID + "&classesid=" + courseID
	return courseProgress,courseCware
}

