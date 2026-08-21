package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"example.com/nursery-cms/service"
	"example.com/nursery-cms/store"
	"example.com/nursery-cms/transport"
)

func main() {
	path := flag.String("db", "nursery-cms.db", "bbolt database path")
	listen := flag.String("listen", "", "HTTP listen address")
	flag.Parse()
	db, err := store.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	svc := service.New(db)
	if *listen != "" {
		server := &http.Server{Addr: *listen, Handler: transport.NewHTTP(svc).Handler()}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := transport.NewCLI(svc).Run(context.Background(), flag.Args(), os.Stdout); err != nil {
		if len(flag.Args()) == 0 {
			fmt.Fprintln(os.Stdout, transport.Usage())
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
