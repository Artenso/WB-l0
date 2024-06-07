package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type natsConfig struct {
	ClusterID string `json:"cluster_id"`
	ClientID  string `json:"client_id"`
}

type httpConfig struct {
	HttpPort string `json:"http_port"`
}

type DbConfig struct {
	Db       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"db_name"`
}

type config struct {
	NatsConfig natsConfig `json:"nats_config"`
	HttpConfig httpConfig `json:"http_config"`
	DbConfig   DbConfig   `json:"db_config"`
}

var globalConfig config

// Read reads config from file and parse it into config struct
func Read(pathToConfig string) {

	configFile, err := os.ReadFile(pathToConfig)
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}

	err = json.Unmarshal(configFile, &globalConfig)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %s", err)
	}
}

// GetClusterID returns cluster_ID
func GetClusterID() string {
	return globalConfig.NatsConfig.ClusterID
}

// GetClientID returns client_ID
func GetClientID() string {
	return globalConfig.NatsConfig.ClientID
}

// GetHttpPort returns http port
func GetHttpPort() string {
	return globalConfig.HttpConfig.HttpPort
}

// GetDbDSN creates connection string and return it
func GetDbDSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		globalConfig.DbConfig.Db,
		globalConfig.DbConfig.User,
		globalConfig.DbConfig.Password,
		globalConfig.DbConfig.Host,
		globalConfig.DbConfig.Port,
		globalConfig.DbConfig.DbName,
	)

}
