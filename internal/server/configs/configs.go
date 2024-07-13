package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port       = GetConfigs("httpPort") //os.Getenv("HTTP_PORT")
	Redis_Addr = GetConfigs("redisAddr")
	DB_id      = 0
	P          = GetConfigs("path")
	Url        = GetConfigs("url")
)

func GetConfigs(param string) string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Data, exists := os.LookupEnv(param)

	if exists {
		Data = os.Getenv(param)
		log.Printf("%s is %s", param, Data)
	} else {
		log.Printf("%s is missing", param)
	}

	return Data

}
