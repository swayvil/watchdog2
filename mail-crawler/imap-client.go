package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"

	//	"fmt"
	//	"sync"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
)

type ImapClient struct {
	c        *client.Client
	location *time.Location
}

const timeLayout = "2006-01-02 15:04:05"

//var instance *ImapClient
//var once sync.Once

func NewImapClient() *ImapClient {
	// Connect to server
	c, err := client.DialTLS(GetConfigInstance().Imap.Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	//defer c.Logout()

	// Login
	if err := c.Login(GetConfigInstance().Imap.User, GetConfigInstance().Imap.Password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of messages: %d\n", mbox.Messages)
	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}

	location, err := time.LoadLocation(GetConfigInstance().Mail.Since.TimeZone)
	if err != nil {
		fmt.Println(err)
	}
	return &ImapClient{c, location}
}

func (imapClient *ImapClient) getSinceDateTime() time.Time {
	latest := GetLatestTimestampInserted()
	if latest == nil {
		fmt.Println("Search mails since default date")
		return time.Date(GetConfigInstance().Mail.Since.Year, time.Month(GetConfigInstance().Mail.Since.Month), GetConfigInstance().Mail.Since.Day, 0, 0, 0, 0, imapClient.location)
	}
	return *latest
}

func (imapClient *ImapClient) GetMessages() {
	criteria := imap.NewSearchCriteria()
	criteria.Since = imapClient.getSinceDateTime()
	uids, err := imapClient.c.Search(criteria)
	if err != nil {
		log.Println(err)
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
		matched, err := regexp.MatchString(GetConfigInstance().Mail.SubjectPattern, msg.Envelope.Subject)
		if err != nil {
			log.Println(err)
		}
		if matched {
			imapClient.parseMail(msg, section)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

func (imapClient *ImapClient) parseMail(msg *imap.Message, section *imap.BodySectionName) {
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
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			kind, _, _ := p.Header.ContentType()
			var camera string
			var timestamp time.Time

			if kind == "text/html" {
				read, err := ioutil.ReadAll(p.Body)
				if err != nil {
					log.Fatal(err)
				}
				camera, timestamp = imapClient.parseBody(string(read))
			}
			if kind == "image/jpeg" {
				read, err := ioutil.ReadAll(p.Body)
				if err != nil {
					log.Fatal(err)
				}
				InsertSnapshot(camera, timestamp, resizeImage(read), read)
			}
		}
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

func (imapClient *ImapClient) parseBody(body string) (string, time.Time) {
	r := regexp.MustCompile(GetConfigInstance().Mail.BodyPattern)
	match := r.FindAllStringSubmatch(body, -1)
	t := time.Date(strToInt(match[0][2]), time.Month(strToInt(match[0][3])), strToInt(match[0][4]), strToInt(match[0][5]), strToInt(match[0][6]), strToInt(match[0][7]), 0, imapClient.location)
	return string(match[0][1]), t
}

func (imapClient *ImapClient) Logout() {
	imapClient.c.Logout()
}
