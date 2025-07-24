# SSH Command Executor Activity

A Flogo activity for executing SSH commands on remote servers and returning results with exit codes.

## Installation

### Flogo CLI
```bash
flogo install github.com/mcommiss-tibco/flogo-extensions/ssh
```

### Third-party Go dep tools
```bash
go get github.com/mcommiss-tibco/flogo-extensions/ssh
```

## Schema
Settings, Inputs and Outputs:

```json
{
  "inputs": [
    {
      "name": "sshServername",
      "type": "string",
      "required": true
    },
    {
      "name": "sshServerPort",
      "type": "integer",
      "required": false,
      "value": 22
    },
    {
      "name": "sshUsername",
      "type": "string",
      "required": true
    },
    {
      "name": "sshPassword",
      "type": "string",
      "required": true
    },
    {
      "name": "sshCommand",
      "type": "string",
      "required": true
    }
  ],
  "outputs": [
    {
      "name": "resultCode",
      "type": "integer"
    },
    {
      "name": "resultText",
      "type": "string"
    }
  ]
}
```

## Inputs
| Input         | Type    | Required | Description |
|:--------------|:--------|:---------|:------------|
| sshServername | string  | true     | The hostname or IP address of the SSH server |
| sshServerPort | integer | false    | The SSH port number (default: 22) |
| sshUsername   | string  | true     | Username for SSH authentication |
| sshPassword   | string  | true     | Password for SSH authentication |
| sshCommand    | string  | true     | The command to execute on the remote server |

## Outputs
| Output     | Type    | Description |
|:-----------|:--------|:------------|
| resultCode | integer | Exit code from command execution (0 = success, non-zero = error) |
| resultText | string  | Output text from the command execution |

## Configuration Examples

### Simple Command Execution
```json
{
  "id": "ssh_activity",
  "name": "SSH Command Executor",
  "activity": {
    "ref": "github.com/mcommiss-tibco/flogo-extensions/ssh",
    "input": {
      "sshServername": "example.com",
      "sshUsername": "myuser",
      "sshPassword": "mypassword",
      "sshCommand": "ls -la"
    }
  }
}
```

### Custom Port Configuration
```json
{
  "id": "ssh_activity",
  "name": "SSH Command Executor",
  "activity": {
    "ref": "github.com/mcommiss-tibco/flogo-extensions/ssh",
    "input": {
      "sshServername": "example.com",
      "sshServerPort": 2222,
      "sshUsername": "myuser",
      "sshPassword": "mypassword",
      "sshCommand": "df -h"
    }
  }
}
```

### Complex Command with Pipes
```json
{
  "id": "ssh_activity",
  "name": "SSH Command Executor",
  "activity": {
    "ref": "github.com/mcommiss-tibco/flogo-extensions/ssh",
    "input": {
      "sshServername": "example.com",
      "sshUsername": "myuser",
      "sshPassword": "mypassword",
      "sshCommand": "ps aux | grep nginx | wc -l"
    }
  }
}
```

## Features

- **Secure SSH Connection**: Uses golang.org/x/crypto/ssh for secure connections
- **Exit Code Support**: Returns actual command exit codes for proper error handling
- **Combined Output**: Captures both stdout and stderr in a single output
- **Connection Timeout**: 30-second timeout for SSH connections
- **Input Validation**: Validates all required parameters before execution
- **Comprehensive Logging**: Debug and error logging for troubleshooting

## Error Handling

The activity handles various error scenarios:

- **Connection Errors**: Returns `resultCode: -1` for SSH connection failures
- **Authentication Errors**: Returns `resultCode: -1` for authentication failures
- **Command Execution Errors**: Returns the actual exit code from the failed command
- **Missing Parameters**: Returns `resultCode: -1` with descriptive error messages

## Security Considerations

⚠️ **Important Security Notes:**

1. **Host Key Verification**: Currently uses `InsecureIgnoreHostKey()`. For production use, implement proper host key verification.
2. **Password Storage**: Consider using SSH key-based authentication instead of passwords for enhanced security.
3. **Command Injection**: Validate and sanitize commands to prevent injection attacks.
4. **Network Security**: Ensure SSH connections are made over secure networks.

## Testing

Run the unit tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Examples

### Check System Uptime
```json
{
  "sshServername": "server.example.com",
  "sshUsername": "admin",
  "sshPassword": "secure_password",
  "sshCommand": "uptime"
}
```

### Check Disk Space
```json
{
  "sshServername": "server.example.com",
  "sshUsername": "admin",
  "sshPassword": "secure_password",
  "sshCommand": "df -h /"
}
```

### Restart a Service
```json
{
  "sshServername": "server.example.com",
  "sshUsername": "admin",
  "sshPassword": "secure_password",
  "sshCommand": "sudo systemctl restart nginx"
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](../LICENSE) file for details.
