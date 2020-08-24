package main

import "github.com/kataras/golog"

func main(){
	var SysLog = New("System", "debug")
	SysLog.Info("gg88g88")
	golog
	//// Default is "[ERRO]"
	//golog.ErrorText("|ERROR|", "100",)
	//// Default is "[WARN]"
	//golog.WarnText("|WARN|", "32")
	//// Default is "[INFO]"
	//golog.InfoText("|INFO|", "34")
	//// Default is "[DBUG]"
	//golog.DebugText("|DEBUG|", "33")
	//
	////Business as usual...
	//
	//golog.SetTimeFormat("2006/01/02 15:04:05")
	//
	//var SuccessLevel golog.Level = 6
	//// Register our level, just three fields.
	//golog.Levels[SuccessLevel] = &golog.LevelMetadata{
	//	Name:      "success",
	//	RawText:     "[SUCC]",
	//	ColorfulText: "32", // Green
	//}
	//myLogger := golog.New()
	//// create a new golog logger
	//
	//l := golog.ParseLevel("success")
	//fmt.Println(l)
	//
	//// "disable"
	//// "fatal"
	//// "error"
	//// "warn"
	//// "info"
	//// "debug"
	//
	//myLogger.SetLevel("debug")
	//myLogger.Println("This is a raw message, no levels, no colors.")
	//myLogger.Info("This is an info message, with colors (if the output is terminal)")
	//myLogger.Warn("This is a warning message")
	//myLogger.Error("This is an error message")
	//myLogger.Debug("This is a debug message")
	//myLogger.Logf(SuccessLevel, "This is a success log message with green color")

}
func New(service string, level string) *golog.Logger{
	log:= golog.New()
	//l := golog.ParseLevel("debug")
	//log.SetLevel(l)
	log.SetTimeFormat("2006/01/02 15:04:18")
	return log

}

