package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Firmware FirmwareConfig
	MQTT    MQTTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Type string
	Path string
	DSN  string
}

type FirmwareConfig struct {
	StoragePath string
}

type MQTTConfig struct {
	Enabled  bool
	Broker   string
	Port     int
	Username string
	Password string
	ClientID string
}

func Load() *Config {
	v := viper.New()

	v.SetDefault("server.port", "8080")
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.path", "./data/lightota.db")
	v.SetDefault("firmware.storage_path", "./firmwares")
	v.SetDefault("mqtt.enabled", true)
	v.SetDefault("mqtt.broker", "localhost")
	v.SetDefault("mqtt.port", 1883)
	v.SetDefault("mqtt.username", "")
	v.SetDefault("mqtt.password", "")
	v.SetDefault("mqtt.client_id", "lightota-server")

	v.AutomaticEnv()

	cfg := &Config{}
	cfg.Server.Port = v.GetString("server.port")
	cfg.Database.Type = v.GetString("database.type")
	cfg.Database.Path = v.GetString("database.path")
	cfg.Database.DSN = v.GetString("database.dsn")
	cfg.Firmware.StoragePath = v.GetString("firmware.storage_path")
	cfg.MQTT.Enabled = v.GetBool("mqtt.enabled")
	cfg.MQTT.Broker = v.GetString("mqtt.broker")
	cfg.MQTT.Port = v.GetInt("mqtt.port")
	cfg.MQTT.Username = v.GetString("mqtt.username")
	cfg.MQTT.Password = v.GetString("mqtt.password")
	cfg.MQTT.ClientID = v.GetString("mqtt.client_id")

	log.Printf("Config loaded: port=%s, db=%s, mqtt=%v", cfg.Server.Port, cfg.Database.Path, cfg.MQTT.Enabled)

	return cfg
}
