package rfc1847

import (
	"testing"
	"io"
	"bytes"
)

func TestEncryptedMessage(t * testing.T){
	output, input := io.Pipe()
	message := NewEncryptedMessage(input, "pgp")

	if message == nil {
		t.Fail()
	}

	buff := new(bytes.Buffer)
	
	go func(){
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
		
	}()

	sync := make(chan bool)

	go func(){
		buff.ReadFrom(output)
		t.Logf("\n%s", buff.String())
		sync <- true
	}()
	
	<- sync
	
}
