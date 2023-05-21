// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
)

type level string

// Log levels
const (
	// LogLevelSuccess represents logs at the SUCCESS level
	LogLevelSuccess level = "SUCCESS"

	// LogLevelFail represents logs at the FAIL level
	LogLevelFail level = "FAIL"

	// LogLevelPanic represents logs at the PANIC level
	LogLevelPanic level = "PANIC"

	// LogLevelInfo represents logs at the INFO level
	LogLevelInfo level = "INFO"

	// LogLevelWarn represents logs at the WARN level
	LogLevelWarn level = "WARN"

	// LogLevelDebug represents logs at the DEBUG level
	LogLevelDebug level = "DEBUG"
)

// Color levels
const (
	// logColorSuccess represents the color for logs at the SUCCESS level
	logColorSuccess = "\u001B[42m" // Green

	// logColorFail represents the color for logs at the FAIL level
	logColorFail = "\u001B[43m" // Orange

	// logColorPanic represents the color for logs at the PANIC level
	logColorPanic = "\u001B[41m" // Red

	// logColorInfo represents the color for logs at the INFO level
	logColorInfo = "\u001B[44m" // Blue

	// logColorWarn represents the color for logs at the WARN level
	logColorWarn = "\u001B[45m" // Purple

	// logColorDebug represents the color for logs at the DEBUG level
	logColorDebug = "\u001B[47m" // Gray

	// logColorReset is used to reset the log color
	logColorReset = "\033[0m" // Reset

	// logColorSign is used for the color of the signature label
	logColorSign = "\033[90m" // Light Gray
)

// ANSIColorCode returns the ANSI color code associated with the log level.
func (l level) ANSIColorCode() string {
	var LogColor = map[level]string{
		LogLevelSuccess: logColorSuccess,
		LogLevelFail:    logColorFail,
		LogLevelPanic:   logColorPanic,
		LogLevelInfo:    logColorInfo,
		LogLevelWarn:    logColorWarn,
		LogLevelDebug:   logColorDebug,
	}
	return LogColor[l]
}

// loggedTransport is custom Transport that logs request information.
type loggedTransport struct {
	transport http.RoundTripper
	logger    *log.Logger
}

// RoundTrip implements the RoundTrip method of the http.RoundTripper interface.
func (t *loggedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// When the actual request is made
	startTime := time.Now()
	response, err := t.transport.RoundTrip(req)
	duration := time.Since(startTime)

	// Select log level based on request duration
	logLevel := LogLevelSuccess
	if duration > TimeoutShort {
		logLevel = LogLevelWarn
	}

	// Record request log
	consoleLog(t.logger, logLevel, response.StatusCode, req.Method, req.URL.String(), fmt.Sprintf("Request took %s", duration))

	return response, err
}

// sign returns a signature string for the generated content.
func sign() string {
	return fmt.Sprintf("%s   # generate by %s.%s", logColorSign, Title, logColorReset)
}

// levelText returns the formatted text representation of the log level.
// It applies the corresponding ANSI color code to the level text.
func levelText(l level) string {
	logColorStart := l.ANSIColorCode()
	return fmt.Sprintf("%s[%s]%s", logColorStart, l, logColorReset)
}

// consoleLog is an auxiliary function that outputs log information with
// a level prefix according to the log level and color.
func consoleLog(logger *log.Logger, level level, statusCode int, method, url, message string) {
	logger.Printf("| %20s | %18s | [%d] [%s] %s | %s %s", fileLocation(2), levelText(level), statusCode, method, url, message, sign())
}

// fileLocation returns the file location in the format "filename:line",
// indicating the file name and line number where the function is called from.
// The 'depth' parameter specifies the number of stack frames to skip.
// It uses the runtime.Caller function to retrieve the stack information.
func fileLocation(depth int) string {
	// Retrieve the stack information
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}

	// Get the base file name
	base := filepath.Base(file)

	return fmt.Sprintf("%s:%d", base, line)
}

// ChalkObj writes a log entry with the specified level and object value.
// It uses reflection to extract the value from the object parameter.
// The 'level' parameter represents the log level.
// The 'obj' parameter is the object to be logged.
// It returns the updated Client instance.
func (c *Client[T]) ChalkObj(level level, obj any) *Client[T] {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	c.Config.Logger.Printf("| %20s | %18s | %#v\n", fileLocation(3), levelText(level), v.Interface())
	return c
}

// ChalkStr writes a log entry with the specified level and string value.
// The 'level' parameter represents the log level.
// The 's' parameter is the string to be logged.
// It returns the updated Client instance.
func (c *Client[T]) ChalkStr(level level, s string) *Client[T] {
	c.Config.Logger.Printf("| %20s | %18s | %s\n", fileLocation(3), levelText(level), s)
	return c
}

// ChalkInt writes a log entry with the specified level and integer value.
// The 'level' parameter represents the log level.
// The 'n' parameter is the integer to be logged.
// It returns the updated Client instance.
func (c *Client[T]) ChalkInt(level level, n int) *Client[T] {
	c.Config.Logger.Printf("| %20s | %18s | %d\n", fileLocation(3), levelText(level), n)
	return c
}

// ChalkPrintf writes a formatted log entry with the specified level and arguments.
// The 'level' parameter represents the log level.
// The 'format' parameter is the format string for the log message.
// The 'args' parameter contains the arguments to be formatted.
// It returns the updated Client instance.
func (c *Client[T]) ChalkPrintf(level level, format string, args ...any) *Client[T] {
	message := fmt.Sprintf(format, args...)
	if (level != LogLevelFail && level != LogLevelPanic) || isEmpty(c.Exception.CodeLocation) {
		c.Config.Logger.Printf("| %20s | %18s | %s\n", fileLocation(3), levelText(level), message)
	} else {
		c.Config.Logger.Printf("| %20s | %18s | %s\n", c.Exception.CodeLocation, levelText(level), message)
	}
	return c
}
