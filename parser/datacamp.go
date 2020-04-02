package parser

import (
	"fmt"
	"icourse/utils"
	"strconv"
)

//// Datacamp下载部分
func DownloadDatacamp(url string,options string){
	courseName := GetDCName(url)
	fmt.Println(courseName)
	ExtractDCVideo(url)
	//switch options{
	//case "all":
	//	DownloadDCAll(url,courseName)
	//case "video":
	//	DownloadDCVideos(url,courseName)
	//case "subtitle":
	//	DownloadDCsubtitles(url,courseName)
	//}

}


//得到课程名称
func GetDCName(url string)string {
	return utils.MatchAll(url,`courses/(.*)`)[0][1]
}

//从起始地址得到带视频章节的下载页面,返回视频列表和pdf课件列表
func GetDCStartURLs(starturl string) ([]string,[]string) {
	//从起始地址中提取课程名称
	courseName := GetDCName(starturl)
	//拼接出正确的地址
	url := "https://campus.datacamp.com/courses/"+courseName
	data:=utils.HttpGet(url)
	//fmt.Print(url)
	var urlList []string
	var pdfList []string
	content := utils.MatchAll(data,`VideoExercise.*?(https://campus.datacamp.com/courses/.*?)&quot;`)
	ppts := utils.MatchAll(data,`slides_link.*?(https.*?pdf)`)
	//fmt.Print(ppts)
	for i,item := range(content){
		if(i< len(content)/2) {
			urlList = append(urlList, item[1])
		}
	}
	//pdfList 有重复的
	for i,item := range(ppts){
		if(i< len(content)/2) {
			pdfList = append(pdfList, item[1])
		}
	}
	return urlList,pdfList
}

//提取单个视频的projector_key
func GetPK(url string)(string){
	data:=utils.HttpGet(url)
	PK := utils.MatchAll(data,`(https://projector.datacamp.com/\?projector_key=.*?)&quot;`)[0][1]
	return PK
}

//根据project key 返回视频和字幕的下载地址
func GetDCVideo(url string)(string,string){
	data:=utils.HttpGet(url)
	videoUrl:= "https:"+utils.MatchAll(data,`video_mp4_link.*?(//.*?\.mp4)`)[0][1]
	subtitleUrl := utils.MatchAll(data,`subtitle_vtt_link.*?(https.*?vtt)`)[0][1]
	return videoUrl,subtitleUrl
}

////返回需要下载的pdf文件地址
//func ExtractDCPPTs(url string)([]utils.File){
//	_,pdfList := GetDCStartURLs(url)
//	var pdfs []utils.File
//	//for _,item :=range(pdfList){
//	//	fileName := utils.MatchAll(item,`.*/(.*?pdf)`)[0][1]
//	//	pdfs = append(pdfs, utils.File{item,fileName})
//	//}
//	utils.WriteFile("pdflist.txt",pdfList)
//	return pdfs
//}

//返回待下载的视频文件列表和字幕文件列表
func ExtractDCVideo(url string)([]utils.File,[]utils.File){
	//需要下载的文件集合
	var videos []utils.File
	var subtitles []utils.File
	urlList,pdfList := GetDCStartURLs(url)
	//放置下载信息的列表
	var  downloadList []string
	var  renameList []string
	//fmt.Println(urlList)
	//对于每一个视频
	for i,item := range(urlList){
		PK := GetPK(item)
		//添加下载链接
		videoUrl,subtitleUrl := GetDCVideo(PK)
		downloadList = append(downloadList,videoUrl )
		downloadList = append(downloadList,subtitleUrl )
		fileName := utils.MatchAll(videoUrl,`.*/(.*?).mp4`)[0][1]
		newName := strconv.Itoa(i+1)
		subtitleName := utils.MatchAll(subtitleUrl,`.*/(.*?.vtt)`)[0][1]
		//添加重命名命令
		renameCM := "ren "+subtitleName+" "+newName+".vtt"
		renameVideo := "ren "+fileName+" "+newName+".mp4"
		fmt.Println(renameCM)
		renameList = append(renameList,renameCM)
		renameList = append(renameList,renameVideo)
		videos = append(videos, utils.File{videoUrl,newName+".mp4"})
		subtitles = append(subtitles, utils.File{subtitleUrl,newName+".vtt"})
	}
	utils.WriteFile("downloadList.txt",downloadList)
	utils.WriteFile("rename.bat",renameList)
	utils.WriteFile("pdflist.txt",pdfList)
	return videos,subtitles
}
