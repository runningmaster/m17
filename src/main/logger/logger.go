package logger

// Logger is common interface for logging.
// See https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g for details.
type Logger interface {
	Printf(string, ...interface{})
}
