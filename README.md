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
