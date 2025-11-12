package main

import (
	"context"
	"fmt"
	"time"

	"go.viam.com/rdk/components/board"

	"go.viam.com/rdk/robot/client"
)

const Y_AO_PIN = "y"
const X_AO_PIN = "x"
const SELECT_PIN = "7"

func ADC_joystick(machine *client.RobotClient) error {
	pi, err := board.FromProvider(machine, "pi")
	if err != nil {
		return fmt.Errorf("failed to get board provider 'pi': %w", err)
	}

	y_ao_pin, err := pi.AnalogByName(Y_AO_PIN)
	if err != nil {
		return fmt.Errorf("failed to get analog pin by name %q: %w", Y_AO_PIN, err)
	}
	x_ao_pin, err := pi.AnalogByName(X_AO_PIN)
	if err != nil {
		return fmt.Errorf("failed to get analog pin by name %q: %w", X_AO_PIN, err)
	}
	selectPin, err := pi.GPIOPinByName(SELECT_PIN)
	if err != nil {
		return fmt.Errorf("failed to get GPIO pin by name %q: %w", SELECT_PIN, err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		yVal, err := y_ao_pin.Read(context.Background(), map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("failed to read Y AO pin (%q): %w", Y_AO_PIN, err)
		}
		xVal, err := x_ao_pin.Read(context.Background(), map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("failed to read X AO pin (%q): %w", X_AO_PIN, err)
		}
		selectVal, err := selectPin.Get(context.Background(), map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("failed to get select pin (%q): %w", SELECT_PIN, err)
		}
		fmt.Println("Select Value:", selectVal)
		fmt.Println("Y AO Value:", yVal.Value)
		fmt.Println("X AO Value:", xVal.Value)
	}

	return nil
}
