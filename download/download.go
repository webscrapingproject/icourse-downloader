package download

import (
	"fmt"
	"github.com/schollz/progressbar"
	"icourse/config"
	"icourse/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)



//批量下载file数组里的文件
func DownloadFiles(files []utils.File){
	for _,file:=range(files){
		//下载到指定文件夹
		fmt.Println("\n"+filepath.Join(config.OutputPath,file.FilePATH))
		DownloadFile(file.FileURL,filepath.Join(config.OutputPath,file.FilePATH))
	}
}
//根据网络URL获得文件的大小
func getFileSize(url string) int64 {
	res, err := http.Head(url)
	utils.Check(err)
	size, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	downloadSize := int64(size)
	return downloadSize
}

// 还没有实现，在这里占个位
func aria2Download(url string,filePath string){
}


func DownloadFile(url string,filePath string){
	//获取需要下载的文件大小
	dataSize:=getFileSize(url)
	//获取需要写入的信息
	res, err := http.Get(url)
	utils.Check(err)
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	//body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	//检查文件是否存在
	if utils.FileExists(filePath){
		fmt.Printf("file already exists, skipping\n")
		return
	}
	//检查目录是否存在
	if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
		//建立目录
		_ = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	}
	//建立进程条,设置参数显示下载速度和下载进度
	bar := progressbar.NewOptions(
		int(dataSize),
		progressbar.OptionSetBytes(int(dataSize)),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
	)

	// 创建文件
	dest, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Can't create %s: %v\n", filePath, err)
		return
	}
	defer dest.Close()
	// 从reader读入文件
	out := io.MultiWriter(dest, bar)
	_, _ = io.Copy(out, res.Body)
}
