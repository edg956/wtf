package go_imap

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/wtfutil/wtf/modules/mail/model"
	"testing"
)

type getSequenceSetTestCase struct {
	config   *Config
	mailbox  *imap.MailboxStatus
	expected *imap.SeqSet
}

var getSequenceSetTestCases = []getSequenceSetTestCase{
	getSequenceSetTestCase{
		&Config{Page: 0, PageSize: 10},
		&imap.MailboxStatus{Messages: 9},
		&imap.SeqSet{Set: []imap.Seq{imap.Seq{Start: uint32(1), Stop: uint32(9)}}},
	},
	getSequenceSetTestCase{
		&Config{Page: 0, PageSize: 10},
		&imap.MailboxStatus{Messages: 11},
		&imap.SeqSet{Set: []imap.Seq{imap.Seq{Start: uint32(2), Stop: uint32(11)}}},
	},
	getSequenceSetTestCase{
		&Config{Page: 0, PageSize: 10},
		&imap.MailboxStatus{Messages: 15},
		&imap.SeqSet{Set: []imap.Seq{imap.Seq{Start: uint32(6), Stop: uint32(15)}}},
	},
	getSequenceSetTestCase{
		&Config{Page: 1, PageSize: 10},
		&imap.MailboxStatus{Messages: 15},
		&imap.SeqSet{Set: []imap.Seq{imap.Seq{Start: uint32(1), Stop: uint32(5)}}},
	},
	getSequenceSetTestCase{
		&Config{Page: 1, PageSize: 10},
		&imap.MailboxStatus{Messages: 5},
		&imap.SeqSet{Set: []imap.Seq{}},
	},
}

func TestGetSequenceSet(t *testing.T) {
	for _, test := range getSequenceSetTestCases {
		if output := getSequenceSet(test.mailbox, test.config); output.String() != test.expected.String() {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}

func FakeFetchFunc(set *imap.SeqSet, items []imap.FetchItem, messages chan *imap.Message) error {
	defer close(messages)
	for _, s := range set.Set {
		for i := s.Start; i <= s.Stop; i++ {
			messages <- &imap.Message{
				SeqNum:   i,
				Envelope: &imap.Envelope{Subject: fmt.Sprintf("Subject %d", i)},
			}
		}
	}

	return nil
}

func TestListMessages(t *testing.T) {
	messages, err := listMessages(FakeFetchFunc, &imap.MailboxStatus{Messages: 15}, &Config{Page: 0, PageSize: 10})

	if err != nil {
		t.Errorf("Error %q", err)
	}

	if len(messages) != 10 {
		t.Errorf("Expected 10 messages, got %d", len(messages))
	}

	if messages[0].SeqNum != 6 {
		t.Errorf("Expected first message to have SeqNum 6, got %d", messages[0].SeqNum)
	}

	if messages[9].SeqNum != 15 {
		t.Errorf("Expected last message to have SeqNum 15, got %d", messages[9].SeqNum)
	}
}

func FakeListFunc(ref, name string, mailboxes chan *imap.MailboxInfo) error {
	defer close(mailboxes)
	for i := 1; i <= 5; i++ {
		mailboxes <- &imap.MailboxInfo{
			Name: fmt.Sprintf("Mailbox %d", i),
		}
	}

	return nil
}

func TestListMailboxes(t *testing.T) {
	mailboxes, err := listMailboxes(FakeListFunc, 5)

	if err != nil {
		t.Errorf("Error %q", err)
	}

	if len(mailboxes) != 5 {
		t.Errorf("Expected 5 messages, got %d", len(mailboxes))
	}

	if mailboxes[0].Name != "Mailbox 1" {
		t.Errorf("Expected first mailbox to have name 'Mailbox 1', got %q", mailboxes[0].Name)
	}

	if mailboxes[4].Name != "Mailbox 5" {
		t.Errorf("Expected first mailbox to have name 'Mailbox 5', got %q", mailboxes[0].Name)
	}
}

type FakeImapClient struct{}

func (c *FakeImapClient) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	return &imap.MailboxStatus{Messages: 15}, nil
}

func (c *FakeImapClient) Login(username, password string) error {
	return nil
}

func (c *FakeImapClient) List(ref, name string, mailboxes chan *imap.MailboxInfo) error {
	defer close(mailboxes)

	mboxes := []*imap.MailboxInfo{
		&imap.MailboxInfo{Name: "Mailbox 1"},
		&imap.MailboxInfo{Name: "Mailbox 2"},
	}

	for _, mbox := range mboxes {
		mailboxes <- mbox
	}

	return nil
}

func (c *FakeImapClient) Fetch(set *imap.SeqSet, items []imap.FetchItem, messages chan *imap.Message) error {
	defer close(messages)

	result := []*imap.Message{
		&imap.Message{
			Envelope: &imap.Envelope{
				Subject: "Subject 1",
				From: []*imap.Address{
					&imap.Address{
						PersonalName: "John Doe",
						MailboxName:  "john",
						HostName:     "example.com",
					},
				},
				To: []*imap.Address{
					&imap.Address{
						PersonalName: "Jane Doe",
						MailboxName:  "jane",
						HostName:     "example.com",
					},
				},
			},
		},
	}

	for _, msg := range result {
		messages <- msg
	}

	return nil
}

func TestClientListMailboxes(t *testing.T) {
	client := &Client{
		Config:     Config{PageSize: 10, Page: 0},
		IMAPClient: &FakeImapClient{},
	}

	mailboxes, err := client.GetMailboxes()

	if err != nil {
		t.Errorf("Error %q", err)
	}

	if len(mailboxes) != 2 {
		t.Errorf("Expected 2 mailboxes, got %d", len(mailboxes))
	}

	if mailboxes[0].Name != "Mailbox 1" {
		t.Errorf("Expected first mailbox to have name 'Mailbox 1', got %q", mailboxes[0].Name)
	}

	if mailboxes[1].Name != "Mailbox 2" {
		t.Errorf("Expected second mailbox to have name 'Mailbox 2', got %q", mailboxes[1].Name)
	}
}

func TestClientListMessages(t *testing.T) {
	client := &Client{
		Config:     Config{PageSize: 10, Page: 0},
		IMAPClient: &FakeImapClient{},
	}

	mailbox := model.Mailbox{Name: "Mailbox 1"}
	messages, err := client.GetMessages(mailbox)

	if err != nil {
		t.Errorf("Error %q", err)
	}

	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if messages[0].Subject != "Subject 1" {
		t.Errorf("Expected first message to have subject 'Subject 1', got %q", messages[0].Subject)
	}

	if messages[0].From.Name != "John Doe" {
		t.Errorf("Expected first message to have from 'John Doe', got %q", messages[0].From.Name)
	}

	if messages[0].From.Email != "john@example.com" {
		t.Errorf("Expected first message to have from `john@example.com`, got %q", messages[0].From.Email)
	}

	if messages[0].To[0].Name != "Jane Doe" {
		t.Errorf("Expected first message to have to 'Jane Doe', got %q", messages[0].To[0].Name)
	}

	if messages[0].To[0].Email != "jane@example.com" {
		t.Errorf("Expected first message to have to `jane@example.com`, got %q", messages[0].To[0].Email)
	}
}
