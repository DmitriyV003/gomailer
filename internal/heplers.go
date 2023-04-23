package mailer

func (c *Config) sendMail(msg Message) {
	c.Wait.Add(1)
	c.Mailer.MailerChan <- msg
}
