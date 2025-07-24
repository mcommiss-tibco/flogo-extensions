package ssh

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)
	assert.NotNil(t, act)
}

func TestNew(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, a)
}

func TestEval_MissingServername(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("sshUsername", "testuser")
	tc.SetInput("sshPassword", "testpass")
	tc.SetInput("sshCommand", "echo test")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	resultCode := tc.GetOutput("resultCode")
	resultText := tc.GetOutput("resultText")
	assert.Equal(t, -1, resultCode)
	assert.Contains(t, resultText, "sshServername is required")
}

func TestEval_MissingUsername(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("sshServername", "localhost")
	tc.SetInput("sshPassword", "testpass")
	tc.SetInput("sshCommand", "echo test")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	resultCode := tc.GetOutput("resultCode")
	resultText := tc.GetOutput("resultText")
	assert.Equal(t, -1, resultCode)
	assert.Contains(t, resultText, "sshUsername is required")
}

func TestEval_MissingPassword(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("sshServername", "localhost")
	tc.SetInput("sshUsername", "testuser")
	tc.SetInput("sshCommand", "echo test")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	resultCode := tc.GetOutput("resultCode")
	resultText := tc.GetOutput("resultText")
	assert.Equal(t, -1, resultCode)
	assert.Contains(t, resultText, "sshPassword is required")
}

func TestEval_MissingCommand(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("sshServername", "localhost")
	tc.SetInput("sshUsername", "testuser")
	tc.SetInput("sshPassword", "testpass")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	resultCode := tc.GetOutput("resultCode")
	resultText := tc.GetOutput("resultText")
	assert.Equal(t, -1, resultCode)
	assert.Contains(t, resultText, "sshCommand is required")
}

func TestEval_DefaultPort(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("sshServername", "invalid-host")
	tc.SetInput("sshUsername", "testuser")
	tc.SetInput("sshPassword", "testpass")
	tc.SetInput("sshCommand", "echo test")
	// Not setting sshServerPort, should default to 22

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	resultCode := tc.GetOutput("resultCode")
	resultText := tc.GetOutput("resultText")
	assert.Equal(t, -1, resultCode)
	assert.Contains(t, resultText, "SSH Error")
}

func TestEval_CustomPort(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("sshServername", "invalid-host")
	tc.SetInput("sshServerPort", 2222)
	tc.SetInput("sshUsername", "testuser")
	tc.SetInput("sshPassword", "testpass")
	tc.SetInput("sshCommand", "echo test")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	resultCode := tc.GetOutput("resultCode")
	resultText := tc.GetOutput("resultText")
	assert.Equal(t, -1, resultCode)
	assert.Contains(t, resultText, "SSH Error")
}

func TestInput_FromMap(t *testing.T) {
	input := &Input{}
	values := map[string]interface{}{
		"sshServername": "testhost",
		"sshServerPort": 2222,
		"sshUsername":   "testuser",
		"sshPassword":   "testpass",
		"sshCommand":    "echo test",
	}

	err := input.FromMap(values)
	assert.Nil(t, err)
	assert.Equal(t, "testhost", input.SSHServername)
	assert.Equal(t, 2222, input.SSHServerPort)
	assert.Equal(t, "testuser", input.SSHUsername)
	assert.Equal(t, "testpass", input.SSHPassword)
	assert.Equal(t, "echo test", input.SSHCommand)
}

func TestInput_ToMap(t *testing.T) {
	input := &Input{
		SSHServername: "testhost",
		SSHServerPort: 2222,
		SSHUsername:   "testuser",
		SSHPassword:   "testpass",
		SSHCommand:    "echo test",
	}

	values := input.ToMap()
	assert.Equal(t, "testhost", values["sshServername"])
	assert.Equal(t, 2222, values["sshServerPort"])
	assert.Equal(t, "testuser", values["sshUsername"])
	assert.Equal(t, "testpass", values["sshPassword"])
	assert.Equal(t, "echo test", values["sshCommand"])
}

func TestOutput_FromMap(t *testing.T) {
	output := &Output{}
	values := map[string]interface{}{
		"resultCode": 0,
		"resultText": "success",
	}

	err := output.FromMap(values)
	assert.Nil(t, err)
	assert.Equal(t, 0, output.ResultCode)
	assert.Equal(t, "success", output.ResultText)
}

func TestOutput_ToMap(t *testing.T) {
	output := &Output{
		ResultCode: 0,
		ResultText: "success",
	}

	values := output.ToMap()
	assert.Equal(t, 0, values["resultCode"])
	assert.Equal(t, "success", values["resultText"])
}
