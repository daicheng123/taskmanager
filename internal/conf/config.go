package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"taskmanager/internal/consts"
)

var defaultConfig *Config

type Config struct {
	Web      *WebConfig   `yaml:"web"`
	DataBase *DBConfig    `yaml:"database"`
	Redis    *RedisConfig `yaml:"redis"`
	Log      *LogConfig   `yaml:"log"`
}

type WebConfig struct {
	Mode string `yaml:"mode"`
	Port uint   `yaml:"port"`
}

type DBConfig struct {
	Driver     string `yaml:"driver"`
	Address    string `yaml:"address"`
	Port       uint   `yaml:"port"`
	Db         string `yaml:"db"`
	DbUser     string `yaml:"dbUser"`
	DbPassword string `yaml:"dbPassword"`
}

type RedisConfig struct {
	Address     string `yaml:"address"`
	UsePassword bool   `yaml:"usePassword"`
	Password    string `yaml:"password"`
}

type LogConfig struct {
	DebugMode bool   `yaml:"debugMode"`
	LogLevel  string `yaml:"logLevel"`
	LogPath   string `yaml:"logPath"`
	//PrintFile bool   `yaml:"printFile"`
	//PrintLine bool   `yaml:"printLine"`
	//PrintFunc bool   `yaml:"printFunc"`
}

func LoadConf() {
	configFile := os.Getenv(consts.AppManagerConfPath)
	if len(configFile) == 0 {
		panic("配置文件路径未设置")
		return
	}

	configFile, err := filepath.Abs(configFile)
	if err != nil {
		panic(err)
	}

	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(config, defaultConfig)
	if err != nil {
		panic(err)
	}
}

func GetWebMode() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.Web.Mode
}

func GetWebPort() uint {
	if defaultConfig == nil {
		return 0
	}
	return defaultConfig.Web.Port
}

func GetDBDriver() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.DataBase.Driver
}

func GetDBAddress() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.DataBase.Address
}

func GetDBPort() uint {
	if defaultConfig == nil {
		return 0
	}
	return defaultConfig.DataBase.Port
}

func GetDBUser() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.DataBase.DbUser
}

func GetDbPassword() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.DataBase.DbPassword
}

func GetLogDebugMode() bool {
	if defaultConfig == nil {
		return false
	}
	return defaultConfig.Log.DebugMode
}

func GetLogLevel() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.Log.LogLevel
}

func GetLogPath() string {
	if defaultConfig == nil {
		return ""
	}
	return defaultConfig.Log.LogPath
}
