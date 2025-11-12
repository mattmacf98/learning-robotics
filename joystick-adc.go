package learningrobotics

import (
	"context"
	"errors"
	"fmt"

	"go.viam.com/rdk/components/board"
	sensor "go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	JoystickAdc = resource.NewModel("mattmacf", "learning-robotics", "joystick-adc")
)

func init() {
	resource.RegisterComponent(sensor.API, JoystickAdc,
		resource.Registration[sensor.Sensor, *JoystickAdcConfig]{
			Constructor: newJoystickAdcJoystickAdc,
		},
	)
}

type JoystickAdcConfig struct {
	YAOPin        string `json:"y_ao_pin"`
	XAOPin        string `json:"x_ao_pin"`
	SelectGPIOPin string `json:"select_gpio_pin"`
	BoardName     string `json:"board_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *JoystickAdcConfig) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.YAOPin == "" {
		return nil, nil, errors.New("y_ao_pin is required")
	}
	if cfg.XAOPin == "" {
		return nil, nil, errors.New("x_ao_pin is required")
	}
	if cfg.SelectGPIOPin == "" {
		return nil, nil, errors.New("select_gpio_pin is required")
	}
	if cfg.BoardName == "" {
		return nil, nil, errors.New("board_name is required")
	}
	return nil, nil, nil
}

type joystickAdcJoystickAdc struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *JoystickAdcConfig

	cancelCtx  context.Context
	cancelFunc func()

	yAOPin        board.Analog
	xAOPin        board.Analog
	selectGPIOPin board.GPIOPin
}

func newJoystickAdcJoystickAdc(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	conf, err := resource.NativeConfig[*JoystickAdcConfig](rawConf)
	if err != nil {
		return nil, err
	}

	return NewJoystickAdc(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewJoystickAdc(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *JoystickAdcConfig, logger logging.Logger) (sensor.Sensor, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	piBoard, err := board.FromProvider(deps, conf.BoardName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	yAOPin, err := piBoard.AnalogByName(conf.YAOPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	xAOPin, err := piBoard.AnalogByName(conf.XAOPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	selectGPIOPin, err := piBoard.GPIOPinByName(conf.SelectGPIOPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	s := &joystickAdcJoystickAdc{
		name:          name,
		logger:        logger,
		cfg:           conf,
		cancelCtx:     cancelCtx,
		cancelFunc:    cancelFunc,
		yAOPin:        yAOPin,
		xAOPin:        xAOPin,
		selectGPIOPin: selectGPIOPin,
	}

	return s, nil
}

func (s *joystickAdcJoystickAdc) Name() resource.Name {
	return s.name
}

func (s *joystickAdcJoystickAdc) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	yValue, err := s.yAOPin.Read(ctx, extra)
	if err != nil {
		return nil, err
	}
	xValue, err := s.xAOPin.Read(ctx, extra)
	if err != nil {
		return nil, err
	}
	selectValue, err := s.selectGPIOPin.Get(ctx, extra)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"y": yValue.Value, "x": xValue.Value, "select": !selectValue}, nil
}

func (s *joystickAdcJoystickAdc) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *joystickAdcJoystickAdc) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
