package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	app.Mailer = app.CreateMailer()
	go app.listenForMail()

	go app.listenForShutdown()
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

func (c *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.shutdown()
	os.Exit(0)
}

func (c *Config) shutdown() {
	c.InfoLog.Print("Shutdown app")
	c.Wait.Wait()
	c.Mailer.DoneChan <- true
	close(c.Mailer.MailerChan)
	close(c.Mailer.ErrorChan)
	close(c.Mailer.DoneChan)
	c.InfoLog.Print("closing channels and shutdown app")
}

func (c *Config) CreateMailer() Mail {
	errChan := make(chan error)
	mailerChan := make(chan Message, 100)
	doneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Username:    "",
		Password:    "",
		Encryption:  "none",
		FromAddress: "info@mycompany.com",
		FromName:    "Info",
		Wait:        c.Wait,
		MailerChan:  mailerChan,
		ErrorChan:   errChan,
		DoneChan:    doneChan,
	}

	return m
}
