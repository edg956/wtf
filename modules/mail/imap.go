package mail

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type imapSettings struct {
	address  string
	username string
	password string
}

type IMAPClient struct {
	client goImapClient
}

type goImapClient interface {
	Login(username, password string) error
	List(ref, name string, mailboxes chan *imap.MailboxInfo) error
	Status(name string, items []imap.StatusItem) (*imap.MailboxStatus, error)
}

func newImapClient(settings *imapSettings) (*IMAPClient, error) {
	c, err := client.DialTLS(settings.address, nil)

	if err != nil {
		return nil, err
	}

	err = c.Login(settings.username, settings.password)

	if err != nil {
		return nil, err
	}

	return &IMAPClient{
		client: c,
	}, nil
}

func (client *IMAPClient) listMailboxes(config *pageConfig) ([]*mailbox, error) {
	mailboxes := make(chan *imap.MailboxInfo, config.pageSize)
	done := make(chan error, 1)
	defer close(done)

	go func() {
		done <- client.client.List("", "*", mailboxes)
	}()

	if err := <-done; err != nil {
		return nil, err
	}

	mailboxesArray := make([]*mailbox, 0, len(mailboxes))

	for mbox := range mailboxes {
		mailboxStatus, err := client.client.Status(mbox.Name, []imap.StatusItem{imap.StatusMessages, imap.StatusUnseen})
		if err != nil {
			mailboxStatus = &imap.MailboxStatus{Unseen: 0, Messages: 0}
		}
		mailboxesArray = append(mailboxesArray, &mailbox{
			name:     mbox.Name,
			unread:   int(mailboxStatus.Unseen),
			messages: int(mailboxStatus.Messages),
		})
	}

	return mailboxesArray, nil
}
