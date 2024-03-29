// Package database is a simple database with a cleaner when elements are too old.
package database

import (
	"time"
)

const (
	defaultLifetime   time.Duration = time.Minute * 5
	defaultCleanCycle time.Duration = defaultLifetime / 3
)
