package atc

import "github.com/adolphlxm/atc/logs"

// Initialize logs.
func initLogs(){
	logFile := &logs.File{
		LogDir:        AppConfig.DefaultString("log.dir", ""),
		MaxSize:       uint64(AppConfig.DefaultInt("log.maxsize", 0)),
		Buffersize:    AppConfig.DefaultInt("log.buffersize", 0),
		FlushInterval: uint64(AppConfig.DefaultInt("log.flushinterval", 0)),
	}
	err := logs.SetLogger(Aconfig.LogOutput, logFile)
	if err != nil {
		panic(err)
	}

	if Aconfig.Debug {
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(Aconfig.LogLevel)
	}
}
