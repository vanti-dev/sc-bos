package meteremail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/vanti-dev/sc-bos/pkg/auto/meteremail/config"
	"net/http"
	"net/smtp"
	"strings"
)

func sendEmail(dst config.Destination, attrs Attrs) error {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(dst.Attachments) > 0
	p := dst.Parsed

	// Subject
	var subj strings.Builder
	if err := p.SubjectTemplate.Execute(&subj, attrs); err != nil {
		return err
	}
	buf.WriteString(fmt.Sprintf("Subject: %s\n", subj.String()))

	// From
	buf.WriteString(fmt.Sprintf("From: %s\n", p.From.String()))

	// To
	var addrs []string
	for _, a := range p.To {
		addrs = append(addrs, a.Address)
		buf.WriteString(fmt.Sprintf("To: %s\n", a.String()))
	}

	buf.WriteString("MIME-Version: 1.0\n")

	// add body part
	mainMessageBoundary := "MixedBoundaryString"
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n", mainMessageBoundary))
	buf.WriteString("\n--" + mainMessageBoundary + "\n")

	attachmentBoundaryString := "AttachmentBoundaryString"
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/related; boundary=\"%s\"\n", attachmentBoundaryString))
	buf.WriteString("\n--" + attachmentBoundaryString + "\n")

	htmlMessageBoundary := "HtmlBoundaryString"
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\n", htmlMessageBoundary))
	buf.WriteString("\n--" + htmlMessageBoundary + "\n")

	buf.WriteString("Content-Type: text/html; charset=utf-8\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\n")
	if err := p.BodyTemplate.Execute(buf, attrs); err != nil {
		return err
	}
	buf.WriteString("--" + htmlMessageBoundary + "--")

	if withAttachments {
		boundary := attachmentBoundaryString

		for k, v := range dst.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
		}

		buf.WriteString(fmt.Sprintf("\n\n--%s--\n", boundary))
	}

	buf.WriteString("\n--" + mainMessageBoundary + "--")

	auth := smtp.PlainAuth("", p.Username, p.Password, dst.Host)
	addr := dst.Addr()
	err := smtp.SendMail(addr, auth, p.Username, addrs, buf.Bytes())

	if err != nil {
		return err
	}

	return nil
}
