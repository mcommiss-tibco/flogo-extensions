# Try Activity

A simple Flogo activity for testing purposes that demonstrates basic input/output processing with configurable behavior.

## Installation

### Flogo CLI
```bash
flogo install github.com/mcommiss-tibco/flogo-extensions/try
```

### Third-party Go dep tools
```bash
go get github.com/mcommiss-tibco/flogo-extensions/try
```

## Schema
Settings, Inputs and Outputs:

```json
{
  "settings": [
    {
      "name": "metricType",
      "type": "boolean",
      "required": true,
      "value": "false"
    }
  ],
  "inputs": [
    {
      "name": "inputString",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "outputString",
      "type": "string"
    }
  ]
}
```

## Settings
| Setting     | Type    | Description |
|:------------|:--------|:------------|
| metricType  | boolean | Boolean setting to toggle processing behavior |

## Inputs
| Input       | Type   | Description |
|:------------|:-------|:------------|
| inputString | string | Input string for testing purposes |

## Outputs
| Output       | Type   | Description |
|:-------------|:-------|:------------|
| outputString | string | Output string for testing purposes |

## Configuration Examples
### Simple
Configure the activity to pass through input unchanged:
```json
{
  "id": "try_activity",
  "name": "Try Activity",
  "activity": {
    "ref": "github.com/mcommiss-tibco/flogo-extensions/try",
    "settings": {
      "metricType": false
    }
  }
}
```

### With Processing
Configure the activity to process the input:
```json
{
  "id": "try_activity",
  "name": "Try Activity",
  "activity": {
    "ref": "github.com/mcommiss-tibco/flogo-extensions/try",
    "settings": {
      "metricType": true
    }
  }
}
```

## Behavior
- When `metricType` is `false`: The activity passes through the input string with minimal modification
- When `metricType` is `true`: The activity processes and transforms the input string

## Testing
Run the unit tests:
```bash
go test ./...
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the BSD 3-Clause License - see the [LICENSE](../LICENSE) file for details.
