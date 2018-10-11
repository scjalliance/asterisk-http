package main

import (
	"flag"
	"os"

	"github.com/gentlemanautomaton/bindflag"
)

// Config holds a set of asterisk HTTP server configuration values.
type Config struct {
	URLPrefix string
	DataPath  string
}

// DefaultConfig holds the default configuration values.
var DefaultConfig = Config{
	DataPath: "data",
}

// ParseEnv will parse environment variables and apply them to the
// configuration.
func (c *Config) ParseEnv() {
	var (
		prefix, hasPrefix     = os.LookupEnv("URLPREFIX")
		dataPath, hasDataPath = os.LookupEnv("DATA")
	)

	if hasPrefix {
		c.URLPrefix = prefix
	}
	if hasDataPath {
		c.DataPath = dataPath
	}
}

// ParseArgs parses the given argument list and applies them to the
// configuration.
func (c *Config) ParseArgs(args []string, errorHandling flag.ErrorHandling) error {
	fs := flag.NewFlagSet("", errorHandling)
	c.Bind(fs)
	return fs.Parse(args)
}

// Bind will bind the given flag set to the configuration.
func (c *Config) Bind(fs *flag.FlagSet) {
	fs.Var(bindflag.String(&c.URLPrefix), "urlprefix", "URL Prefix")
	fs.Var(bindflag.String(&c.DataPath), "path", "path of asset directory on local filesystem")
}
