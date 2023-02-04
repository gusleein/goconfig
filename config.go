package config

import (
	log "github.com/gusleein/golog"
	"github.com/spf13/viper"
	"net"
	"strconv"
)

var (
	InterfaceIP = getMyIP()
	DebugMode   bool
)

// эта обертка нужна, чтобы логировать отсутствие параметра
var config *viper.Viper

func Init(env string) {
	config = viper.New()

	config.AddConfigPath("./config") // path to folder

	config.SetConfigName(env) // (without extension)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	// viper сам конвертирует on/off yes/no в true/false
	DebugMode = config.GetBool("debug")
}

func GetString(key string) (value string) {
	value = config.GetString(key)
	if value == "" {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	return
}

func GetInt(key string) int {
	// все как строку берем
	strValue := config.GetString(key)
	if strValue == "" {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	value64, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		log.Errorw("config: invalid type of int", "key", key, "value", strValue)
		return 0
	}
	return int(value64)
}

func GetStringSlice(key string) []string {
	slice := config.GetStringSlice(key)
	if slice == nil || len(slice) == 0 {
		log.Errorw("config: value is empty, key: "+key, "key", key)
	}
	return slice
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
