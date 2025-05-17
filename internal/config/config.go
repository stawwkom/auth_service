package config

import (
	"fmt"
	"github.com/joho/godotenv" // Для загрузки переменных из .env файла
	"github.com/spf13/viper"   // Для управления конфигурацией из yaml и env
	"log"
	"os"
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
}

// Глобальная переменная для хранения загруженной конфигурации
var Cfg *Config

// Load загружает конфигурацию из файлов и переменных окружения
func Load() error {
	// Загружаем .env файл (если он есть), чтобы переменные окружения были доступны.
	_ = godotenv.Load("../../.env")

	// Определяем текущую среду (например, dev, prod, staging) из переменной окружения APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local" // Если не задана, используем "local" по умолчанию
	}

	// Настраиваем viper для чтения файла конфигурации с именем среды (например, local.yaml)
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")          // по типу конфига и считаного env установливаем конфиг
	viper.AddConfigPath("../../configs") // Папка, где хранятся конфиги

	viper.AutomaticEnv() // Подключаем автоматическое чтение env

	// Считываем конфиг из файла
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Распарсиваем конфигурацию в структуру Config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Сохраняем глобально для удобного доступа из других пакетов
	Cfg = &cfg

	// Логируем успешную загрузку конфигурации
	log.Printf("Loaded config for %s environment", env)
	return nil
}
