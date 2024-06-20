package model

import "fmt"

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Message struct {
	Body string `json:"body"`
}

type Envelope struct {
	From    Contact   `json:"from"`
	To      []Contact `json:"to"`
	Content *Message  `json:"content"`
	Subject string    `json:"subject"`
}

type Mailbox struct {
	Name     string     `json:"name"`
	Messages []Envelope `json:"messages"`
}

func (contact *Contact) String() string {
	if contact.Name != "" {
		return fmt.Sprintf("%s <%s>", contact.Name, contact.Email)
	}
	return contact.Email
}
