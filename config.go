package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/gentlemanautomaton/bindflag"
)

// Config holds a set of asterisk HTTP server configuration values.
type Config struct {
	URLPrefix        string
	DataPath         string
	Logging          bool
	DirectoryListing bool
}

// DefaultConfig holds the default configuration values.
var DefaultConfig = Config{
	DataPath: "data",
}

// ParseEnv will parse environment variables and apply them to the
// configuration.
func (c *Config) ParseEnv() {
	var (
		prefix, hasPrefix     = os.LookupEnv("URL_PREFIX")
		dataPath, hasDataPath = os.LookupEnv("DATA")
		logging, hasLogging   = os.LookupEnv("LOGGING")
		listing, hasListing   = os.LookupEnv("DIRECTORY_LISTING")
	)

	if hasPrefix {
		c.URLPrefix = prefix
	}
	if hasDataPath {
		c.DataPath = dataPath
	}
	if hasLogging {
		if val, err := strconv.ParseBool(logging); err != nil {
			c.Logging = val
		}
	}
	if hasListing {
		if val, err := strconv.ParseBool(listing); err != nil {
			c.DirectoryListing = val
		}
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
	fs.Var(bindflag.Bool(&c.Logging), "logging", "should all requests should be logged?")
	fs.Var(bindflag.Bool(&c.DirectoryListing), "dirlist", "should directory listings be generated?")
}
