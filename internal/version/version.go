package version

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

//go:embed VERSION
var versionString string

type showVersion bool

func (sv *showVersion) IsBoolFlag() bool { return true }
func (sv *showVersion) String() string   { return fmt.Sprintf("%v", *sv) }
func (sv *showVersion) Set(v string) error {
	isTrue, err := strconv.ParseBool(v)
	if isTrue {
		fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), versionString)
		os.Exit(0)
	}
	return err
}

// AddVersionFlag installs a -version command-line flag in the global flag set.
func AddVersionFlag() {
	var sv showVersion
	flag.CommandLine.Var(&sv, "version", "show build information and exit")
}
