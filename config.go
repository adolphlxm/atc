package atc

import (
	"errors"
	"github.com/adolphlxm/atc/utils"
	"github.com/lxmgo/config"
	"os"
	"path/filepath"
)

var (
	Aconfig   *Config
	AppConfig *appConfig
)

type Config struct {
	Runmode string
	Debug   bool
	AppName string

	// HTTP/Websocket
	HTTPSupport      bool
	HTTPAddr         string
	HTTPPort         string
	HTTPSCertFile    string
	HTTPSKeyFile     string
	HTTPQTimeout     int
	HTTPReadTimeout  int
	HTTPWriteTimeout int

	PostMaxMemory int64

	// Front
	FrontSupport   bool
	FrontDir       []string
	FrontDirectory bool
	FrontSuffix    string
	FrontHost      string

	// Thrift
	ThriftSupport       bool
	ThriftDebug         bool
	ThriftQTimeout      int
	ThriftClientTimeout int
	LogSupport          bool
	LogLevel            int
	LogOutput           string

	// Grpc
	GrpcSupport bool

	// Orm
	OrmSupport      bool
	ThriftAddr      string
	ThriftPort      string
	ThriftSecure    string
	ThriftProtocol  string
	ThriftTransport string

	// Log
	OrmLogLevel   string
	OrmAliasNames []string

	// queue
	QueuePublisherSupport bool
	QueueConsumerSupport  bool

	// mongodb
	MgoSupport bool
}

// Parsing the configuration
func ParseConfig(confName, runmode string) error {
	var err error

	AppConfig, err = NewAppConfig(confName)
	if err != nil {
		return err
	}

	Aconfig = &Config{
		Runmode:          "local",
		Debug:            false,
		AppName:          "ATC",
		HTTPAddr:         "",
		HTTPPort:         "9000",
		HTTPQTimeout:     60,
		HTTPReadTimeout:  0,
		HTTPWriteTimeout: 0,

		PostMaxMemory: 1 << 26, // 64MB

		FrontSupport:   false,
		FrontDir:       []string{"index", "assets"},
		FrontDirectory: false,
		FrontSuffix:    "html",
		FrontHost:      "",

		ThriftSupport:       false,
		ThriftDebug:         false,
		ThriftQTimeout:      300, // 5min
		ThriftClientTimeout: 10,
		ThriftAddr:          "",
		ThriftPort:          "9090",
		ThriftSecure:        "false",
		ThriftProtocol:      "tbinary",
		ThriftTransport:     "tframed",

		GrpcSupport: false,

		LogSupport: true,
		LogLevel:   0,
		LogOutput:  "stdout",

		OrmSupport:    false,
		OrmLogLevel:   "LOG_OFF",
		OrmAliasNames: []string{},

		QueuePublisherSupport: false,
		QueueConsumerSupport:  false,

		MgoSupport: false,
	}

	if runmode != "" {
		Aconfig.Runmode = runmode
	}

	Aconfig.HTTPSupport = AppConfig.DefaultBool("http.support", Aconfig.HTTPSupport)
	Aconfig.HTTPAddr = AppConfig.DefaultString("http.addr", Aconfig.HTTPAddr)
	Aconfig.HTTPPort = AppConfig.DefaultString("http.port", Aconfig.HTTPPort)
	Aconfig.HTTPQTimeout = AppConfig.DefaultInt("http.qtimeout", Aconfig.HTTPQTimeout)
	Aconfig.HTTPReadTimeout = AppConfig.DefaultInt("http.readtimeout", Aconfig.HTTPReadTimeout)
	Aconfig.HTTPWriteTimeout = AppConfig.DefaultInt("http.readtimeout", Aconfig.HTTPWriteTimeout)

	Aconfig.PostMaxMemory = int64(AppConfig.DefaultInt("post.maxmemory", int(Aconfig.PostMaxMemory)))
	Aconfig.Debug = AppConfig.DefaultBool("app.debug", Aconfig.Debug)
	Aconfig.AppName = AppConfig.DefaultString("app.name", Aconfig.AppName)
	Aconfig.FrontSupport = AppConfig.DefaultBool("front.support", Aconfig.FrontSupport)
	Aconfig.FrontDir = AppConfig.DefaultStrings("front.dir", Aconfig.FrontDir)
	Aconfig.FrontDirectory = AppConfig.DefaultBool("front.directory", Aconfig.FrontDirectory)
	Aconfig.FrontSuffix = AppConfig.DefaultString("front.suffix", Aconfig.FrontSuffix)
	Aconfig.FrontHost = AppConfig.DefaultString("front.host", Aconfig.FrontHost)

	Aconfig.ThriftSupport = AppConfig.DefaultBool("thrift.support", Aconfig.ThriftSupport)
	Aconfig.ThriftDebug = AppConfig.DefaultBool("thrift.debug", Aconfig.ThriftDebug)
	Aconfig.ThriftQTimeout = AppConfig.DefaultInt("thrift.qtimeout", Aconfig.ThriftQTimeout)
	Aconfig.ThriftClientTimeout = AppConfig.DefaultInt("thrift.client.timeout", Aconfig.ThriftClientTimeout)
	Aconfig.ThriftAddr = AppConfig.DefaultString("thrift.addr", Aconfig.ThriftAddr)
	Aconfig.ThriftPort = AppConfig.DefaultString("thrift.port", Aconfig.ThriftPort)
	Aconfig.ThriftSecure = AppConfig.DefaultString("thrift.secure", Aconfig.ThriftSecure)
	Aconfig.ThriftProtocol = AppConfig.DefaultString("thrift.protocol", Aconfig.ThriftProtocol)
	Aconfig.ThriftTransport = AppConfig.DefaultString("thrift.transport", Aconfig.ThriftTransport)

	Aconfig.GrpcSupport = AppConfig.DefaultBool("grpc.support", Aconfig.GrpcSupport)

	Aconfig.LogSupport = AppConfig.DefaultBool("log.support", Aconfig.LogSupport)
	logLevel := AppConfig.DefaultString("log.level", "LevelDebug")
	switch logLevel {
	case "LevelFatal":
		Aconfig.LogLevel = 0
	case "LevelError":
		Aconfig.LogLevel = 1
	case "LevelWarn":
		Aconfig.LogLevel = 2
	case "LevelNotice":
		Aconfig.LogLevel = 3
	case "LevelInfo":
		Aconfig.LogLevel = 4
	case "LevelTrace":
		Aconfig.LogLevel = 5
	case "LevelDebug":
		Aconfig.LogLevel = 6
	}
	Aconfig.LogOutput = AppConfig.DefaultString("log.output", Aconfig.LogOutput)

	Aconfig.OrmSupport = AppConfig.DefaultBool("orm.support", Aconfig.OrmSupport)
	Aconfig.OrmLogLevel = AppConfig.DefaultString("orm.log.level", Aconfig.OrmLogLevel)
	Aconfig.OrmAliasNames = AppConfig.DefaultStrings("orm.aliasnames", Aconfig.OrmAliasNames)

	Aconfig.QueuePublisherSupport = AppConfig.DefaultBool("queue.publisher.support", Aconfig.QueuePublisherSupport)
	Aconfig.QueueConsumerSupport = AppConfig.DefaultBool("queue.consumer.support", Aconfig.QueueConsumerSupport)

	Aconfig.MgoSupport = AppConfig.DefaultBool("mgo.support", Aconfig.MgoSupport)
	return nil
}

type appConfig struct {
	config config.ConfigInterface
}

func NewAppConfig(confName string) (*appConfig, error) {
	conf, err := config.NewConfig(confName)
	if err != nil {
		return nil, err
	}
	return &appConfig{conf}, nil
}
func (a *appConfig) Set(key string, value string) error {
	return a.config.Set(Aconfig.Runmode+"::"+key, value)
}
func (a *appConfig) String(key string) string {
	return a.config.String(Aconfig.Runmode + "::" + key)
}
func (a *appConfig) Strings(key string) []string {
	return a.config.Strings(Aconfig.Runmode + "::" + key)
}

func (a *appConfig) Bool(key string) (bool, error) {
	return a.config.Bool(Aconfig.Runmode + "::" + key)
}

func (a *appConfig) Int(key string) (int, error) {
	return a.config.Int(Aconfig.Runmode + "::" + key)
}

func (a *appConfig) DefaultString(key string, defaultVal string) string {
	if v := a.String(key); v != "" {
		return v
	}
	return defaultVal
}
func (a *appConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := a.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}
func (a *appConfig) DefaultBool(key string, defaultVal bool) bool {
	if b, err := a.Bool(key); err == nil {
		return b
	}
	return defaultVal
}
func (a *appConfig) DefaultInt(key string, defaultVal int) int {
	if b, err := a.Int(key); err == nil {
		return b
	}
	return defaultVal
}

// Initialize config.
func initConfig(configFile, runMode string) {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !utils.FileExists(configFile) {
		configFile, err = matchingConfig(workPath)
		if err != nil {
			panic(err)
		}
	}

	err = ParseConfig(configFile, runMode)
	if err != nil {
		panic(err)
	}
}

func matchingConfig(workPath string) (string, error) {
	configFile := filepath.Join(workPath, "conf", "app.ini")
	if !utils.FileExists(configFile) {
		configFile = filepath.Join("../", "conf", "app.ini")
		if !utils.FileExists(configFile) {
			configFile = filepath.Join("../../", "conf", "app.ini")
			if !utils.FileExists(configFile) {
				configFile = filepath.Join("../../../", "conf", "app.ini")
				if !utils.FileExists(configFile) {
					return "", errors.New(configFile + ": no such file or directory.")
				}
			}
		}
	}
	return configFile, nil
}
