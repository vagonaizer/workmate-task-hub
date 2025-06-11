package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// DBType — тип хранилища.
type DBType int

const (
	DBInMemory DBType = iota
	DBPostgres
)

func (t DBType) String() string {
	switch t {
	case DBInMemory:
		return "inmemory"
	case DBPostgres:
		return "postgres"
	default:
		return "unknown"
	}
}

// EnvType — тип окружения.
type EnvType string

const (
	EnvDevelopment EnvType = "dev"
	EnvProduction  EnvType = "prod"
	EnvTest        EnvType = "test"
)

// ServerConfig — конфиг HTTP-сервера.
type ServerConfig struct {
	Port string
}

// LoggerConfig — конфиг логгера.
type LoggerConfig struct {
	Level string
}

// DBConfig — конфиг для хранилища.
type DBConfig struct {
	Type DBType
	DSN  string
}

// TaskConfig — конфиг для задач.
type TaskConfig struct {
	DefaultDuration time.Duration
}

// AppConfig — основной конфиг приложения.
// -- Содержит конфиги для всех компонентов приложения через композицию.
//
// Пока не обнаружил минусов подобного подхода, таким образом
// мы избегаем огромной мешанины полей в конфиге.
type AppConfig struct {
	AppName    string
	AppVersion string
	Env        EnvType
	Server     ServerConfig
	Logger     LoggerConfig
	DB         DBConfig
	Task       TaskConfig
}

// ParseDBType — парсинг типа хранилища из строки.
func ParseDBType(s string) DBType {
	switch strings.ToLower(s) {
	case "inmemory":
		return DBInMemory
	case "postgres":
		return DBPostgres
	default:
		return DBInMemory
	}
}

// Методы для проверки окружения
func (c *AppConfig) IsDevelopment() bool {
	return c.Env == EnvDevelopment
}

func (c *AppConfig) IsProduction() bool {
	return c.Env == EnvProduction
}

func (c *AppConfig) IsTest() bool {
	return c.Env == EnvTest
}

func (c *AppConfig) IsDebug() bool {
	return c.Env == EnvDevelopment || c.Env == EnvTest
}

// LoadConfig — инициализация и загрузка конфига через viper.
func LoadConfig() *AppConfig {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Дефолты
	viper.SetDefault("env", string(EnvDevelopment))
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("db.type", "inmemory")
	viper.SetDefault("db.dsn", "")
	viper.SetDefault("task.defaultduration", "5m")
	viper.SetDefault("appname", "task-hub")
	viper.SetDefault("appversion", "1.0.0")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found: %v (using env/defaults)", err)
	}

	return &AppConfig{
		AppName:    viper.GetString("appname"),
		AppVersion: viper.GetString("appversion"),
		Env:        EnvType(viper.GetString("env")),
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
		},
		Logger: LoggerConfig{
			Level: viper.GetString("logger.level"),
		},
		DB: DBConfig{
			Type: ParseDBType(viper.GetString("db.type")),
			DSN:  viper.GetString("db.dsn"),
		},
		Task: TaskConfig{
			DefaultDuration: viper.GetDuration("task.defaultduration"),
		},
	}
}
