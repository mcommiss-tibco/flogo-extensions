package try

import (
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Settings represents the activity settings
type Settings struct {
	MetricType bool `md:"metricType,required"`
}

// Input represents the input data structure
type Input struct {
	InputString string `md:"inputString"`
}

// Output represents the output data structure
type Output struct {
	OutputString string `md:"outputString"`
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
	logger.Debugf("Creating Try Activity with settings: %+v", settings)

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
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}

	a.logger.Debugf("Try Activity executing with input: %+v", input)

	// Process the input based on the metricType setting
	var outputString string
	if a.settings.MetricType {
		// If metricType is true, transform the input string
		outputString = fmt.Sprintf("Processed (MetricType enabled): %s", input.InputString)
		a.logger.Infof("MetricType is enabled, processing input string")
	} else {
		// If metricType is false, just pass through the input
		outputString = fmt.Sprintf("Passed through (MetricType disabled): %s", input.InputString)
		a.logger.Infof("MetricType is disabled, passing through input string")
	}

	output := &Output{
		OutputString: outputString,
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	a.logger.Debugf("Try Activity completed with output: %+v", output)

	return true, nil
}
