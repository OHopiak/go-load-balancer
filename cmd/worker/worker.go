package main

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/OHopiak/fractal-load-balancer/worker"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"strconv"
)

type Config struct {
	IP string
	Port int
	MasterIP     string
	MasterPort   int
	DbDriver     string
	DbConnString string
}

func (c Config) MasterHost() core.Host {
	return core.Host{
		IP:   c.MasterIP,
		Port: c.MasterPort,
	}
}

func (c Config) DbConfig() core.DatabaseConfig {
	return core.DatabaseConfig{
		Driver:           c.DbDriver,
		ConnectionString: c.DbConnString,
	}
}

func defaultConfig() *Config {
	return &Config{
		IP: "localhost",
		MasterIP:     "localhost",
		MasterPort:   8000,
		DbDriver:     "postgres",
		DbConnString: "host=localhost port=5432 user=orest dbname=load_balancer password=password123$ sslmode=disable",
	}
}

// LoadConfigFromEnv loads config using env vars
//
// WORKER_IP
// WORKER_PORT
// MASTER_IP
// MASTER_PORT
// DB_DRIVER
// DB_CONN_STRING
//
func LoadConfigFromEnv() *Config {
	config := defaultConfig()
	workerIP := os.Getenv("WORKER_IP")
	workerPortString := os.Getenv("WORKER_PORT")
	masterIP := os.Getenv("MASTER_IP")
	masterPortString := os.Getenv("MASTER_PORT")
	dbDriver := os.Getenv("DB_DRIVER")
	dbConnString := os.Getenv("DB_CONN_STRING")
	if workerIP != "" {
		config.IP = workerIP
	}
	if workerPortString != "" {
		port, err := strconv.ParseInt(workerPortString, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		config.Port = int(port)
	}
	if masterIP != "" {
		config.MasterIP = masterIP
	}
	if masterPortString != "" {
		port, err := strconv.ParseInt(masterPortString, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		config.MasterPort = int(port)
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

	server := worker.New(config.MasterHost(), config.DbConfig())
	if config.IP != "" {
		server = server.WithIP(config.IP)
	}
	if config.Port != 0 {
		server = server.WithPort(config.Port)
	}

	server.Server.Logger.Info("Starting a worker")
	// Start server
	server.Start()
}
