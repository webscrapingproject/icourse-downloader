package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"icourse/config"
	"icourse/download"
	"icourse/utils"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Icourse163File struct{
	contentID string
	typeID string
	ID string
}
// 从初始网址中提取出总的文件信息
func HttpPostTermID(ID string)string {
	requestBody := url.Values{
		"callCount":{"1"},
		"scriptSessionId": {"${scriptSessionId}190"},
		"httpSessionId":{"c662f9cfbe7241b09927ff837c5d2ddc"},
		"c0-scriptName":{"CourseBean"},
		"c0-methodName":{"getMocTermDto"},
		"c0-id":{"0"},
		"c0-param0": {"number:"+ID},
		"c0-param1": {"number:0"},
		"c0-param2": {"boolean:true"},
		"batchId": {strconv.Itoa(int(time.Now().Unix()* 1000))},
	}
	return utils.HttpPostForm(config.GetMocTermDto,requestBody)
}

//从总的文件信息中提取出单个文件,返回三种文件的list
func ExtractIcourse163Res (text string,cookie string)([]utils.File,[]utils.File,[]utils.File){
	//三种不同类型的资源
	videoList := []utils.File{}
	pdfList := []utils.File{}
	richTextList := []utils.File{}
	//每一个章节,章节名为chapter[2]
	chapters := utils.MatchAll(text,"homeworks=[a-z0-9]*;.*id=([0-9]*).+name=\"(.*)\";")
	for _,chapter := range chapters{
		//每一个小节, 小节名为lesson[2]
		lessons := utils.MatchAll(text,"chapterId=" + chapter[1] + ".*contentType=1.*id=([0-9]*).*name=\"(.*)\".*test")
		for _,lesson := range lessons{
			//每个小节的视频，可能有多个,可能没有
			videos := utils.MatchAll(text, "contentId=([0-9]*).+contentType=(1).*id=([0-9]*).+lessonId=" +
				lesson[1] + ".*name=\"(.*)\"")
			for num,video:= range videos{
				videoURL:= parseIcourse163Video(video[3],cookie,1)
				videoPath := utils.Unicode2utf8(filepath.Join(chapter[2],lesson[2],strconv.Itoa(num+1)+ "-" + video[4] + ".mp4")) //添加文件编号以及后缀
				videoList = append(videoList,utils.File{FilePATH:videoPath,FileURL:videoURL})
			}
			//每个小节的pdf，可能多个,也可能没有
			pdfs := utils.MatchAll(text,"contentId=([0-9]*).*contentType=(3).*id=([0-9]*).*lessonId=" + lesson[1] + ".*name=\"(.*)\"")
			for num,pdf := range(pdfs){
				pdfURL := parseIcourse163PDF(Icourse163File{contentID:pdf[1],typeID:pdf[2],ID:pdf[3]})
				pdfPath := utils.Unicode2utf8(filepath.Join(chapter[2],lesson[2],strconv.Itoa(num +1 )+ "-"+pdf[4] + ".pdf"))
				pdfList = append(pdfList,utils.File{FilePATH:pdfPath,FileURL:pdfURL})
			}
			//每个小节的富文本文件，可能多个，也可能没有
			richTexts := utils.MatchAll(text,"contentId=([0-9]*).*contentType=(4).*id=([0-9]*).*jsonContent=(.*?);.*lessonId=" + lesson[1] + ".*name=\"(.*)\"")
			for _,richText := range richTexts{
				//直接可以构造附件的下载地址
				//防止只有富文本没有文件的情况
				if (richText[4] != "null") {
					//fmt.Println(richTexts)
					richTextUrl, fileName := parseIcourse163RichText(richText[4])
					richTextPath := utils.Unicode2utf8(filepath.Join(chapter[2], lesson[2], fileName))
					richTextList = append(richTextList, utils.File{FileURL: richTextUrl, FilePATH: richTextPath})
				}
			}

		}
	}
	return videoList,pdfList,richTextList
}

//从单个pdf文件的id等信息提取其下载地址（需要构造post请求）
func parseIcourse163File(content Icourse163File)string{
	requestBody := url.Values{
		"callCount":{"1"},
		"scriptSessionId": {"${scriptSessionId}190"},
		"httpSessionId":{"c47c239c3d414133b309532a7a8c2783"},
		"c0-scriptName":{"CourseBean"},
		"c0-methodName":{"getLessonUnitLearnVo"},
		"c0-id":{"0"},
		"c0-param0": {"number:"+content.contentID},
		"c0-param1": {"number:"+content.typeID},
		"c0-param2": {"number:0"},
		"c0-param3": {"number:"+content.ID},
		"batchId": {strconv.Itoa(int(time.Now().Unix()* 1000))},
	}
	return utils.HttpPostForm(config.GetLessonUnitLearnVo,requestBody)
}


// 从单个视频的id等信息提取其下载地址（需要构造post请求）
func parseIcourse163VideoFile(ID string, cookie string)string{
	requestBody := url.Values{
		"bizId": { ID },
		"bizType" : {"1"},
		"contentType": {"1"},
	}
	csrfKey := GetCsrfKey(cookie)
	sigText:= utils.HttpPostCookie(config.GetResourceToken+csrfKey,cookie,requestBody)
	signature := utils.MatchAll(sigText,"signature\":\"(.*?)\"")[0][1]
	videoId := utils.MatchAll(sigText,"videoId\":([0-9]*?),")[0][1]
	videoDownloadUrl := fmt.Sprintf("%s?videoId=%s&signature=%s&clientType=1", config.GetVideo, videoId, signature)
	return utils.HttpGet(videoDownloadUrl)
}

//从视频文件的解析结果中提取网址,输入清晰度，输出下载地址,默认hd
func parseIcourse163Video(contentID string,cookie string,quality int) string {
	text := parseIcourse163VideoFile(contentID,cookie)
	//优先下载清晰度高的视频
	hlsUrl := utils.MatchAll(text,"\"videoUrl\":\"(http.*?)\",")
	if (len(hlsUrl) > quality ){
		return hlsUrl[quality-1][1]
	} else if (hlsUrl != nil){
		return hlsUrl[0][1]
	}
	return ""
}

//从富文本文件的jsonText出发，返回文件下载链接以及文件名
func parseIcourse163RichText(jsontext string)(string,string){
	fileName := utils.MatchAll(jsontext,`fileName.{2}:.{2}(.*)\\\"}`)[0][1]
	nosKey := utils.MatchAll(jsontext,`nosKey.{2}:.{2}(.*?)\\`)[0][1]
	//将参数连接成网址
	params := url.Values{}
	params.Add("fileName",fileName)
	params.Add("nosKey",nosKey)
	richTextUrl := "https://www.icourse163.org/course/attachment.htm"+params.Encode()
	return richTextUrl,fileName
}


//从pdf文件的解析结果中获取下载地址
func parseIcourse163PDF(content Icourse163File)string{
	text := parseIcourse163File(content)
	//fmt.Println(text)
	PDFUrl := utils.MatchAll(text,"textOrigUrl:\"(http.*.pdf)")[0][1]
	return PDFUrl
}

//从起始的网页地址，得到课程名称
func GetIcourse163Name(url string)string{
	data:=utils.HttpGet(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	utils.Check(err)
	courseName := doc.Find("head > title").Text()
	return courseName
}

//从cookie中得到csrfKey
func GetCsrfKey(cookie string)string{
	return utils.MatchAll(cookie,"NTESSTUDYSI=(.*?);")[0][1];
}

//总的解析函数
func DownloadIcourse163(url string,options string,cookie string)bool{
	if utils.FileExists(cookie){
		data, _ := ioutil.ReadFile(cookie)
		cookie = string(data)
	}
	courseName := GetChinesemoocName(url)
	termId := utils.MatchAll(url,"tid=([0-9]*)")
	if termId != nil{
		termIDNum := termId[0][1]
		//所有文件信息
		text := HttpPostTermID(termIDNum)
		Icourse163Video,Icourse163PPT,Icourse163RichText := ExtractIcourse163Res(text,cookie)
		switch options{
			//下载所有文件
			case "all":
				utils.M3u8DownloadList(Icourse163Video,courseName)
				download.DownloadFiles(Icourse163PPT,courseName)
				download.DownloadFiles(Icourse163RichText,courseName)
			//只下载视频文件
			case "video":
				utils.M3u8DownloadList(Icourse163Video,courseName)
			//只下载ppt文件
			case "PPT":
				download.DownloadFiles(Icourse163PPT,courseName)
			//只下载富文本文件
			case "RichText":
				download.DownloadFiles(Icourse163RichText,courseName)
		}

	} else{
		//网址不符合格式
		fmt.Printf("this website %s is not supported now",url)
		return true
	}
	return true
}

