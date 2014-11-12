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

func NewEncryptedMessage(w io.Writer, protocol string) * EncryptedMessage {
	controlOutput, controlInput := io.Pipe()
	encryptedOutput, encryptedInput := io.Pipe()

	multipartWriter := multipart.NewWriter(w)

	message := &EncryptedMessage{w: multipartWriter,
		protocol: protocol,
		Control: controlInput,
		Encrypted: encryptedInput	}

	go func(control, encrypted io.Reader, m * EncryptedMessage){
		controlPart, _ := m.w.CreatePart(protocolHeader(m.protocol))
		io.Copy(controlPart, control)
		encryptPart, _ := m.w.CreatePart(bodyHeader())
		io.Copy(encryptPart, encrypted)
		m.w.Close()
	}(controlOutput, encryptedOutput, message)

	return message
}
