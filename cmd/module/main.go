package main

import (
	"learningrobotics"

	sensor "go.viam.com/rdk/components/sensor"
	sw "go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

func main() {
	// ModularMain can take multiple APIModel arguments, if your module implements multiple models.
	module.ModularMain(
		resource.APIModel{generic.API, learningrobotics.RgbLed},
		resource.APIModel{generic.API, learningrobotics.LightSwitch},
		resource.APIModel{sensor.API, learningrobotics.UltrasonicSensor},
		resource.APIModel{sensor.API, learningrobotics.JoystickAdc},
		resource.APIModel{sw.API, learningrobotics.RgbPq},
		resource.APIModel{generic.API, learningrobotics.PriorityQueueSwitch},
		resource.APIModel{generic.API, learningrobotics.EventSystem},
	)
}
