package flag

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	showVersion  bool
	tomlFilepath string
	logger       *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
}

func ArgParse(version, revision string) (*Config, error) {
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.BoolVar(&showVersion, "v", false, "show application version")
	fl.StringVar(&tomlFilepath, "f", "", "TOML configuration filepath")
	fl.Parse(os.Args[1:])

	if showVersion {
		fmt.Println(fmt.Sprintf(`version: %s (revision %s)`, version, revision))
		os.Exit(1)
	}
	if tomlFilepath == "" {
		fmt.Println(`option "-f" is requirement`)
		os.Exit(2)
	}

	return ParseConfig(tomlFilepath)
}
