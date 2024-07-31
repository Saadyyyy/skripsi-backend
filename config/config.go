package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	SERVERPORT int
	DBPORT     int
	DBHOST     string
	DBUSERNAME string
	DBPASSWORD string
	DBNAME     string
}

func InitConfig() *AppConfig {
	var res = new(AppConfig)
	res = loadConfig()

	if res == nil {
		logrus.Fatal("Config : Cannot start program, failed to load configuration")
		return nil
	}

	return res
}

func loadConfig() *AppConfig {
	var res = new(AppConfig)

	godotenv.Load(".env")

	if val, found := os.LookupEnv("SERVERPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : invalid port value,", err.Error())
			return nil
		}
		res.SERVERPORT = port
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : invalid db port value,", err.Error())
			return nil
		}
		res.DBPORT = port
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		res.DBHOST = val
	}

	if val, found := os.LookupEnv("DBUSERNAME"); found {
		res.DBUSERNAME = val
	}

	if val, found := os.LookupEnv("DBPASSWORD"); found {
		res.DBPASSWORD = val
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		res.DBNAME = val
	}

	return res

}
