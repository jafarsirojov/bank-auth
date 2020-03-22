package main

import (
	"flag"
	"os"
)

var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)

//-host 0.0.0.0 -port 9011 -dsn postgres://user:pass@localhost:5302/app
const ENV_PORT = "PORT"
const ENV_DSN = "DATABASE_URL"
const ENV_HOST = "HOST"

func FromFlagOrEnv(flag string, env string) (value string, ok bool) {
	if flag != "" {
		return flag, true
	}

	return os.LookupEnv(env)
}

var secret = []byte("top secret")
