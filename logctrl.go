package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/mholt/archiver"
)

type PollLogLvCb func() int

var g_advLog *Logger = nil

type LogCtrl struct {
	CompressMethod string // 压缩方式
	Directory      string // 文件路径
	TrigerSize     int64  // 触发压缩的大小
	//普通日志
	Std        *Logger //普通日志对象
	LogName    string  // 日志服务名
	FileName   string  // 原始文件名
	CompresCnt int
	//重要日志
	AdvLog        *Logger // 重要日志对象
	LogNameAdv    string  // 重要日志服务名
	FileNameAdv   string  // 原始文件名
	CompresAdvCnt int
	LogAdvEn      bool

	StdOut        bool  // 是否终端打印日志
	AllZipMaxSize int64 // 压缩文件最大总大小
	ZipMaxCount   int64 // 压缩文件最大个数
	logLvChk      PollLogLvCb //动态更新外部应用设置的日志等级

}

func getFileSize(FileName string) int64 {
	fi, err := os.Stat(FileName)
	if err != nil {
		return -1
	}
	//fmt.Println("file size is ",fi.Size(),err)
	return fi.Size()
}

func fileRename(src_file, dst_file string) bool {
	err := os.Rename(src_file, dst_file) //重命名,验证到windows重命名文件会异常，暂不在windows用
	if err != nil {
		//如果重命名文件失败,则输出错误 file rename Error!
		fmt.Println("[Logctrl] file rename Error:", err)
		return false
	}
	return true
}

func fileZip(src_file, zip_file string) bool {
	// 压缩文件
	err := archiver.Archive([]string{src_file}, zip_file) //验证到windows压缩文件会异常，暂不在windows用
	if err != nil {
		fmt.Println("[Logctrl] File zip fail,err:", err)
		return false
	}
	err = os.Remove(src_file) //删除文件
	if err != nil {
		fmt.Println("[Logctrl] File remove fail: ", err)
		return false
	}
	return true
}

var DIRCHAR string = "/"

func init() {
	sysType := runtime.GOOS
	if sysType == "windows" {
		// windows系统
		DIRCHAR = "\\"
	}
}
func (this *LogCtrl) resetLogWriter(l *Logger, stdout bool, fileName string) {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	l.SetFlags(Ldate | Ltime | Lshortfile | Lmicroseconds)

	writers := []io.Writer{file}
	if stdout {
		writers = append(writers, os.Stdout)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	l.SetOutput(fileAndStdoutWriter)
}

/**
*	@brief 初始化函数
*   参数: FileName   日志文件路径文件名
*   参数: params  为可变参数，可不填，亦可填多个
*        int类型参数   --日志等级说明：  FATAL:1  ERROR:2  WARN:3  INFO:4 DEBUG:5(默认) TRACE:6
*        string类型参数--配置字符串说明： "ADV" 支持重要日志单独输出（WARN以上的等级会同时输出到指定的error文件里）
 */
func (this *LogCtrl) Init(FileName string, params ...interface{}) {

	var level int = 0
	var configs string = ""
	//配置参数解析
	for _, u := range params {
		switch u.(type) {
		case int:
			level = u.(int)
		case int64:
			level = int(u.(int64))
		case int32:
			level = int(u.(int32))
		case string:
			configs += u.(string) + "|"
		default:
			fmt.Printf("Unknow type:%v\n", reflect.TypeOf(u))
		}
	}
	//获取配置
	var advEnable bool = false
	if strings.Contains(configs, "ADV") {
		advEnable = true
	}
	//初始化
	this.LogInit(FileName, true, level, advEnable)
}

func parseParam(params ...interface{}) (level int, advEnable bool) {
	configs := ""
	for _, u := range params {
		switch u.(type) {
		case int:
			level = u.(int)
		case int64:
			level = int(u.(int64))
		case int32:
			level = int(u.(int32))
		case string:
			configs += u.(string) + "|"
		default:
			fmt.Printf("Unknow type:%v\n", reflect.TypeOf(u))
		}
	}
	if strings.Contains(configs, "ADV") {
		advEnable = true
	}
	fmt.Printf("[logctrl] %v Lv: %v  COnfigs:%v\n", len(params), level, configs)
	return
}

func (this *LogCtrl) LogInit(FileName string, StdOut bool, level int, advLogEn bool) {
	if len(FileName) == 0 {
		FileName = "default.log"
		this.Init(FileName)
	}
	//level, advLogEn := parseParam(params...)
	if level > 0 {
		fmt.Printf("Find set level=%v\n", level)
		std.LogLevel(level)
	}
	this.Std = std
	this.StdOut = StdOut
	this.FileName = FileName
	this.CompressMethod = "zip"            // 默认zip格式
	this.TrigerSize = 20 * 1024 * 1024     // 默认20MB进行压缩
	this.AllZipMaxSize = 100 * 1024 * 1024 // 默认压缩大小超过100MB，删除最前的日志
	this.ZipMaxCount = 30                  // 默认压缩包个数最大20
	n := strings.LastIndexByte(FileName, byte(DIRCHAR[0]))
	if n < 0 {
		this.Directory = "." // 当前文件夹
		n = 0
	} else {
		this.Directory = FileName[0:n]
	}
	s := strings.LastIndexByte(FileName, '.')
	if s < 0 {
		this.LogName = FileName[n:len(FileName)]
	} else {
		if s <= n {
			this.LogName = FileName
		} else {
			this.LogName = FileName[n+1 : s]
		}
	}
	if advLogEn {
		this.LogAdvEn = advLogEn
		this.LogNameAdv = this.LogName + "_error"
		this.FileNameAdv = strings.Replace(this.FileName, this.LogName, this.LogNameAdv, -1)
		this.AdvLog = New(os.Stderr, "", LstdFlags)
		g_advLog = this.AdvLog
		//this.resetLogWriter(this.AdvLog, false, this.FileNameAdv)
	}
	//this.resetLogWriter(this.Std, this.StdOut, this.FileName)

	fmt.Println("LogFile: ", this.FileName)
	fmt.Println("LogDir: ", this.Directory)
	fmt.Println("LogMeth: ", this.CompressMethod)
	fmt.Println("LogName: ", this.LogName)
	if advLogEn {
		fmt.Println("LogNameAdv: ", this.LogNameAdv)
	}
}
func (this *LogCtrl) doLogJobs(obj *Logger, stdout bool, logName, fileName string) (zip int) {
	//本函数检查日志文件是否达到压缩条件，有则压缩并重置新的日志文件
	zip = 0
	filsize := getFileSize(fileName)
	if filsize > this.TrigerSize {
		//fmt.Println("Start zip..")
		header := logName + ".log@" + time.Now().Format("20060102_150405")
		newlfile := this.Directory + DIRCHAR + header + ".log"
		zipFile := this.Directory + DIRCHAR + header + "." + this.CompressMethod
		if fileRename(fileName, newlfile) {
			this.resetLogWriter(obj, stdout, fileName)
			fileZip(newlfile, zipFile)
			zip = 1
		}
	}
	return
}
func (this *LogCtrl) Run() {
	if len(this.LogName) == 0 {
		this.Init("default.log")
	}
	//进来首次检查是否需要清理日志压缩文件
	this.doclear(this.LogName)
	if this.LogAdvEn {
		this.doclear(this.LogNameAdv)
	}

	this.resetLogWriter(this.Std, this.StdOut, this.FileName)
	if this.LogAdvEn {
		this.resetLogWriter(this.AdvLog, false, this.FileNameAdv)
	}
	var CurLoglv int
	for {
		zipcnt := this.doLogJobs(this.Std, this.StdOut, this.LogName, this.FileName)
		if zipcnt > 0 {
			this.CompresCnt += zipcnt
			this.doclear(this.LogName) //有新增压缩包，则需判断是否清除普通日志
		}
		if this.LogAdvEn {
			//如果开启了重要日志，则需处理重要日志的事务
			zipcnt = this.doLogJobs(this.AdvLog, false, this.LogNameAdv, this.FileNameAdv)
			if zipcnt > 0 {
				this.CompresAdvCnt += zipcnt
				this.doclear(this.LogNameAdv) //有新增压缩包，则需判断是否清除重要日志
			}
		}
		if this.logLvChk != nil {
			loglev := this.logLvChk()
			if loglev != CurLoglv {
				LogLevel(loglev)
				CurLoglv = loglev
			}
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func (this *LogCtrl) ResetCompressSize(s int64 /*字节*/) {
	if s < 100*1024 {
		return
	}
	this.TrigerSize = s
	Println("[Logctrl] Logctrl TriggerSize = ", this.TrigerSize)
}
func (this *LogCtrl) SetClearSize(s int64 /*字节*/) {
	if s < 100*1024 {
		return
	}
	this.AllZipMaxSize = s
	Println("[Logctrl] Logctrl AllZipMaxSize = ", this.AllZipMaxSize)
}
func (this *LogCtrl) SetZipMaxCount(count int64) {
	if count < 30 {
		return
	}
	this.ZipMaxCount = count
	Println("[Logctrl] Logctrl ZipMaxCount = ", this.ZipMaxCount)
}

func (this *LogCtrl) SetLogLevelChkCb(lchk PollLogLvCb) {
	if lchk != nil {
		this.logLvChk = lchk
	}
}

func getFileModTime(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

type zipfile struct {
	FileName       string
	PathName       string
	LastModifyTime int64
	FileSize       int64
}

func (this *LogCtrl) doclear(logName string) {
	fileArray := []string{}
	//zipFInfo := make([]zipfile, 0)
	pwd := this.Directory
	//获取当前目录下的所有文件或目录信息
	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			//fmt.Println(path) //打印path信息
			//fmt.Println(info.Name()) //打印文件或目录名
			if fileArray != nil {
				fileArray = append(fileArray, info.Name())
			}
		}
		return nil
	})
	type LogFile struct {
		Target  string
		Size    int64
		ModifyT int64
	}
	delTargets := []LogFile{}

	zipCount := int64(0)
	zipSize := int64(0)
	for _, f := range fileArray {
		if !strings.HasPrefix(f, logName+".log@") {
			//fmt.Printf("f %s, not contain %s\n", f, this.LogName)
			continue
		}
		if !strings.HasSuffix(f, "."+this.CompressMethod) {
			//fmt.Printf("f %s, not contain %s\n", f, "." + this.CompressMethod)
			continue
		}

		PathName := this.Directory + DIRCHAR + f
		LastModifyTime := getFileModTime(PathName)
		FileSize := getFileSize(PathName)
		zipSize = zipSize + FileSize
		t := LogFile{
			Target:  PathName,
			Size:    FileSize,
			ModifyT: LastModifyTime,
		}
		delTargets = append(delTargets, t)
		zipCount++
	}
	if len(delTargets) > 1 {
		//按时间顺序排列
		sort.SliceStable(delTargets, func(i int, j int) bool {
			return delTargets[i].ModifyT < delTargets[j].ModifyT
		})
	}
	// fmt.Printf("zipCount: %d , conut: %d\n", zipCount, len(delTargets))
	// fmt.Printf("TotalZie: %d , MaxSize: %d\n", zipSize, this.AllZipMaxSize)
	//按时间顺序判断是否删除
	for len(delTargets) > int(this.ZipMaxCount) || zipSize > this.AllZipMaxSize {
		delTarget := delTargets[0]
		fmt.Println("[logctrl]clear log file:", delTarget)
		err := os.Remove(delTarget.Target) //删除文件
		if err != nil {
			Println("[Logctrl] File remove fail: ", err)
			return
		}
		zipSize -= delTarget.Size
		if len(delTargets) > 1 {
			delTargets = delTargets[1:]
		}
	}

	return
}
