package enum

type LogLevelType uint8

const (
	LogInfoLevel  LogLevelType = 1
	LogWarnLevel  LogLevelType = 2
	LogErrorLevel LogLevelType = 3
)

func (logLevelType LogLevelType) String() string {
	switch logLevelType {
	case LogInfoLevel:
		return "Info"
	case LogWarnLevel:
		return "Warn"
	case LogErrorLevel:
		return "Error"
	default:
		return "Info"
	}
}
