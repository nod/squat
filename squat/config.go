package squat

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

type RuntimeConfig struct {
	QueueUrl   string `envcfg:"URL"`
	Region     string `envcfg:"REGION"`
	ShowConfig bool
	Verbose    bool
}

// override with environment variables if present
func parseEnv(cfg *RuntimeConfig) {
	ctype := reflect.TypeOf(*cfg)
	for i := 0; i < ctype.NumField(); i++ {
		fld := ctype.Field(i)
		envkey := fld.Tag.Get("envcfg")
		if envkey == "" {
			continue
		}
		envnm := fmt.Sprintf("SQUAT_%s", envkey)
		val := os.Getenv(envnm)
		if val == "" {
			continue
		}
		cf := reflect.ValueOf(cfg).Elem().Field(i)
		cf.SetString(val)
	}
}

func BuildRuntimeConfig() *RuntimeConfig {
	const (
		defaultPtn  = "squat"
		defaultNone = ""
		regHelp     = "aws region. env: SQUAT_REGION"
		showHelp    = "show parsed config, vars and env, and exit"
		verboseHelp = "be a bit louder about what is happening"
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <sqs-queue-url>\n\n",
			os.Args[0],
		)
		fmt.Fprint(os.Stderr,
			"  Options can be overriden with environment variables\n\n",
		)
		flag.PrintDefaults()
	}
	cfg := &RuntimeConfig{}
	flag.StringVar(&cfg.Region, "region", defaultNone, regHelp)
	flag.BoolVar(&cfg.ShowConfig, "config", false, showHelp)
	flag.BoolVar(&cfg.Verbose, "verbose", false, verboseHelp)
	flag.Parse()
	// override from env
	parseEnv(cfg)
	cfg.QueueUrl = flag.Arg(0)
	// should we show config and exit?
	if cfg.ShowConfig {
		fmt.Fprintf(os.Stderr, "CONFIG\n%+v\n", cfg)
		os.Exit(1)
	}
	// if we don't have a streamname, error
	if cfg.QueueUrl == "" {
		fmt.Fprint(os.Stderr, "ERR missing sqs-queue-url\n\n")
		flag.Usage()
		os.Exit(1)
	}
	return cfg
}
