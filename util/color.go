package util

import (
	"github.com/fatih/color"
)

// ColorInfo returns a new function that returns info-colorized (green) strings for the
// given arguments with fmt.Sprint().
var ColorInfo = color.New(color.FgGreen).SprintFunc()

// ColorStatus returns a new function that returns status-colorized (blue) strings for the
// given arguments with fmt.Sprint().
var ColorStatus = color.New(color.FgBlue).SprintFunc()

// ColorWarning returns a new function that returns warning-colorized (yellow) strings for the
// given arguments with fmt.Sprint().
var ColorWarning = color.New(color.FgYellow).SprintFunc()

// ColorError returns a new function that returns error-colorized (red) strings for the
// given arguments with fmt.Sprint().
var ColorError = color.New(color.FgRed).SprintFunc()

// ColorBold returns a new function that returns bold-colorized (bold) strings for the
// given arguments with fmt.Sprint().
var ColorBold = color.New(color.Bold).SprintFunc()

// ColorAnswer returns a new function that returns answer-colorized (cyan) strings for the
// given arguments with fmt.Sprint().
var ColorAnswer = color.New(color.FgCyan).SprintFunc()
