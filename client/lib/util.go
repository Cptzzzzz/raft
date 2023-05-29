package lib

import "github.com/spf13/viper"

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	Hosts = viper.GetStringSlice("hosts")
	Judge = viper.GetString("judge")
}
