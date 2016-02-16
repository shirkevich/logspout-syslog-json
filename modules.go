package main

import (
	_ "github.com/gliderlabs/logspout/adapters/raw"
	_ "github.com/shirkevich/logspout-syslog-json/adapters/json"
	_ "github.com/shirkevich/logspout-syslog-json/adapters/syslog"
	_ "github.com/gliderlabs/logspout/httpstream"
	_ "github.com/gliderlabs/logspout/routesapi"
	_ "github.com/gliderlabs/logspout/transports/tcp"
	_ "github.com/gliderlabs/logspout/transports/udp"
)
