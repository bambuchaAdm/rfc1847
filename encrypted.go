package rfc1847

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
)

var encryptedMIMEType = "multipart"
var encryptedMIMESubType = "encrypted"
var encryptedMIMEContentType = encryptedMIMEType + "/" + encryptedMIMESubType

type EncryptedMessage struct {
	w * multipart.Writer
	protocol string
	Control io.WriteCloser
	Encrypted io.WriteCloser
	Header textproto.MIMEHeader
}

func protocolHeader(protocol string) textproto.MIMEHeader  {
	return map[string][]string{ "Content-Type": []string{ protocol }}
}

func bodyHeader() textproto.MIMEHeader  {
	return map[string][]string{ "Content-Type": []string{ "application/octet-stream" }}
}

func NewEncryptedMessage(w io.Writer, protocol string) * EncryptedMessage {
	controlOutput, controlInput := io.Pipe()
	encryptedOutput, encryptedInput := io.Pipe()

	multipartWriter := multipart.NewWriter(w)

	header := map[string][]string{"Content-Type" : []string{
		"multipart/encrypted",
		fmt.Sprintf("boundary=\"%s\"", multipartWriter.Boundary()),
		fmt.Sprintf("protocol=\"%s\"", protocol)}	}

	message := &EncryptedMessage{
		w: multipartWriter,
		protocol: protocol,
		Control: controlInput,
		Encrypted: encryptedInput,
	  Header : header }

	go func(control, encrypted io.Reader, m * EncryptedMessage){		
		if controlPart, err := m.w.CreatePart(protocolHeader(m.protocol)); err == nil {
			if _, err := io.Copy(controlPart, control); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
		if encryptPart, err := m.w.CreatePart(bodyHeader()); err == nil {
			if _, err := io.Copy(encryptPart, encrypted); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
		if err := m.w.Close(); err != nil {
			panic(err)
		}
	}(controlOutput, encryptedOutput, message)

	return message
}
