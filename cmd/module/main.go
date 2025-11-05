package main

import (
	"learningrobotics"

	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

func main() {
	// ModularMain can take multiple APIModel arguments, if your module implements multiple models.
	module.ModularMain(
		resource.APIModel{generic.API, learningrobotics.RgbLed},
		resource.APIModel{generic.API, learningrobotics.LightSwitch},
	)
}
