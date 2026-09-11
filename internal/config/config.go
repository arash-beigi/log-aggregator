package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv      string
	NATSUrl     string
	NATSSubject string
	MongoURI    string
	MongoDBName string
	MongoCol    string
	WorkerCount int
	BatchSize   int
}

func LoadConfig() *Config {
	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		NATSUrl:     getEnv("NATS_URL", "nats://localhost:4222"),
		NATSSubject: getEnv("NATS_SUBJECT", "logs.>"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://root:rootpassword@localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "log_aggregator"),
		MongoCol:    getEnv("MONGO_COLLECTION", "logs"),
		WorkerCount: getEnvAsInt("WORKER_COUNT", 5),
		BatchSize:   getEnvAsInt("BATCH_SIZE", 100),
	}
}
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
