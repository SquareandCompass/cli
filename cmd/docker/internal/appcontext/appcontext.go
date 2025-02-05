package appcontext

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	appContextCache context.Context
	appContextOnce  sync.Once
)

// Context returns a static context that reacts to termination signals of the
// running process. Useful in CLI tools.
func Context() context.Context {
	appContextOnce.Do(func() {
		signals := make(chan os.Signal, 2048)
		signal.Notify(signals, terminationSignals...)

		const exitLimit = 3
		retries := 0

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		appContextCache = ctx

		go func() {
			for {
				<-signals
				cancel()
				retries++
				if retries >= exitLimit {
					logrus.Errorf("got %d SIGTERM/SIGINTs, forcing shutdown", retries)
					os.Exit(1)
				}
			}
		}()
	})
	return appContextCache
}
