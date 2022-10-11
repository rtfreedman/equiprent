package util

import (
	"equiprent/internal/util/config"
	"equiprent/internal/util/flags"
	"equiprent/internal/util/log"
)

// Initialize utilities in appropriate order
func Initialize() {
	flags.Initialize()
	config.Initialize()
	log.Initialize()
}

// Stop all utilities
func Stop() {
	log.Stop()
}
