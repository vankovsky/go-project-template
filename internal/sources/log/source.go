package log

import (
	"fmt"
	"log"
	"os"

	"go-project-template/internal/pkg/configs"

	"github.com/getsentry/sentry-go"
)

var useSentry bool

func Init(cfg *configs.Root) error {
	useSentry = cfg.Sentry.DSN != ""

	if useSentry {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         cfg.Sentry.DSN,
			Environment: cfg.App.Environment,
			Release:     configs.Version,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
}

type logNameSpace struct {
	i, e      *log.Logger
	hub       *sentry.Hub
	nameSpace string
}

func New(nameSpace string) Logger {
	l := logNameSpace{
		nameSpace: nameSpace,
		i:         log.New(os.Stdout, "I ", log.LstdFlags|log.Lmicroseconds),
		e:         log.New(os.Stdout, "E ", log.LstdFlags|log.Lmicroseconds),
	}

	if useSentry {
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("space", nameSpace)
		})
		l.hub = localHub
	}
	return &l
}

func (l *logNameSpace) Info(format string, args ...interface{}) {
	l.i.Printf("%v: %v", l.nameSpace, fmt.Sprintf(format, args...))
}
func (l *logNameSpace) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.e.Printf("%v: %v", l.nameSpace, msg)
	if useSentry {
		l.hub.CaptureMessage(msg)
	}
}
