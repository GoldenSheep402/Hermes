package healthcheck

import (
	"github.com/juanjiTech/jin"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
	"net/http"
	"strings"
	"sync"
)

var (
	defaultHealthChecker   = map[string]Checker{}
	defaultHealthCheckerMu sync.Mutex
)

func RegisterHealthChecker(checkers ...Checker) error {
	for _, checker := range checkers {
		err := func() error {
			defaultHealthCheckerMu.Lock()
			defer defaultHealthCheckerMu.Unlock()
			if _, ok := defaultHealthChecker[checker.Name()]; ok {
				return ErrConflictCheckerName
			}
			defaultHealthChecker[checker.Name()] = checker
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	ErrConflictCheckerName = errors.New("health checker names conflict")
	ErrCheckFailed         = errors.New("health check failed")
)

var sg singleflight.Group

func Register(e *jin.Engine) {
	e.GET("/healthz", NewHandler(lo.MapToSlice(defaultHealthChecker, func(name string, checker Checker) Checker {
		return checker
	})...))
}

func NewHandler(checkers ...Checker) jin.HandlerFunc {
	handlerID := ulid.Make().String()

	return func(c *jin.Context) {
		verbose := c.Request.URL.Query().Has("verbose")
		excludes := strings.Split(c.Request.URL.Query().Get("exclude"), ",")

		body, err, _ := sg.Do(handlerID, func() (interface{}, error) {
			status := NewStatus(len(checkers))
			var eg errgroup.Group

			for _, checker := range checkers {
				eg.Go(func() error {
					name := checker.Name()
					if len(excludes) > 0 {
						for _, excludedName := range excludes {
							if excludedName == name {
								status.Set(name, StatusExcluded)
								return nil
							}
						}
					}

					if _, ok := status.Get(name); ok {
						return ErrConflictCheckerName
					}

					if checker.Pass() {
						status.Set(name, StatusOK)
					} else {
						status.Set(name, StatusError)
						return ErrCheckFailed
					}

					return nil
				})
			}
			err := eg.Wait()
			return status.String(verbose), err
		})
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
		_, _ = c.Writer.WriteString(body.(string))
	}
}
