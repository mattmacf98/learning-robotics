package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/robot/client"
	"go.viam.com/utils/rpc"
)

const RED_LED_PIN = "32"
const GREEN_LED_PIN = "33"
const BLUE_LED_PIN = "35"

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
	logger.Info("Resources:")
	logger.Info(machine.ResourceNames())

	// Note that the pin supplied is a placeholder. Please change this to a valid pin.
	// pi
	pi, err := board.FromRobot(machine, "pi")
	if err != nil {
		logger.Error(err)
		return
	}

	red_pin, err := pi.GPIOPinByName(RED_LED_PIN)
	if err != nil {
		logger.Error(err)
		return
	}
	blue_pin, err := pi.GPIOPinByName(BLUE_LED_PIN)
	if err != nil {
		logger.Error(err)
		return
	}
	green_pin, err := pi.GPIOPinByName(GREEN_LED_PIN)
	if err != nil {
		logger.Error(err)
		return
	}

	for range 10 {
		makeRed(red_pin, blue_pin, green_pin)
		time.Sleep(time.Millisecond * 100)
		makeGreen(red_pin, blue_pin, green_pin)
		time.Sleep(time.Millisecond * 100)
		makeBlue(red_pin, blue_pin, green_pin)
		time.Sleep(time.Millisecond * 100)
	}
	turnOff(red_pin, blue_pin, green_pin)
}

func makeRed(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), true, map[string]interface{}{})
	b.Set(context.Background(), false, map[string]interface{}{})
	g.Set(context.Background(), false, map[string]interface{}{})
}

func makeGreen(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), false, map[string]interface{}{})
	b.Set(context.Background(), false, map[string]interface{}{})
	g.Set(context.Background(), true, map[string]interface{}{})
}

func makeBlue(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), false, map[string]interface{}{})
	b.Set(context.Background(), true, map[string]interface{}{})
	g.Set(context.Background(), false, map[string]interface{}{})
}

func turnOff(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), false, map[string]interface{}{})
	b.Set(context.Background(), false, map[string]interface{}{})
	g.Set(context.Background(), false, map[string]interface{}{})
}
