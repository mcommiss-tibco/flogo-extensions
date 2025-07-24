package try

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
	ctx := test.NewActivityInitContext(map[string]interface{}{
		"metricType": true,
	}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, a)
}

func TestEval_MetricTypeEnabled(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{
		"metricType": true,
	}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("inputString", "test input")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := tc.GetOutput("outputString")
	assert.Equal(t, "Processed (MetricType enabled): test input", output)
}

func TestEval_MetricTypeDisabled(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{
		"metricType": false,
	}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("inputString", "test input")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := tc.GetOutput("outputString")
	assert.Equal(t, "Passed through (MetricType disabled): test input", output)
}

func TestEval_EmptyInput(t *testing.T) {
	ctx := test.NewActivityInitContext(map[string]interface{}{
		"metricType": true,
	}, nil)

	a, err := New(ctx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(a.Metadata())
	tc.SetInput("inputString", "")

	done, err := a.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := tc.GetOutput("outputString")
	assert.Equal(t, "Processed (MetricType enabled): ", output)
}
