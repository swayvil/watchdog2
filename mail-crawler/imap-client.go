package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
)

type imapClient struct {
	c         *client.Client
	location  *time.Location
	importing bool // True = is already importing new mails
}

const timeLayout = "2006-01-02 15:04:05"

// Singleton
var instance *imapClient
var once sync.Once

func getImapClient() *imapClient {
	once.Do(func() {
		location, err := time.LoadLocation(getConfigInstance().Mail.Since.TimeZone)
		if err != nil {
			fmt.Println(err)
		}
		instance = &imapClient{nil, location, false}
	})
	return instance
}

func (imapClient *imapClient) connect() {
	// Connect to server
	c, err := client.DialTLS(getConfigInstance().Imap.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	//defer c.Logout()

	// Login
	if err := c.Login(getConfigInstance().Imap.User, getConfigInstance().Imap.Password); err != nil {
		log.Println("Can't connect to IMAP server with user: " + getConfigInstance().Imap.User)
		log.Fatal(err)
	}
	log.Println("Logged in")
	imapClient.selectMailBox("INBOX")
	imapClient.c = c
}

func (imapClient *imapClient) getSinceDateTime() time.Time {
	latest := selectLatestTimestampInserted()
	if latest == nil {
		fmt.Println("Search mails since default date")
		return time.Date(getConfigInstance().Mail.Since.Year, time.Month(getConfigInstance().Mail.Since.Month), getConfigInstance().Mail.Since.Day, 0, 0, 0, 0, imapClient.location)
	}
	return *latest
}

func (imapClient *imapClient) selectMailBox(mailbox string) {
	mbox, err := imapClient.c.Select(mailbox, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of messages: %d\n", mbox.Messages)
	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
}

func (imapClient *imapClient) importMessages() bool {
	if imapClient.importing { // Already importing
		return false
	}

	imapClient.importing = true
	imapClient.connect()

	criteria := imap.NewSearchCriteria()
	criteria.Since = imapClient.getSinceDateTime()
	uids, err := imapClient.c.Search(criteria)
	if err != nil {
		log.Println(err)
		imapClient.importing = false
		return false
	}
	fmt.Printf("Search result is: %d\n", len(uids))

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}

	messages := make(chan *imap.Message, 100) // Sends to a buffered channel block only when the buffer is full
	done := make(chan error, 1)
	go func() {
		done <- imapClient.c.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		matched, err := regexp.MatchString(getConfigInstance().Mail.SubjectPattern, msg.Envelope.Subject)
		if err != nil {
			log.Println("Bad subject: " + msg.Envelope.Subject)
			log.Println(err)
		}
		if matched {
			imapClient.parseMail(msg, section)
		} else {
			log.Printf("Bad subject, email ignored:  %s\n", msg.Envelope.Subject)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
	imapClient.logout()
	imapClient.importing = false
	return true
}

func (imapClient *imapClient) parseMail(msg *imap.Message, section *imap.BodySectionName) {
	r := msg.GetBody(section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	// Create a new mail reader
	mailReader, err := message.Read(r)
	if err != nil {
		log.Fatal(err)
	}

	// Print some info about the message
	// header := mailReader.Header
	// if date, err := header.Date(); err == nil {
	// 	log.Println("Date:", date)
	// }
	// if from, err := header.AddressList("From"); err == nil {
	// 	log.Println("From:", from)
	// }
	// if to, err := header.AddressList("To"); err == nil {
	// 	log.Println("To:", to)
	// }
	// if subject, err := header.Subject(); err == nil {
	// 	log.Println("Subject:", subject)
	// }

	if mr := mailReader.MultipartReader(); mr != nil {
		// This is a multipart message
		var camera string
		var timestamp time.Time
		var photo []byte
		var photoSmall []byte

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			kind, _, _ := p.Header.ContentType()

			if kind == "text/html" {
				read, err := ioutil.ReadAll(p.Body)
				if err != nil {
					log.Fatal(err)
				}
				camera, timestamp = imapClient.parseBody(string(read))
			}
			if kind == "image/jpeg" {
				photo, err = ioutil.ReadAll(p.Body)
				if err != nil {
					log.Fatal(err)
				}
				photoSmall = resizeImage(photo)
			}
		}
		insertSnapshot(camera, timestamp, photoSmall, photo)
	}

	// Process each message's part
	// for {
	// 	p, err := mailReader.NextPart()
	// 	if err == io.EOF {
	// 		log.Println("EOF")
	// 		//break
	// 	} else if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	switch h := p.Header.(type) {
	// 	case *mail.InlineHeader:
	// 		// This is the message's text (can be plain-text or HTML)
	// 		b, _ := ioutil.ReadAll(p.Body)
	// 		imapClient.parseBody(string(b))
	// 	case *mail.AttachmentHeader:
	// 		// This is an attachment
	// 		h.Filename()
	// 		filename, _ := h.Filename()
	// 		log.Println("Got attachment: %v", filename)
	// 	}
	//}
}

func strToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func (imapClient *imapClient) parseBody(body string) (string, time.Time) {
	r := regexp.MustCompile(getConfigInstance().Mail.BodyPattern)
	match := r.FindAllStringSubmatch(body, -1)
	t := time.Date(strToInt(match[0][2]), time.Month(strToInt(match[0][3])), strToInt(match[0][4]), strToInt(match[0][5]), strToInt(match[0][6]), strToInt(match[0][7]), 0, imapClient.location)
	return string(match[0][1]), t
}

func (imapClient *imapClient) logout() {
	imapClient.c.Logout()
}

// func (imapClient *imapClient) nbMailsToFetch() int {
// 	criteria := imap.NewSearchCriteria()
// 	criteria.Since = imapClient.getSinceDateTime()
// 	uids, err := imapClient.c.Search(criteria)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return len(uids)
// }

// eg. imapClient.deleteMessages(time.Date(2020, 5, 2, 0, 0, 0, 0, imapClient.location))
func (imapClient *imapClient) deleteMessages(mailbox string, to time.Time) {
	imapClient.selectMailBox(mailbox)
	criteria := imap.NewSearchCriteria()
	criteria.Before = to
	uids, err := imapClient.c.Search(criteria)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%d messages found to be deleted\n", len(uids))

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	// First mark the messages as deleted
	flags := []interface{}{imap.DeletedFlag}
	item := imap.FormatFlagsOp(imap.AddFlags, false)
	if err := imapClient.c.Store(seqset, item, flags, nil); err != nil {
		log.Fatal(err)
	}

	// Then delete them
	if err := imapClient.c.Expunge(nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Messages have been deleted")
}

func (imapClient *imapClient) printMailboxes() {
	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- imapClient.c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}
}
