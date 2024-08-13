package app

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddr  string
	ServerPort uint16
}

func LoadConfig() Config {

	// Configurar o viper
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	// Lê o arquivo de configuração
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Erro ao ler o arquivo de configuração: %v", err)
	}

	// Obter valores do arquivo de configuração
	redisAddress := viper.GetString("redis_address")
	serverPort := viper.GetUint16("server_port")
	
	// Cria uma instância de Config p/retornar
	cfg := Config{
		RedisAddr:  redisAddress,
		ServerPort: serverPort,
	}

	return cfg
}