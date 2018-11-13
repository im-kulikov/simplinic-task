package misc

import "time"

var (
	BuildTime    = time.Now().Format(time.RFC3339Nano)
	BuildVersion = "dev"
)
