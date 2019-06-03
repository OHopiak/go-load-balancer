package main

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/OHopiak/fractal-load-balancer/master"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"strconv"
)

type Config struct {
	IP           string
	Port         int
	DbDriver     string
	DbConnString string
}

func (c Config) DbConfig() core.DatabaseConfig {
	return core.DatabaseConfig{
		Driver:           c.DbDriver,
		ConnectionString: c.DbConnString,
	}
}

func defaultConfig() *Config {
	return &Config{
		IP:           "localhost",
		Port:         8000,
		DbDriver:     "postgres",
		DbConnString: "host=localhost port=5432 user=orest dbname=load_balancer password=password123$ sslmode=disable",
	}
}

// LoadConfigFromEnv loads config using env vars
//
// MASTER_IP
// MASTER_PORT
// DB_DRIVER
// DB_CONN_STRING
//
func LoadConfigFromEnv() *Config {
	config := defaultConfig()
	masterIP := os.Getenv("MASTER_IP")
	masterPortString := os.Getenv("MASTER_PORT")
	dbDriver := os.Getenv("DB_DRIVER")
	dbConnString := os.Getenv("DB_CONN_STRING")
	if masterIP != "" {
		config.IP = masterIP
	}
	if masterPortString != "" {
		port, err := strconv.ParseInt(masterPortString, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		config.Port = int(port)
	}
	if dbDriver != "" {
		config.DbDriver = dbDriver
	}
	if dbConnString != "" {
		config.DbConnString = dbConnString
	}

	return config
}

func main() {
	config := LoadConfigFromEnv()
	server := master.New(config.DbConfig())
	if config.IP != "" {
		server = server.WithIP(config.IP)
	}
	if config.Port != 0 {
		server = server.WithPort(config.Port)
	}

	server.Server.Logger.Info("Starting a master")
	// Start server
	server.Start()
}
