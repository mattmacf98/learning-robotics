package main

import (
	"context"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/robot/client"
)

const LIGHT_OUTPUT_PIN = "8"
const ON_BUTTON_INPUT_PIN = "38"
const OFF_BUTTON_INPUT_PIN = "40"

func LightSwitch(machine *client.RobotClient) error {
	pi, err := board.FromProvider(machine, "pi")
	if err != nil {
		return err
	}
	lightOutputPin, err := pi.GPIOPinByName(LIGHT_OUTPUT_PIN)
	if err != nil {
		return err
	}

	onButtonInputPin, err := pi.GPIOPinByName(ON_BUTTON_INPUT_PIN)
	if err != nil {
		return err
	}
	offButtonInputPin, err := pi.GPIOPinByName(OFF_BUTTON_INPUT_PIN)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second / 30) // 30 frames per second
	defer ticker.Stop()

	for range ticker.C {
		lightOnValue, err := onButtonInputPin.Get(context.Background(), map[string]interface{}{})
		if err != nil {
			return err
		}
		if !lightOnValue {
			lightOutputPin.Set(context.Background(), true, map[string]interface{}{})
		}

		lightOffValue, err := offButtonInputPin.Get(context.Background(), map[string]interface{}{})
		if err != nil {
			return err
		}
		if !lightOffValue {
			lightOutputPin.Set(context.Background(), false, map[string]interface{}{})
		}
	}

	return nil
}
