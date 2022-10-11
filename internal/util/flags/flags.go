package flags

import "github.com/spf13/pflag"

// Flags for the API
var (
	Port   = pflag.Int("port", 8002, "specify port for api")
	Dev    = pflag.Bool("dev", false, "indicate whether or not to run in dev mode")
	Config = pflag.String("config", "", "supply a config file to the api")
)

// Initialize flags
func Initialize() {
	pflag.Parse()
}
