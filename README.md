# Icouse-Downloader

icourse-downloader可以根据课程链接下载[爱课程网](https://www.icourses.cn/home/)、 [华文慕课](http://www.chinesemooc.org)、[中国大学MOOC](https://www.icourse163.org) 以及
[datacamp](https://www.datacamp.com/)上的视频以及课件文档等

![repo-size](https://img.shields.io/github/repo-size/webscrapingproject/icourse-downloader) ![release](https://img.shields.io/github/v/release/webscrapingproject/icourse-downloader)

## 1. 使用方法

下载编译好的exe文件（windows平台）或者二进制文件（Linux和mac平台），在cmd或者终端里执行（链接可换成其它课程）

### 1.1 爱课程

```bash
./icourse http://www.icourses.cn/sCourse/course_6447.html
```

### 1.2 华文慕课

```bash
./icourse http://www.chinesemooc.org/mooc/4880
```

### 1.3 中国大学MOOC

注意tid必须存在

pdf直接下载，视频生成使用ffmpeg下载的脚本，可自行执行下载：

```bash
./icourse https://www.icourse163.org/learn/SDU-1001907001?tid=1003113029
```

### 1.4 datacamp

默认获得视频文件、字幕文件以及课件文件的下载地址保存在文件中，并生成重命名的批处理文件

```bash
./icourse https://learn.datacamp.com/courses/data-visualization-with-ggplot2-part-3
```

## 2. 参数说明

### 2.1 爱课程下载参数

```bash
icourse -co <option(all,most,videoPPT,assignments,testPaper,shareResource)> -o <outputPath> <url-of-icourse>
```

基本参数解释如下:

```
<option> 下载内容的选择，仅支持单选
all 下载全部内容
most 下载课件视频、课程作业以及课程试卷
videoPPT 仅下载课件视频
assignments 仅下载课程作业
testPaper 仅下载课程试卷
shareResource 仅下载其它公开资源
<outputPath> 指定下载路径
<url-of-icourse> 课程主页的链接，格式同：http://www.icourses.cn/sCourse/course_6447.html
```

### 2.2 华文慕课下载参数

```bash
icourse -co <option{all, video , PPT}> -o <outputPath> <url-of-chineseMooc>
```

基本参数解释如下:

```
<option> 下载内容的选择，仅支持单选
all 下载全部内容
video 仅下载课件视频
PPT 仅下载课程课件
```

### 2.3 中国大学MOOC下载参数

```bash
icourse -co <option{all, video , PPT ， RichText}> -o <outputPath> <url-of-icourse163>
```

基本参数解释如下:

```
<option> 下载内容的选择，仅支持单选
all 下载全部内容
video 仅下载课件视频
videoPPT 仅下载课程课件
RichText 仅下载课程富文本附件
```

### 2.4 datacamp下载参数

```bash
icourse -o <outputPath> <url-of-datacamp>
```


## 相关博文

1. <https://blog.univerone.com/post/34-datacamp-video-download/>

2. <https://blog.univerone.com/post/23-go-icourse-downloader/>
