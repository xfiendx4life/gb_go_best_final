package config

type Config interface {
	// init config from source env or file
	InitConfig(source string) (err error)
	// returns path to csv file
	GetPath() (path string)
	// get commands from config
	// GetOps() (ops Operation[] )
	// read config from file with path
	ReadConfig(path string) (err error)
}