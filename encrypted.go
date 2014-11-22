package rfc1847

import (
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
}

func protocolHeader(protocol string) textproto.MIMEHeader  {
	contentType := make([]string,1)
	contentType[0] = protocol
	return map[string][]string{ "ContentType": contentType }
}

func bodyHeader() textproto.MIMEHeader  {
	contentType := make([]string,1)
	contentType[0] = "application/octet-stream"
	return map[string][]string{ "ContentType": contentType }
}

func NewEncryptedMessage(w io.WriteCloser, protocol string) * EncryptedMessage {
	controlOutput, controlInput := io.Pipe()
	encryptedOutput, encryptedInput := io.Pipe()

	multipartWriter := multipart.NewWriter(w)

	message := &EncryptedMessage{w: multipartWriter,
		protocol: protocol,
		Control: controlInput,
		Encrypted: encryptedInput	}

	go func(control, encrypted io.Reader, m * EncryptedMessage){
		controlPart, err := m.w.CreatePart(protocolHeader(m.protocol))
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(controlPart, control); err != nil {
			panic(err)
		}
		encryptPart, err := m.w.CreatePart(bodyHeader())
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(encryptPart, encrypted); err != nil {
			panic(err)
		}
		if err := m.w.Close(); err != nil { panic(err) }
		w.Close()
	}(controlOutput, encryptedOutput, message)

	return message
}
