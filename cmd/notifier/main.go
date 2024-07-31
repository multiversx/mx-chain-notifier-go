package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/common/logging"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/notifier"
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

	apiConfigFile = cli.StringFlag{
		Name:  "api-config",
		Usage: "The path for the api config",
		Value: "./config/api.toml",
	}

	workingDirectory = cli.StringFlag{
		Name:  "working-directory",
		Usage: "This flag specifies the directory where the eventNotifier proxy will store logs.",
		Value: "",
	}

	// deprecated - keep it for backwards compatibility
	// TODO: remove flag
	apiType = cli.StringFlag{
		Name:  "api-type",
		Usage: "Deprecated: This flag specifies the api type, it defines the way in which it will expose the events. Options: rabbit-api | notifier",
		Value: "notifier",
	}

	publisherType = cli.StringFlag{
		Name:  "publisher-type",
		Usage: "This flag specifies the publisher type, it defines the way in which it will expose the events. Options: " + common.MessageQueuePublisherType + " | " + common.WSPublisherType,
		Value: common.MessageQueuePublisherType,
	}

	checkDuplicates = cli.BoolTFlag{
		Name:  "check-duplicates",
		Usage: "Boolean option to check the duplicates. Set this to '--check-duplicates=false' to disable the check",
	}
)

// appVersion should be populated at build time using ldflags
// Usage examples:
// linux/mac:
//
//	go build -i -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)"
//
// windows:
//
//	for /f %i in ('git describe --tags --long --dirty') do set VERS=%i
//	go build -i -v -ldflags="-X main.appVersion=%VERS%"
var appVersion = "undefined"

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = cliHelpTemplate
	app.Name = "MultiversX events notifier"
	app.Usage = "An events notifier service to handle blockchain events to subscribers"
	app.Flags = []cli.Flag{
		logLevel,
		logSaveFile,
		generalConfigFile,
		apiConfigFile,
		workingDirectory,
		apiType,
		publisherType,
		checkDuplicates,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The MultiversX Team",
			Email: "contact@multiversx.com",
		},
	}

	machineID := core.GetAnonymizedMachineID(app.Name)
	app.Version = fmt.Sprintf("%s/%s/%s-%s/%s", appVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH, machineID)
	app.Action = startEventNotifierProxy

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startEventNotifierProxy(ctx *cli.Context) error {
	log.Info("starting eventNotifier proxy...")

	cfgs, err := readConfigs(ctx)
	if err != nil {
		return err
	}

	fileLogging, err := initLogger(&cfgs.Flags)
	if err != nil {
		return err
	}

	notifierRunner, err := notifier.NewNotifierRunner(cfgs)
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

func readConfigs(ctx *cli.Context) (*config.Configs, error) {
	flagsConfig, err := getFlagsConfig(ctx)
	if err != nil {
		return nil, err
	}

	mainConfig, err := config.LoadMainConfig(flagsConfig.GeneralConfigPath)
	if err != nil {
		return nil, err
	}

	if ctx.IsSet(checkDuplicates.Name) {
		mainConfig.General.CheckDuplicates = ctx.GlobalBool(checkDuplicates.Name)
	}

	apiConfig, err := config.LoadAPIConfig(flagsConfig.APIConfigPath)
	if err != nil {
		return nil, err
	}

	return &config.Configs{
		MainConfig:      *mainConfig,
		ApiRoutesConfig: *apiConfig,
		Flags:           *flagsConfig,
	}, nil
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
	flagsConfig.APIConfigPath = ctx.GlobalString(apiConfigFile.Name)
	flagsConfig.PublisherType = ctx.GlobalString(publisherType.Name)

	// TODO: remove deprecated flag
	if ctx.IsSet(apiType.Name) {
		// override for backwards compatibility
		flagsConfig.PublisherType, err = handleAPIType(ctx)
		if err != nil {
			return nil, err
		}
	}

	return flagsConfig, nil
}

// TODO: remove deprecated flag
func handleAPIType(ctx *cli.Context) (string, error) {
	apiType := ctx.GlobalString(apiType.Name)
	switch apiType {
	case "rabbit-api":
		return common.MessageQueuePublisherType, nil
	case "notifier":
		return common.WSPublisherType, nil
	default:
		return "", common.ErrInvalidAPIType
	}
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
