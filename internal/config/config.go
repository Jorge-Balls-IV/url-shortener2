package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	env "github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	ServerHTTP  `yaml:"http_server" env-required:"true"`
}

type ServerHTTP struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"SERVER_HTTP_PASSWORD"`
	// Тэг env устанавливает название для переменной окружения, в которой будет храниться папроль в тайном мете гитхаба
}

//пишем функцию чтения файла с конфигом и заполнения структуры
//Начинается с префикса Must, потому что мы не возвращаем ошибку, а сразу паникуем

// Сначала подгружаем путь к конфигу из переменной окружения
func MustLoadEnv() {
	err := env.Load("./path.env")
	if err != nil {
		log.Fatalf("cannot load path.env: %s", err)
	}
}

func MustLoad() *Config {

	MustLoadEnv()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("empty config path")
	}

	//Проверяем, существет ли файл конфигурации
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist on path: %s", configPath)
	}

	var cnfg Config

	if err := cleanenv.ReadConfig(configPath, &cnfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cnfg
}
