package sentry

import (
	"fmt"
	"github.com/GoldenSheep402/Hermes/conf"
	"github.com/getsentry/sentry-go"
)

func Init() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: conf.Get().SentryDsn,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
}
