package main

import (
	"log"
	"sync"
)

type Config struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Mailer   Mail
}
