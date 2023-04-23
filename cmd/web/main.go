package main

import (
	mailer "github.com/dmitriyv003/gomailer/internal"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	wg := sync.WaitGroup{}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := mailer.Config{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
	}

	app.Mailer = app.CreateMailer()
	go app.ListenForMail()

	go listenForShutdown(&app)
	app.Serve()
}

func listenForShutdown(c *mailer.Config) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Shutdown()
	os.Exit(0)
}
