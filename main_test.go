package log

/*
import (
	"fmt"
	"time"

	"github.com/sheazp/log"
)

var lct *log.LogCtrl

func main_test() {
	lct = &log.LogCtrl{}
	logFile := "./log/mylog.go.log"
	//lct.Init(logFile) //默认DEBUG等级
	//lct.Init(logFile, log.INFO) //初始化为INFO等级
	//lct.Init(logFile, "ADV")    //开启重要日志配置,默认DEBUG等级
	lct.Init(logFile, log.ERROR, "ADV") //初始化为ERROR等级，并且开启重要日志配置
	lct.ResetCompressSize(101 * 1024)
	lct.SetZipMaxCount(5)
	lct.SetClearSize(100 * 1024) // 500MB
	go lct.Run()

	go func() {
		for {
			str := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
			log.Errorf("-------------\n")
			log.Fatalf("Fatalf...%v\n", str)
			log.Infof("Infof...%v\n", str)
			log.Debugf("Debugf...%v\n", str)
			log.Errorf("Errorf...%v\n", str)
			log.Printf("Printf...%v\n", str)
			log.Warnf("Warnf...%v\n", str)
			log.Tracef("Tracef...%v\n", str)
			time.Sleep(time.Duration(50) * time.Millisecond)
		}
	}()

	for {
		var lv int = 0
		fmt.Scan(&lv)
		fmt.Printf("set level change:%v -> %v\n", log.LogLevel(0) , log.LogLevel(lv))
	}
}
*/
