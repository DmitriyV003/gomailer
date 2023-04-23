package mailer

func (c *Config) SendMail(msg Message) {
	c.Wait.Add(1)
	c.Mailer.MailerChan <- msg
}
