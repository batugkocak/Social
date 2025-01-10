package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	loadErr := godotenv.Load(curDir + "/.env")
	if loadErr != nil {
		log.Fatalln("can't load env file from current directory: " + curDir)
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}
