package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/redirects"
)

const minFileLifeSpan = time.Second
const defaultFileLifeSpan = time.Hour * 24

var log = logger.GetOrCreate("fileLogger")

// FileLogger defines the behaviour of a file logger
type FileLogger interface {
	ChangeFileLifeSpan(newDuration time.Duration) error
	IsInterfaceNil() bool
	Close() error
}

type fileLogger struct {
	chLifeSpanChanged chan time.Duration
	mutFile           sync.Mutex
	currentFile       *os.File
	workingDir        string
	defaultLogsPath   string
	logFilePrefix     string
	cancelFunc        func()
	mutIsClosed       sync.Mutex
	isClosed          bool
}

// NewFileLogging created a new file logger instance
func NewFileLogging(workingDir string, defaultLogsPath string, logFilePrefix string) (*fileLogger, error) {
	fl := &fileLogger{
		workingDir:        workingDir,
		defaultLogsPath:   defaultLogsPath,
		logFilePrefix:     logFilePrefix,
		chLifeSpanChanged: make(chan time.Duration),
		isClosed:          false,
	}
	fl.recreateLogFile()

	runtime.SetFinalizer(fl, func(fileLogHandler *fileLogger) {
		_ = fileLogHandler.currentFile.Close()
	})

	ctx, cancelFunc := context.WithCancel(context.Background())
	go fl.autoRecreateFile(ctx)
	fl.cancelFunc = cancelFunc

	return fl, nil
}

// ChangeFileLifeSpan will set file life span
func (fl *fileLogger) ChangeFileLifeSpan(newDuration time.Duration) error {
	if newDuration < minFileLifeSpan {
		return fmt.Errorf("%w, provided %v", core.ErrInvalidLogFileMinLifeSpan, newDuration)
	}

	fl.mutIsClosed.Lock()
	defer fl.mutIsClosed.Unlock()

	if fl.isClosed {
		return core.ErrFileLoggingProcessIsClosed
	}

	fl.chLifeSpanChanged <- newDuration
	return nil
}

// Close will close the logging file
func (fl *fileLogger) Close() error {
	fl.mutIsClosed.Lock()
	if fl.isClosed {
		fl.mutIsClosed.Unlock()
		return nil
	}

	fl.isClosed = true
	fl.mutIsClosed.Unlock()

	fl.mutFile.Lock()
	err := fl.currentFile.Close()
	fl.mutFile.Unlock()

	fl.cancelFunc()

	return err
}

// IsInterfaceNil returns true if there is no value under the interface
func (fl *fileLogger) IsInterfaceNil() bool {
	return fl == nil
}

func (fl *fileLogger) createFile() (*os.File, error) {
	logDirectory := filepath.Join(fl.workingDir, fl.defaultLogsPath)

	return core.CreateFile(
		core.ArgCreateFileArgument{
			Prefix:        fl.logFilePrefix,
			Directory:     logDirectory,
			FileExtension: "log",
		},
	)
}

func (fl *fileLogger) recreateLogFile() {
	newFile, err := fl.createFile()
	if err != nil {
		log.Error("error creating new log file", "error", err)
		return
	}

	fl.mutFile.Lock()
	defer fl.mutFile.Unlock()

	oldFile := fl.currentFile
	err = logger.AddLogObserver(newFile, &logger.PlainFormatter{})
	if err != nil {
		log.Error("error adding log observer", "error", err)
		return
	}

	errNotCritical := redirects.RedirectStderr(newFile)
	log.LogIfError(errNotCritical, "step", "redirecting std error")

	fl.currentFile = newFile

	if oldFile == nil {
		return
	}

	errNotCritical = oldFile.Close()
	log.LogIfError(errNotCritical, "step", "closing old log file")

	errNotCritical = logger.RemoveLogObserver(oldFile)
	log.LogIfError(errNotCritical, "step", "removing old log observer")
}

func (fl *fileLogger) autoRecreateFile(ctx context.Context) {
	fileLifeSpan := defaultFileLifeSpan
	for {
		select {
		case <-ctx.Done():
			log.Debug("closing fileLogger.autoRecreateFile go routine")
			return
		case <-time.After(fileLifeSpan):
			fl.recreateLogFile()
		case fileLifeSpan = <-fl.chLifeSpanChanged:
			log.Debug("changed log file span", "new value", fileLifeSpan)
		}
	}
}
