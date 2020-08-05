package config
//命令行参数

var (
	//下载文件夹的路径
	OutputPath string
	//下载的内容
	ContentOptions string
	// 显示版本号
	Version bool
	// cookie 内容或者cookie文件地址
	Cookie string
)

const VideoPPT="https://www.icourses.cn/web/sword/portal/shareChapter?cid="
const Assignments="http://www.icourses.cn/web/sword/portal/assignments?cid="
const TestPaper="http://www.icourses.cn/web/sword/portal/testPaper?cid="
const ShareResource="http://www.icourses.cn/web/sword/portal/sharerSource?cid="

//中国大学mooc 三个地址
const GetMocTermDto = "https://www.icourse163.org/dwr/call/plaincall/CourseBean.getLastLearnedMocTermDto.dwr"
const GetVideo = "https://vod.study.163.com/eds/api/v1/vod/video"
const GetResourceToken = "https://www.icourse163.org/web/j/resourceRpcBean.getResourceToken.rpc?csrfKey="
const GetLessonUnitLearnVo = "https://www.icourse163.org/dwr/call/plaincall/CourseBean.getLessonUnitLearnVo.dwr"