package app

import (
	"sync"
)

var version string

var vSetOnce sync.Once

// SetVersion sets the version for the app
func SetVersion(v string) {
	vSetOnce.Do(func() {
		version = v
	})
}

// Version of the app. i.e, the git tag
func Version() string {
	return version
}
