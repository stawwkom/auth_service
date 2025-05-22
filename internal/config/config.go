package config

import (
	"fmt"
	"github.com/joho/godotenv" // Для загрузки переменных из .env файла
	"github.com/spf13/viper"   // Для управления конфигурацией из yaml и env
	"log"
	"os"
	"strconv"
)

// Config — структура, в которую будут загружены все настройки приложения
type Config struct {
	Server struct {
		Host string // Хост для сервера
		Port int    // Порт для сервера
	}
	Log struct {
		Level string // Уровень логирования (debug, info, warn, error)
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
}

// Глобальная переменная для хранения загруженной конфигурации
var Cfg *Config

// Load загружает конфигурацию из файлов и переменных окружения
func Load() error {
	// Загружаем .env файл (если он есть), чтобы переменные окружения были доступны.
	_ = godotenv.Load("../../.env")

	// Определяем текущую среду (например, local, prod, staging) из переменной окружения APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local" // Если не задана, используем "local" по умолчанию
	}

	// Настраиваем viper для чтения файла конфигурации с именем среды (например, local.yaml)
	viper.SetConfigName(env)
	viper.SetConfigType("yaml") // по типу конфига и считаного env установливаем конфиг
	// Добавляем пути к конфигурационным файлам
	viper.AddConfigPath("configs")       // Относительный путь из корня проекта
	viper.AddConfigPath("../configs")    // Относительный путь на уровень выше
	viper.AddConfigPath("../../configs") // Относительный путь на два уровня выше
	viper.AutomaticEnv()                 // Подключаем автоматическое чтение env

	// Считываем конфиг из файла
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Распарсиваем конфигурацию в структуру Config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Переопределяем настройки базы данных из переменных окружения, если они заданы
	if dbName := os.Getenv("PG_DATABASE_NAME"); dbName != "" {
		cfg.Database.DBName = dbName
	}
	if dbUser := os.Getenv("PG_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPassword := os.Getenv("PG_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if dbPort := os.Getenv("PG_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			cfg.Database.Port = port
		}
	}

	// Сохраняем глобально для удобного доступа из других пакетов
	Cfg = &cfg

	// Логируем успешную загрузку конфигурации
	log.Printf("Loaded config for %s environment", env)
	return nil
}
