package main

import (
	"github.com/jafarsirojov/bank-auth/cmd/auth/app"
	"github.com/jafarsirojov/bank-auth/pkg/core/token"
	"github.com/jafarsirojov/bank-auth/pkg/core/users"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jafarsirojov/mux/pkg/mux"
	"log"
	"net"
	"net/http"
)

func main() {
	flag.Parse()
	envPort, ok := FromFlagOrEnv(*port, ENV_PORT)
	if !ok {
		log.Println("can't port")
		return
	}
	envDsn, ok := FromFlagOrEnv(*dsn, ENV_DSN)
	if !ok {
		log.Println("can't dsn")
		return
	}
	envHost, ok := FromFlagOrEnv(*host, ENV_HOST)
	if !ok {
		log.Println("can't host")
		return
	}
	addr := net.JoinHostPort(envHost, envPort)
	log.Println("starting server!")
	log.Printf("host = %s, port = %s\n", envHost, envPort)
	pool, err := pgxpool.Connect(
		context.Background(),
		envDsn,
	)
	if err != nil {
		panic(err)
	}

	userSvc := users.NewService(pool)
	userSvc.Start()
	tokenSvc := token.NewService(secret,pool)
	exactMux := mux.NewExactMux()
	server := app.NewServer(exactMux, pool, secret, tokenSvc, userSvc)
	server.Start()
	panic(http.ListenAndServe(addr, server))
}
