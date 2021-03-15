# Icouse-Downloader

icourse-downloader is a simple tool to download course materials of mooc courses

![repo-size](https://img.shields.io/github/repo-size/webscrapingproject/icourse-downloader) ![release](https://img.shields.io/github/v/release/webscrapingproject/icourse-downloader)


Currently supported wesites:
- [Icourses](https://www.icourses.cn/home/)
- [Chinesemooc](http://www.chinesemooc.org)
- [Icourse163](https://www.icourse163.org)
- [Datacamp](https://www.datacamp.com/)

## 1. Uages

Download pre-compiled .exe file(for windows) or binary file(for Linux and mac)

### 1.1 Icourses

#### Demo

```bash
./icourse http://www.icourses.cn/sCourse/course_6447.html
```

#### Options

```bash
icourse -co <option(all,most,videoPPT,assignments,testPaper,shareResource)> -o <outputPath> <url-of-icourse>
```
```
<option>
all Download all the course materials
most Download videos, coursework and test papers
videoPPT Download videos only
assignments Download coursework only
testPaper Download test papers only
shareResource Download other course materials
<outputPath> Specify the output path
<url-of-icourse> the URL of a specific resource，for example：http://www.icourses.cn/sCourse/course_6447.html
```

### 1.2 Chinesemooc

#### Demo
```bash
./icourse http://www.chinesemooc.org/mooc/4880
```
#### Options

```bash
icourse -co <option{all, video , PPT}> -o <outputPath> <url-of-chineseMooc>
```

```
<option>
all Download all the course materials(Videos + Slides)
video Download videos only
PPT Download slides only
```

### 1.3 Icourse163

#### Demo
Take care that pid number must be provided
PDF can be downloaded automatically, you can download video by executing the generated script file.

```bash
./icourse https://www.icourse163.org/learn/SDU-1001907001?tid=1003113029
```
#### Options
```bash
icourse -co <option{all, video , PPT ， RichText}> -o <outputPath> <url-of-icourse163>
```

```
<option>
all Download all the course materials(Videos + Slides + Rich texts)
video Download videos only
videoPPT Download slides only
RichText Download Rich text files only
```

### 1.4 Datacamp

#### Demo

Generate downloading and renaming scripts for videos, subtitles and slides

```bash
./icourse https://learn.datacamp.com/courses/data-visualization-with-ggplot2-part-3
```
#### Options

```bash
icourse -o <outputPath> <url-of-datacamp>
```

## References

1. <https://blog.univerone.com/post/34-datacamp-video-download/>

2. <https://blog.univerone.com/post/23-go-icourse-downloader/>
