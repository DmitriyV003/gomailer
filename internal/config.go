package mailer

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const webPort = "80"

type Config struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Mailer   Mail
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

func (c *Config) Serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: c.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (c *Config) Shutdown() {
	c.InfoLog.Print("Shutdown app")
	c.Wait.Wait()
	c.Mailer.DoneChan <- true
	close(c.Mailer.MailerChan)
	close(c.Mailer.ErrorChan)
	close(c.Mailer.DoneChan)
	c.InfoLog.Print("closing channels and shutdown app")
}
