package main

import (
	"ginEssential/common"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"github.com/spf13/viper"
)



func main() {
	InitConfig()
	db := common.InitDB()
	defer common.Close(db)


	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	if port != ""{
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

func InitConfig() {
	workDir,_ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil{

	}
}




