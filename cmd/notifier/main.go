package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/common/logging"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/notifier"
	"github.com/urfave/cli"
)

const (
	defaultLogsPath    = "logs"
	logFilePrefix      = "event-notifier"
	logFileLifeSpanSec = 86400
)

var (
	cliHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
	log = logger.GetOrCreate("eventNotifier")

	logLevel = cli.StringFlag{
		Name:  "log-level",
		Usage: "This flag specifies the log level. Options: *:NONE | ERROR | WARN | INFO | DEBUG | TRACE",
		Value: fmt.Sprintf("*:%s", logger.LogInfo.String()),
	}

	logSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving",
	}

	generalConfigFile = cli.StringFlag{
		Name:  "general-config",
		Usage: "The path for the general config",
		Value: "./config/config.toml",
	}

	workingDirectory = cli.StringFlag{
		Name:  "working-directory",
		Usage: "This flag specifies the directory where the eventNotifier proxy will store logs.",
		Value: "",
	}

	apiType = cli.StringFlag{
		Name:  "api-type",
		Usage: "This flag specifies the api type, it defines the way in which it will expose the events. Options: rabbit-api | notifier",
		Value: "notifier",
	}
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = cliHelpTemplate
	app.Name = "Elrond event "
	app.Flags = []cli.Flag{
		logLevel,
		logSaveFile,
		generalConfigFile,
		workingDirectory,
		apiType,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}
	app.Action = startEventNotifierProxy

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startEventNotifierProxy(ctx *cli.Context) error {
	log.Info("starting eventNotifier proxy...")

	flagsConfig, err := getFlagsConfig(ctx)
	if err != nil {
		return err
	}

	fileLogging, err := initLogger(flagsConfig)
	if err != nil {
		return err
	}

	cfg, err := config.LoadConfig(flagsConfig.GeneralConfigPath)
	if err != nil {
		return err
	}
	cfg.Flags = flagsConfig

	notifierRunner, err := notifier.NewNotifierRunner(cfg)
	if err != nil {
		return err
	}

	err = notifierRunner.Start()
	if err != nil {
		return err
	}

	if !check.IfNil(fileLogging) {
		err = fileLogging.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func getFlagsConfig(ctx *cli.Context) (*config.FlagsConfig, error) {
	flagsConfig := &config.FlagsConfig{}

	workingDir, err := getWorkingDir(ctx)
	if err != nil {
		return nil, err
	}
	flagsConfig.WorkingDir = workingDir

	flagsConfig.LogLevel = ctx.GlobalString(logLevel.Name)
	flagsConfig.SaveLogFile = ctx.GlobalBool(logSaveFile.Name)
	flagsConfig.GeneralConfigPath = ctx.GlobalString(generalConfigFile.Name)
	flagsConfig.APIType = ctx.GlobalString(apiType.Name)

	return flagsConfig, nil
}

func initLogger(config *config.FlagsConfig) (logging.FileLogger, error) {
	err := logger.SetLogLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	var fileLogging logging.FileLogger
	if config.SaveLogFile {
		fileLogging, err = logging.NewFileLogging(config.WorkingDir, defaultLogsPath, logFilePrefix)
		if err != nil {
			return fileLogging, err
		}
	}
	if !check.IfNil(fileLogging) {
		err = fileLogging.ChangeFileLifeSpan(time.Second * time.Duration(logFileLifeSpanSec))
		if err != nil {
			return nil, err
		}
	}

	return fileLogging, nil
}

func getWorkingDir(ctx *cli.Context) (string, error) {
	if ctx.IsSet(workingDirectory.Name) {
		return ctx.GlobalString(workingDirectory.Name), nil
	}

	return os.Getwd()
}
