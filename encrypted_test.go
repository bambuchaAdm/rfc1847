package rfc1847

import (
	"testing"
	"bytes"
	"strings"
	"mime/multipart"
	"fmt"
	"time"
	"io"
)

var exampleProtocol = "pgp"

func exampleMessage(t * testing.T) (* EncryptedMessage, * bytes.Buffer) {
	buffer := new(bytes.Buffer)
	message := NewEncryptedMessage(buffer, exampleProtocol)

	if message == nil {
		t.Fail()
	}

	if _, err := message.Control.Write([]byte("Version 1")); err != nil {
		t.Fail()
	}
	if err := message.Control.Close(); err != nil {
		t.Fail()
	}
	if _, err := message.Encrypted.Write([]byte("Some encrypted Message")); err != nil {
		t.Fail()
	}
	if err := message.Encrypted.Close(); err != nil {
		t.Fail()
	}

	time.Sleep(1 * time.Microsecond)

	return message, buffer
}

func TestEncryptedMessageHeader(t * testing.T){
	message,_ := exampleMessage(t)
	headers := message.Header

	contentType := headers["Content-Type"]
	
	if contentType[0] != "multipart/encrypted" {
		t.Errorf("Content type is not multipart/encrypted now is '%s'", contentType[0] )
	}
	if ! strings.HasPrefix(contentType[1], "boundary=") {
		t.Error("Content type has no boundary")
	}
	if contentType[2] != fmt.Sprintf("protocol=\"%s\"", exampleProtocol) {
		t.Error("Content type has no protocol")
	}
}

func TestMetadata(t * testing.T){
	message, buffer := exampleMessage(t)
	t.Log(buffer)
	reader := multipart.NewReader(buffer, message.w.Boundary())
	firstPart, err := reader.NextPart()
	contentType := firstPart.Header["Content-Type"][0]
	if err != nil || contentType != exampleProtocol {
		t.Errorf("Content-Type is '%s'. Should be '%s'", contentType, exampleProtocol)
	}
}

func TestEncryptedData(t * testing.T){
	message, buffer := exampleMessage(t)
	t.Log(buffer)
	reader := multipart.NewReader(buffer, message.w.Boundary())
	_, err := reader.NextPart()
	secondPart, err := reader.NextPart()
	contentType := secondPart.Header["Content-Type"][0]
	if err != nil || contentType  != "application/octet-stream" {
		t.Error("Content-Type is '%s'. Should be '%s'", contentType, "application/octet-stream")
	}
}

func TestOnlyTwoPart(t * testing.T){
	message, buffer := exampleMessage(t)
	reader := multipart.NewReader(buffer, message.w.Boundary())
	part, err := reader.NextPart()
	part, err = reader.NextPart()
	part, err = reader.NextPart()
	
	if err != io.EOF {
		t.Log(part)
		t.Error("There is more parts there should be")
	}
}
