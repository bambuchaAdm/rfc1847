package rfc1847

import (
	"testing"
	"io"
)

func TestEncryptedMessage(t * testing.T){
	output, input := io.Pipe()
	message := NewEncryptedMessage(input, "pgp")

	if message == nil {
		t.Fail()
	}
	
	go func(){
		var err error = nil
		buf := make([]byte, 1024)
		for err == nil {

			
			buf = make([]byte, 1024)
			_, err = output.Read(buf)
			t.Log(string(buf), err)
			if err == io.EOF {
				t.Error("First EOF")
			}
		}
	}()
	
	message.Control.Write([]byte("Version 1"))
	message.Control.Close()
	message.Encrypted.Write([]byte("Some encrypted Message"))
	message.Encrypted.Close()
}
