package main

import (
	"context"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/robot/client"
)

const RED_LED_PIN = "32"
const GREEN_LED_PIN = "33"
const BLUE_LED_PIN = "35"

func RGBLed(machine *client.RobotClient) error {
	pi, err := board.FromRobot(machine, "pi")
	if err != nil {
		return err
	}

	red_pin, err := pi.GPIOPinByName(RED_LED_PIN)
	if err != nil {
		return err
	}
	blue_pin, err := pi.GPIOPinByName(BLUE_LED_PIN)
	if err != nil {
		return err
	}
	green_pin, err := pi.GPIOPinByName(GREEN_LED_PIN)
	if err != nil {
		return err
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
	return nil
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
