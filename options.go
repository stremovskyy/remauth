package remauth

import (
	"time"
)

type Options struct {
	Debug     bool
	CheckUrl  string
	Timeout   int
	CacheTime time.Duration
}
