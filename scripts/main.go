package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/robot/client"
	"go.viam.com/utils/rpc"
)

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %q required but not set", key)
	}
	return val
}

func loadEnvVars() (host, apiKey, apiVal string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	host = mustGetEnv("ROBOT_HOST")
	apiKey = mustGetEnv("API_KEY_NAME")
	apiVal = mustGetEnv("API_KEY_VAL")
	return
}

func main() {
	host, apiKeyName, apiKeyVal := loadEnvVars()
	logger := logging.NewDebugLogger("client")
	machine, err := client.New(
		context.Background(),
		host,
		logger,
		client.WithDialOptions(rpc.WithEntityCredentials(

			apiKeyName,
			rpc.Credentials{
				Type: rpc.CredentialsTypeAPIKey,

				Payload: apiKeyVal,
			})),
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer machine.Close(context.Background())
	err = RgbPriorityQueue(machine)
	if err != nil {
		logger.Fatal(err)
	}
}
