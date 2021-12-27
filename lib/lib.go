package lib

type Dispatcher func(settings SettingsType, stdout, stderr []byte, args ...string) (any, any, error)

type SettingsType struct {
	ExitCode int
}
