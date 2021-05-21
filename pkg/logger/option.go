package logger

import "github.com/rs/zerolog"

// An Option configures a Logger
type Option interface {
	Apply(zerolog.Logger) zerolog.Logger
}

// OptionFunc is a function that configures a zerolog.Logger
type OptionFunc func(zerolog.Logger) zerolog.Logger

// Apply is a function that set value to zerolog.Logger
func (f OptionFunc) Apply(engine zerolog.Logger) zerolog.Logger {
	return f(engine)
}

func WithFields(fields map[string]interface{}) Option {
	return OptionFunc(func(engine zerolog.Logger) zerolog.Logger {
		return engine.With().Fields(fields).Logger()
	})
}

func WithStr(key, value string) Option {
	return OptionFunc(func(engine zerolog.Logger) zerolog.Logger {
		return engine.With().Str(key, value).Logger()
	})
}
