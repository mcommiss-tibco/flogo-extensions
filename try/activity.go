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

// FromMap implements data.StructValue.FromMap
func (i *Input) FromMap(values map[string]interface{}) error {
	if val, exists := values["inputString"]; exists {
		if str, ok := val.(string); ok {
			i.InputString = str
		}
	}
	return nil
}

// ToMap implements data.StructValue.ToMap
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"inputString": i.InputString,
	}
}

// Output represents the output data structure
type Output struct {
	OutputString string `md:"outputString"`
}

// FromMap implements data.StructValue.FromMap
func (o *Output) FromMap(values map[string]interface{}) error {
	if val, exists := values["outputString"]; exists {
		if str, ok := val.(string); ok {
			o.OutputString = str
		}
	}
	return nil
}

// ToMap implements data.StructValue.ToMap
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"outputString": o.OutputString,
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
	// Get input using the proper Flogo API
	inputValue := ctx.GetInput("inputString")
	inputString := ""
	if inputValue != nil {
		if str, ok := inputValue.(string); ok {
			inputString = str
		}
	}

	a.logger.Debugf("Try Activity executing with input: %s", inputString)

	// Process the input based on the metricType setting
	var outputString string
	if a.settings.MetricType {
		// If metricType is true, transform the input string
		outputString = fmt.Sprintf("Processed (MetricType enabled): %s", inputString)
		a.logger.Infof("MetricType is enabled, processing input string")
	} else {
		// If metricType is false, just pass through the input
		outputString = fmt.Sprintf("Passed through (MetricType disabled): %s", inputString)
		a.logger.Infof("MetricType is disabled, passing through input string")
	}

	// Set output using the proper Flogo API
	ctx.SetOutput("outputString", outputString)

	a.logger.Debugf("Try Activity completed with output: %s", outputString)

	return true, nil
}
