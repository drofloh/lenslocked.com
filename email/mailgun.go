package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to lenslocked.com!"
	resetSubject   = "instructions for resetting your password."
	resetBaseURL   = "https://www.lenslocked.com/reset"
)

const welcomeText = `Hi!

Welcome to lenslocked.com!`

const welcomeHTML = `Hi!<br/>
<br/>
Welcome to lenslocked.com!<br/>`

const resetTextTmpl = `Hi!

It appears you have requested a password reset, follow the link below:

%s

If asked for a token, please use the following:

%s

If you didnt request this, please ignore.

Regards, Support.`

const resetHTMLTmpl = `Hi!<br/>
<br/>
It appears you have requested a password reset, follow the link below:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If asked for a token, please use the following:<br/>
<br/>
%s<br/>
<br/>
If you didnt request this, please ignore.<br/>
<br/>
Regards, Support.<br/>`

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		from: "support@lenslocked.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
