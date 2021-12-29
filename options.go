package remauth

import (
	"time"
)

type Options struct {
	Debug     bool
	CheckUrl  string
	Timeout   int
	UseGzip   bool
	CacheTime time.Duration
}
