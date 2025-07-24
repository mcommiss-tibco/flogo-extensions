package ssh

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"golang.org/x/crypto/ssh"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Settings represents the activity settings (empty for this activity)
type Settings struct {
}

// Input represents the input data structure
type Input struct {
	SSHServername string `md:"sshServername"`
	SSHServerPort int    `md:"sshServerPort"`
	SSHUsername   string `md:"sshUsername"`
	SSHPassword   string `md:"sshPassword"`
	SSHCommand    string `md:"sshCommand"`
}

// FromMap implements data.StructValue.FromMap
func (i *Input) FromMap(values map[string]interface{}) error {
	if val, exists := values["sshServername"]; exists {
		if str, ok := val.(string); ok {
			i.SSHServername = str
		}
	}
	if val, exists := values["sshServerPort"]; exists {
		if port, ok := val.(int); ok {
			i.SSHServerPort = port
		} else if portFloat, ok := val.(float64); ok {
			i.SSHServerPort = int(portFloat)
		}
	}
	if val, exists := values["sshUsername"]; exists {
		if str, ok := val.(string); ok {
			i.SSHUsername = str
		}
	}
	if val, exists := values["sshPassword"]; exists {
		if str, ok := val.(string); ok {
			i.SSHPassword = str
		}
	}
	if val, exists := values["sshCommand"]; exists {
		if str, ok := val.(string); ok {
			i.SSHCommand = str
		}
	}
	return nil
}

// ToMap implements data.StructValue.ToMap
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"sshServername": i.SSHServername,
		"sshServerPort": i.SSHServerPort,
		"sshUsername":   i.SSHUsername,
		"sshPassword":   i.SSHPassword,
		"sshCommand":    i.SSHCommand,
	}
}

// Output represents the output data structure
type Output struct {
	ResultCode int    `md:"resultCode"`
	ResultText string `md:"resultText"`
}

// FromMap implements data.StructValue.FromMap
func (o *Output) FromMap(values map[string]interface{}) error {
	if val, exists := values["resultCode"]; exists {
		if code, ok := val.(int); ok {
			o.ResultCode = code
		} else if codeFloat, ok := val.(float64); ok {
			o.ResultCode = int(codeFloat)
		}
	}
	if val, exists := values["resultText"]; exists {
		if str, ok := val.(string); ok {
			o.ResultText = str
		}
	}
	return nil
}

// ToMap implements data.StructValue.ToMap
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"resultCode": o.ResultCode,
		"resultText": o.ResultText,
	}
}

// Activity is the main activity struct
type Activity struct {
	settings *Settings
	logger   log.Logger
}

// New creates a new instance of the activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	logger := ctx.Logger()
	logger.Debugf("Creating SSH Activity")

	act := &Activity{
		settings: settings,
		logger:   logger,
	}

	return act, nil
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval executes the activity
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	// Get input values using the proper Flogo API
	sshServername := ""
	sshServerPort := 22
	sshUsername := ""
	sshPassword := ""
	sshCommand := ""

	if val := ctx.GetInput("sshServername"); val != nil {
		if str, ok := val.(string); ok {
			sshServername = str
		}
	}

	if val := ctx.GetInput("sshServerPort"); val != nil {
		if port, ok := val.(int); ok {
			sshServerPort = port
		} else if portFloat, ok := val.(float64); ok {
			sshServerPort = int(portFloat)
		}
	}

	if val := ctx.GetInput("sshUsername"); val != nil {
		if str, ok := val.(string); ok {
			sshUsername = str
		}
	}

	if val := ctx.GetInput("sshPassword"); val != nil {
		if str, ok := val.(string); ok {
			sshPassword = str
		}
	}

	if val := ctx.GetInput("sshCommand"); val != nil {
		if str, ok := val.(string); ok {
			sshCommand = str
		}
	}

	a.logger.Debugf("SSH Activity executing with server: %s:%d, user: %s", sshServername, sshServerPort, sshUsername)

	// Validate required inputs
	if sshServername == "" {
		err = fmt.Errorf("sshServername is required")
		a.logger.Error(err.Error())
		ctx.SetOutput("resultCode", -1)
		ctx.SetOutput("resultText", err.Error())
		return true, nil
	}

	if sshUsername == "" {
		err = fmt.Errorf("sshUsername is required")
		a.logger.Error(err.Error())
		ctx.SetOutput("resultCode", -1)
		ctx.SetOutput("resultText", err.Error())
		return true, nil
	}

	if sshPassword == "" {
		err = fmt.Errorf("sshPassword is required")
		a.logger.Error(err.Error())
		ctx.SetOutput("resultCode", -1)
		ctx.SetOutput("resultText", err.Error())
		return true, nil
	}

	if sshCommand == "" {
		err = fmt.Errorf("sshCommand is required")
		a.logger.Error(err.Error())
		ctx.SetOutput("resultCode", -1)
		ctx.SetOutput("resultText", err.Error())
		return true, nil
	}

	// Execute SSH command
	resultCode, resultText, err := a.executeSSHCommand(sshServername, sshServerPort, sshUsername, sshPassword, sshCommand)
	if err != nil {
		a.logger.Errorf("SSH command execution failed: %v", err)
		ctx.SetOutput("resultCode", -1)
		ctx.SetOutput("resultText", fmt.Sprintf("SSH Error: %v", err))
		return true, nil
	}

	// Set outputs
	ctx.SetOutput("resultCode", resultCode)
	ctx.SetOutput("resultText", resultText)

	a.logger.Debugf("SSH Activity completed with result code: %d", resultCode)

	return true, nil
}

// executeSSHCommand connects to SSH server and executes the command
func (a *Activity) executeSSHCommand(hostname string, port int, username, password, command string) (int, string, error) {
	// Create SSH client configuration
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, you should verify host keys
		Timeout:         30 * time.Second,
	}

	// Connect to SSH server
	address := net.JoinHostPort(hostname, strconv.Itoa(port))
	a.logger.Debugf("Connecting to SSH server: %s", address)

	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return -1, "", fmt.Errorf("failed to connect to SSH server: %v", err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return -1, "", fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	a.logger.Debugf("Executing SSH command: %s", command)

	// Execute the command
	output, err := session.CombinedOutput(command)
	outputText := strings.TrimSpace(string(output))

	if err != nil {
		// Check if it's an exit error (command returned non-zero exit code)
		if exitError, ok := err.(*ssh.ExitError); ok {
			exitCode := exitError.ExitStatus()
			a.logger.Debugf("SSH command completed with exit code: %d", exitCode)
			return exitCode, outputText, nil
		}
		// Other SSH errors
		return -1, outputText, fmt.Errorf("SSH command execution error: %v", err)
	}

	// Command executed successfully (exit code 0)
	a.logger.Debugf("SSH command completed successfully")
	return 0, outputText, nil
}
