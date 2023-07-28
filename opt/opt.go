package opt

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type pg struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
}

type esIndex struct {
	// add your elasticsearch index here
	// example:
	// Book string `mapstructure:"book"`
}

type es struct {
	Endpoints   []string `mapstructure:"endpoints"`
	Index       esIndex  `mapstructure:"index"`
	Sniff       bool     `mapstructure:"sniff"`
	EnableTrace bool     `mapstructure:"enable_trace"`
}

type etcd struct {
	Endpoints      []string `mapstructure:"endpoints"`
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	EnableResolver bool     `mapstructure:"enable_resolver"`
}

type nebula struct {
	Space     string   `mapstructure:"space"`
	Endpoints []string `mapstructure:"endpoints"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
	PoolSize  int      `mapstructure:"pool_size"`
	LogLevel  string   `mapstructure:"log_level"`
}

type s3Bucket struct {
	// add your bucket here
	// example:
	// Avatar string `mapstructure:"avatar"`
}

type s3 struct {
	Endpoint  string   `mapstructure:"endpoint"`
	AccessKey string   `mapstructure:"access_key"`
	SecretKey string   `mapstructure:"secret_key"`
	Bucket    s3Bucket `mapstructure:"buckets"`
}

type service struct {
	Name       string `mapstructure:"name"`
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`
}

type config struct {
	Project  string `mapstructure:"project"`
	LogLevel string `mapstructure:"log_level"`
}

var (
	configFile string
	Cfg        = new(config)
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600)
	flag.StringVar(&configFile, "c", "etc/config.json", "")
}

func MustInitConfig() {
	flag.Parse()

	var (
		err error
	)

	viper.SetConfigType("json")
	viper.SetConfigFile(configFile)
	if err = viper.ReadInConfig(); err != nil {
		logrus.Panicf("read in config file err: %v", err)
	}

	if err = viper.Unmarshal(Cfg); err != nil {
		logrus.Panicf("unmarshal config file err: %v", err)
	}

	switch strings.ToLower(Cfg.LogLevel) {
	case "err", "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				file := path.Base(frame.File)
				return "", " " + file + ":" + strconv.Itoa(frame.Line)
			},
		})
	}

	logrus.Infof("read in config detail: %+v", *Cfg)
}
