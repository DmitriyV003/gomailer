package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

const webPort = "80"

func main() {
	wg := sync.WaitGroup{}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := Config{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
	}
	app.serve()
}

func (c *Config) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: c.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
