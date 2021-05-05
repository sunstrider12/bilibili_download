package config

import "github.com/spf13/viper"

type GlobalConfig struct {
	RoomInfo []RoomInfo
	Cookie string
	Dir string
}

type RoomInfo struct {
	RoomNum string //房间号
	CheckTime int //多长时间检查一次(秒为单位)
	SaveSpace int64 //文件超过多少MB切割 0为不切割
	BeginTime int //从几点开始 一直等到主播开播,才开始录
	EndTime int //到几点结束 之后就算主播还在播也不在录播
	NeedTicker bool //是否启用begin time 和end time
}

var config GlobalConfig

// ParseConfig 解析配置文件
func ParseConfig(cfg string) {
	viper.SetConfigFile(cfg)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func Config() GlobalConfig  {
	return config
}
