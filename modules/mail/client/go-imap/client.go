package go_imap

import (
	"github.com/emersion/go-imap"
	"github.com/wtfutil/wtf/modules/mail/model"
	_ "github.com/wtfutil/wtf/modules/mail/model"
)

type IMAPClient interface {
	List(ref, name string, mailboxes chan *imap.MailboxInfo) error
	Select(name string, readOnly bool) (*imap.MailboxStatus, error)
	Fetch(set *imap.SeqSet, items []imap.FetchItem, messages chan *imap.Message) error
	Login(username, password string) error
}

type Config struct {
	Page     uint32
	PageSize uint32
}

type Client struct {
	Config     Config
	IMAPClient IMAPClient
}

type FetchFunc func(set *imap.SeqSet, items []imap.FetchItem, messages chan *imap.Message) error
type ListFunc func(ref, name string, mailboxes chan *imap.MailboxInfo) error

func getSequenceSet(mailbox *imap.MailboxStatus, config *Config) *imap.SeqSet {
	seqSet := new(imap.SeqSet)

	offset := config.Page * config.PageSize

	if offset > mailbox.Messages {
		return seqSet
	}

	to := mailbox.Messages - offset
	from := uint32(1)

	if to > config.PageSize {
		from = to - config.PageSize + 1
	}

	seqSet.AddRange(from, to)

	return seqSet
}

func listMailboxes(listFunc ListFunc, numMailboxes uint32) ([]*imap.MailboxInfo, error) {
	mailboxes := make(chan *imap.MailboxInfo, numMailboxes)
	done := make(chan error, 1)
	defer close(done)

	go func() {
		done <- listFunc("", "*", mailboxes)
	}()

	if err := <-done; err != nil {
		return nil, err
	}

	mailboxesArray := make([]*imap.MailboxInfo, 0, len(mailboxes))

	for mailbox := range mailboxes {
		mailboxesArray = append(mailboxesArray, mailbox)
	}

	return mailboxesArray, nil
}

func listMessages(fetchFunc FetchFunc, mailbox *imap.MailboxStatus, config *Config) ([]*imap.Message, error) {
	seqSet := getSequenceSet(mailbox, config)

	messages := make(chan *imap.Message, config.PageSize)
	done := make(chan error, 1)
	defer close(done)

	go func() {
		done <- fetchFunc(seqSet, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	if err := <-done; err != nil {
		return nil, err
	}

	messageArray := make([]*imap.Message, 0, len(messages))
	for message := range messages {
		messageArray = append(messageArray, message)
	}

	return messageArray, nil
}

func (c *Client) GetMailboxes() ([]model.Mailbox, error) {
	mailboxes, err := listMailboxes(c.IMAPClient.List, c.Config.PageSize)
	if err != nil {
		return nil, err
	}

	result := make([]model.Mailbox, len(mailboxes))

	for i, mailbox := range mailboxes {
		result[i] = model.Mailbox{
			Name: mailbox.Name,
		}
	}

	return result, nil
}

func (c *Client) GetMessages(mailbox model.Mailbox) ([]model.Envelope, error) {
	_, err := c.IMAPClient.Select(mailbox.Name, true)
	if err != nil {
		return nil, err
	}

	messages, err := listMessages(c.IMAPClient.Fetch, &imap.MailboxStatus{Messages: 15}, &c.Config)
	if err != nil {
		return nil, err
	}

	result := make([]model.Envelope, len(messages))

	for i, message := range messages {
		toList := make([]model.Contact, len(message.Envelope.To))
		for j, to := range message.Envelope.To {
			toList[j] = model.Contact{
				Name:  to.PersonalName,
				Email: to.Address(),
			}
		}
		result[i] = model.Envelope{
			Subject: message.Envelope.Subject,
			From: model.Contact{
				Name:  message.Envelope.From[0].PersonalName,
				Email: message.Envelope.From[0].Address(),
			},
			To: toList,
		}
	}

	return result, nil
}
