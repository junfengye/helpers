package helpers

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelError = "error"

	LogFormatJSON = "json"
	LogFormatText = "text"
)

var Logger = &logger{}

type logger struct {
}

var zapLogger *zap.Logger

func init() {
	_ = InitLogger(LogLevelInfo, LogFormatJSON, false)
}

func InitLogger(loglevel, logFormat string, isDev bool, configCb ...func(*zap.Config)) error {
	var config zap.Config
	if isDev {
		config = zap.NewDevelopmentConfig()
		if logFormat == LogFormatText {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // for setting color
		}
	} else {
		config = zap.NewProductionConfig()
	}
	config.Encoding = getLogFormat(logFormat)
	config.Level.SetLevel(getLogLevel(loglevel))
	if len(configCb) > 0 {
		configCb[0](&config)
	}
	var err error
	zapLogger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(zapLogger)
	return nil
}

func getLogFormat(logFormat string) string {
	switch logFormat {
	case LogFormatJSON:
		return LogFormatJSON
	case LogFormatText:
		return "console"
	default:
		Logger.LogInfo(internalRequestID, "Invalid log format provided switching to log format json", nil)
		return LogFormatJSON
	}
}

func getLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case LogLevelDebug:
		return zap.DebugLevel
	case LogLevelInfo:
		return zap.InfoLevel
	case LogLevelError:
		return zap.ErrorLevel
	default:
		Logger.LogInfo(internalRequestID, "Invalid log level provided switching to log level info", nil)
		return zap.InfoLevel
	}
}

// Error stores error info
type Error struct {
	message  string
	rawError string
}

// NewError creates a new error object
func NewError(message, rawError string) Error {
	return Error{message: message, rawError: rawError}
}

func (e Error) RawError() string {
	return e.rawError
}

func (e Error) Error() string {
	return e.message
}

// LogErrorEnabled check LogError is enabled
func (l *logger) LogErrorEnabled() bool {
	return zapLogger.Core().Enabled(zapcore.ErrorLevel)
}

// LogError logs the error in the proper format
func (l *logger) LogError(requestID, message string, err error, fields map[string]interface{}) error {
	value, ok := err.(Error)
	if ok {
		return value
	}
	if err == nil {
		err = errors.New("")
	}
	// Log the error
	if fields != nil {
		zapLogger.Error(message, zap.Any("error", err.Error()), zap.String("requestId", requestID), zap.Any("fields", fields))
	} else {
		zapLogger.Error(message, zap.Any("error", err.Error()), zap.String("requestId", requestID))
	}
	return Error{message: message, rawError: err.Error()}
}

// LogWarnEnabled check LogWarn is enabled
func (l *logger) LogWarnEnabled() bool {
	return zapLogger.Core().Enabled(zapcore.WarnLevel)
}

// LogWarn logs the warning message in the proper format
func (l *logger) LogWarn(requestID, message string, fields map[string]interface{}) {
	if fields != nil {
		zapLogger.Warn(message, zap.String("requestId", requestID), zap.Any("fields", fields))
	} else {
		zapLogger.Warn(message, zap.String("requestId", requestID))
	}
}

// LogInfoEnabled check LogInfo is enabled
func (l *logger) LogInfoEnabled() bool {
	return zapLogger.Core().Enabled(zapcore.InfoLevel)
}

// LogInfo logs the info message in the proper format
func (l *logger) LogInfo(requestID, message string, fields map[string]interface{}) {
	if fields != nil {
		zapLogger.Info(message, zap.String("requestId", requestID), zap.Any("fields", fields))
	} else {
		zapLogger.Info(message, zap.String("requestId", requestID))
	}
}

// LogDebugEnabled check LogDebug is enabled
func (l *logger) LogDebugEnabled() bool {
	return zapLogger.Core().Enabled(zapcore.DebugLevel)
}

// LogDebug logs the debug message in proper format
func (l *logger) LogDebug(requestID, message string, fields map[string]interface{}) {
	if fields != nil {
		zapLogger.Debug(message, zap.String("requestId", requestID), zap.Any("fields", fields))
	} else {
		zapLogger.Debug(message, zap.String("requestId", requestID))
	}
}

// LogFatalEnabled check LogFatal is enabled
func (l *logger) LogFatalEnabled() bool {
	return zapLogger.Core().Enabled(zapcore.FatalLevel)
}

// LogFatal logs the fatal message in proper format
func (l *logger) LogFatal(requestID, message string, fields map[string]interface{}) {
	if fields != nil {
		zapLogger.Fatal(message, zap.String("requestId", requestID), zap.Any("fields", fields))
	} else {
		zapLogger.Fatal(message, zap.String("requestId", requestID))
	}
}
