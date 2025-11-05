package learningrobotics

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

var (
	RgbLed = resource.NewModel("mattmacf", "learning-robotics", "rgb-led")
)

func init() {
	resource.RegisterService(generic.API, RgbLed,
		resource.Registration[resource.Resource, *RGBLedConfig]{
			Constructor: newLearningRoboticsRgbLed,
		},
	)
}

type RGBLedConfig struct {
	RedPin    string `json:"red_pin"`
	GreenPin  string `json:"green_pin"`
	BluePin   string `json:"blue_pin"`
	BoardName string `json:"board_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *RGBLedConfig) Validate(path string) ([]string, []string, error) {
	if cfg.RedPin == "" {
		return nil, nil, errors.New("red_pin is required")
	}
	if cfg.GreenPin == "" {
		return nil, nil, errors.New("green_pin is required")
	}
	if cfg.BluePin == "" {
		return nil, nil, errors.New("blue_pin is required")
	}
	if cfg.BoardName == "" {
		return nil, nil, errors.New("board_name is required")
	}
	return nil, nil, nil
}

type learningRoboticsRgbLed struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *RGBLedConfig

	cancelCtx  context.Context
	cancelFunc func()

	redPin   board.GPIOPin
	greenPin board.GPIOPin
	bluePin  board.GPIOPin
}

func newLearningRoboticsRgbLed(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*RGBLedConfig](rawConf)
	if err != nil {
		return nil, err
	}

	return NewRgbLed(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewRgbLed(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *RGBLedConfig, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	board, err := board.FromProvider(deps, conf.BoardName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	redPin, err := board.GPIOPinByName(conf.RedPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	greenPin, err := board.GPIOPinByName(conf.GreenPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	bluePin, err := board.GPIOPinByName(conf.BluePin)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	s := &learningRoboticsRgbLed{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		redPin:     redPin,
		greenPin:   greenPin,
		bluePin:    bluePin,
	}
	return s, nil
}

func (s *learningRoboticsRgbLed) Name() resource.Name {
	return s.name
}

func (s *learningRoboticsRgbLed) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if partyData, ok := cmd["party_mode"]; ok {

		occurences, ok := partyData.(map[string]any)["occurences"]
		if !ok {
			return nil, fmt.Errorf("occurences is required")
		}
		occurStr, ok := occurences.(string)
		if !ok {
			return nil, fmt.Errorf("occurences must be a string")
		}
		occurInt, err := strconv.Atoi(occurStr)
		if err != nil {
			return nil, fmt.Errorf("occurences must be an integer string: %v", err)
		}
		err = s.partyMode(occurInt)
		if err != nil {
			return nil, err
		}
		return map[string]any{"status": "success"}, nil
	}

	if _, ok := cmd["turn_off"]; ok {
		err := s.turnOff()
		if err != nil {
			return nil, err
		}
		return map[string]any{"status": "success"}, nil
	}

	if _, ok := cmd["make_red"]; ok {
		err := s.makeRed()
		if err != nil {
			return nil, err
		}
		return map[string]any{"status": "success"}, nil
	}

	if _, ok := cmd["make_green"]; ok {
		err := s.makeGreen()
		if err != nil {
			return nil, err
		}
		return map[string]any{"status": "success"}, nil
	}

	if _, ok := cmd["make_blue"]; ok {
		err := s.makeBlue()
		if err != nil {
			return nil, err
		}
		return map[string]any{"status": "success"}, nil
	}

	return nil, fmt.Errorf("Unknown command: %v", cmd)
}

func (s *learningRoboticsRgbLed) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

func (s *learningRoboticsRgbLed) makeRed() error {
	err := s.bluePin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.greenPin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.redPin.Set(s.cancelCtx, true, map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (s *learningRoboticsRgbLed) makeGreen() error {
	err := s.bluePin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.greenPin.Set(s.cancelCtx, true, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.redPin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (s *learningRoboticsRgbLed) makeBlue() error {

	err := s.bluePin.Set(s.cancelCtx, true, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.greenPin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.redPin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (s *learningRoboticsRgbLed) turnOff() error {
	err := s.bluePin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.greenPin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = s.redPin.Set(s.cancelCtx, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

func (s *learningRoboticsRgbLed) partyMode(occurences int) error {
	for range occurences {
		err := s.makeRed()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		err = s.makeGreen()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
		err = s.makeBlue()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
	err := s.turnOff()
	if err != nil {
		return err
	}
	return nil
}
