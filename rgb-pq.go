package learningrobotics

import (
	"context"
	"errors"
	"fmt"

	"go.viam.com/rdk/components/board"
	sw "go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	RgbPq = resource.NewModel("mattmacf", "learning-robotics", "rgb-pq")
)

func init() {
	resource.RegisterComponent(sw.API, RgbPq,
		resource.Registration[sw.Switch, *RGBPQConfig]{
			Constructor: newRgbPqRgbPq,
		},
	)
}

type RGBPQConfig struct {
	RedPin    string `json:"red_pin"`
	GreenPin  string `json:"green_pin"`
	BluePin   string `json:"blue_pin"`
	BoardName string `json:"board_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *RGBPQConfig) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
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

type learningRoboticsRgbPq struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *RGBPQConfig

	cancelCtx  context.Context
	cancelFunc func()

	redPin   board.GPIOPin
	greenPin board.GPIOPin
	bluePin  board.GPIOPin
	position uint32
}

func newRgbPqRgbPq(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sw.Switch, error) {
	conf, err := resource.NativeConfig[*RGBPQConfig](rawConf)
	if err != nil {
		return nil, err
	}

	return NewRgbPq(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewRgbPq(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *RGBPQConfig, logger logging.Logger) (sw.Switch, error) {

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

	s := &learningRoboticsRgbPq{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		redPin:     redPin,
		greenPin:   greenPin,
		bluePin:    bluePin,
		position:   0,
	}
	return s, nil
}

func (s *learningRoboticsRgbPq) Name() resource.Name {
	return s.name
}

// SetPosition sets the switch to the specified position.
// Position must be within the valid range for the switch type.
func (s *learningRoboticsRgbPq) SetPosition(ctx context.Context, position uint32, extra map[string]interface{}) error {
	s.position = position

	s.bluePin.Set(ctx, false, extra)
	s.greenPin.Set(ctx, false, extra)
	s.redPin.Set(ctx, false, extra)

	switch s.position {
	case 1:
		err := s.redPin.Set(ctx, true, extra)
		if err != nil {
			return err
		}
	case 2:
		err := s.greenPin.Set(ctx, true, extra)
		if err != nil {
			return err
		}
	case 3:
		err := s.bluePin.Set(ctx, true, extra)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetPosition returns the current position of the switch.
func (s *learningRoboticsRgbPq) GetPosition(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return s.position, nil
}

// GetNumberOfPositions returns the total number of valid positions for this switch, along with their labels.
// Labels should either be nil, empty, or the same length has the number of positions.
func (s *learningRoboticsRgbPq) GetNumberOfPositions(ctx context.Context, extra map[string]interface{}) (uint32, []string, error) {
	return 4, []string{"off", "red", "green", "blue"}, nil
}

func (s *learningRoboticsRgbPq) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *learningRoboticsRgbPq) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
