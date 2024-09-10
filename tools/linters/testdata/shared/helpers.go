package shared

type (
	logger interface {
		Debug(string, ...any)
		Error(string, ...any)
		Info(string, ...any)
		Printf(string, ...any)
		Sync() error
		Trace(string, ...any)
		Warn(string, ...any)
		Verbose() bool
	}

	log struct{}
)

func (l *log) Debug(msg string, fields ...any)  {}
func (l *log) Error(msg string, fields ...any)  {}
func (l *log) Info(msg string, fields ...any)   {}
func (l *log) Printf(msg string, fields ...any) {}
func (l *log) Sync() error                      { return nil }
func (l *log) Trace(msg string, fields ...any)  {}
func (l *log) Warn(msg string, fields ...any)   {}
func (l *log) Verbose() bool                    { return false }

func Logger() logger {
	return &log{}
}
