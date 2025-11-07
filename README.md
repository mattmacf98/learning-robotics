# Module learning-robotics

This module provides some models for learning how to interact with a Raspberry Pi using Viam. It is meant to function as a replacement for popular ARduino courses.

## Model mattmacf:learning-robotics:rgb-led

This model represents an RGB LED driven by three GPIO pins (red, green, and blue) on a supported board. It supports commands to turn the LED off, set it to primary colors, and activate a party mode (color cycling) a specific number of times.

### Configuration

The following attribute template can be used to configure this model:

```json
{
  "board": "<string>",
  "red_pin": "<string>",
  "green_pin": "<string>",
  "blue_pin": "<string>"
}
```

#### Attributes

The following attributes are required for this model:

| Name        | Type   | Inclusion | Description                            |
| ----------- | ------ | --------- | -------------------------------------- |
| `board`     | string | Required  | The name of the board interface to use |
| `red_pin`   | string | Required  | The pin name for the red LED channel   |
| `green_pin` | string | Required  | The pin name for the green LED channel |
| `blue_pin`  | string | Required  | The pin name for the blue LED channel  |

#### Example Configuration

```json
{
  "board": "my-board",
  "red_pin": "32",
  "green_pin": "33",
  "blue_pin": "35"
}
```

### DoCommand

The model implements DoCommand for runtime LED color control and party mode effects.

#### Example DoCommands

Turn the LED off:

```json
{
  "turn_off": true
}
```

Set the LED to red:

```json
{
  "make_red": true
}
```

Set the LED to green:

```json
{
  "make_green": true
}
```

Set the LED to blue:

```json
{
  "make_blue": true
}
```

Activate party mode (cycle through colors a set number of times):

```json
{
  "party_mode": {
    "occurences": "5"
  }
}
```

Note: `occurences` should be a string representing an integer number of cycles.

## Model mattmacf:learning-robotics:light-switch

This model represents a light switch system driven by two buttons (on and off) and one output pin to control a light. The system continuously monitors button inputs and controls the light accordingly. It is designed as a generic service to teach basic input/output control with physical buttons.

### Configuration

The following attribute template can be used to configure this model:

```json
{
  "board_name": "<string>",
  "light_output_pin": "<string>",
  "on_button_input_pin": "<string>",
  "off_button_input_pin": "<string>"
}
```

#### Attributes

The following attributes are required for this model:

| Name                   | Type   | Inclusion | Description                                      |
| ---------------------- | ------ | --------- | ------------------------------------------------ |
| `board_name`           | string | Required  | The name of the board interface to use           |
| `light_output_pin`     | string | Required  | The pin name that controls the light output      |
| `on_button_input_pin`  | string | Required  | The pin name connected to the "on" button input  |
| `off_button_input_pin` | string | Required  | The pin name connected to the "off" button input |

#### Example Configuration

```json
{
  "board_name": "my-board",
  "light_output_pin": "32",
  "on_button_input_pin": "36",
  "off_button_input_pin": "38"
}
```

### Behavior

This model runs continuously in the background, polling the button inputs at 30 frames per second. When the on button is pressed (low signal), the light output is set to high. When the off button is pressed (low signal), the light output is set to low.

## Model mattmacf:learning-robotics:ultrasonic-sensor

This model represents an ultrasonic distance sensor (such as the HC-SR04) that measures distance to objects using ultrasonic pulses. The sensor uses a trigger pin to emit sound pulses and an echo interrupt pin to detect the returning echo, calculating distance based on the time difference.

### Configuration

The following attribute template can be used to configure this model:

```json
{
  "board_name": "<string>",
  "trigger_pin": "<string>",
  "echo_interrupt_pin": "<string>"
}
```

#### Attributes

The following attributes are required for this model:

| Name                 | Type   | Inclusion | Description                                         |
| -------------------- | ------ | --------- | --------------------------------------------------- |
| `board_name`         | string | Required  | The name of the board interface to use              |
| `trigger_pin`        | string | Required  | The pin name used to trigger ultrasonic pulses      |
| `echo_interrupt_pin` | string | Required  | The digital interrupt pin name for detecting echoes |

#### Example Configuration

```json
{
  "board_name": "my-board",
  "trigger_pin": "16",
  "echo_interrupt_pin": "18"
}
```

### Readings

The sensor implements the standard `Readings()` method which returns the distance to the nearest object in meters. The distance is calculated using the time between sending the ultrasonic pulse and receiving its echo, based on the speed of sound (343 m/s).

#### Example Response

```json
{
  "distance": 0.523
}
```

The `distance` value is in meters. For example, 0.523 meters equals approximately 52.3 centimeters.
