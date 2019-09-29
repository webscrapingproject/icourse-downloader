# Icouse-Downloader
icourse-downloader可以根据课程链接下载[爱课程网](https://www.icourses.cn/home/)
上的视频以及课件文档等
![](https://img.shields.io/github/repo-size/webscrapingproject/icourse-downloader)

## 1.使用方法
下载编译好的exe文件（windows平台）或者二进制文件（Linux和mac平台），在cmd或者终端里执行（链接可换成其它课程）
下载爱课程的视频课程：
```bash
./icourse http://www.icourses.cn/sCourse/course_6447.html
```
下载华文慕课的课程：
```bash
./icourse http://www.chinesemooc.org/mooc/4880
```
## 2.参数说明
### 2.1 爱课程下载参数
```
icourse -co <option(all,most,videoPPT,assignments,testPaper,shareResource)> -o <outputPath> <url-of-icourse>
```
基本参数解释如下:
```
<option>	下载内容的选择，仅支持单选
all	下载全部内容
most	下载课件视频、课程作业以及课程试卷
videoPPT	仅下载课件视频
assignments	仅下载课程作业
testPaper	仅下载课程试卷
shareResource	仅下载其它公开资源
<outputPath>	指定下载路径
<url-of-icourse>	课程主页的链接，格式同：http://www.icourses.cn/sCourse/course_6447.html
```
### 2.2 华文慕课下载参数
```
icourse -co <option{all, video , PPT}> -o <outputPath> <url-of-icourse>
```
基本参数解释如下:
```
<option>	下载内容的选择，仅支持单选
all	下载全部内容
video	仅下载课件视频
PPT	仅下载课程课件
```

