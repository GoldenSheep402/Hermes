package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/juanjiTech/jframe/conf"
)

func Init() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: conf.Get().SentryDsn,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
}
