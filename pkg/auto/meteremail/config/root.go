package config

import (
	"encoding/json"
	"fmt"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type ParsedDestination struct {
	Addr            string
	Username        string
	Password        string
	From            *mail.Address
	To              []*mail.Address
	SubjectTemplate *template.Template
	BodyTemplate    *template.Template
}

type Destination struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"` // defaults to From.Address
	jsontypes.Password
	// todo: TLS config

	From string   `json:"from,omitempty"` // RFC 5322 address, the address part used for auth against Host
	To   []string `json:"to,omitempty"`   // RFC 5322 address

	SendTime *jsontypes.Schedule `json:"sendTime,omitempty"` // defaults to midnight on Monday mornings: "0 0 * * 1"

	SubjectTemplate jsontypes.String `json:"subjectTemplate,omitempty"`
	BodyTemplate    jsontypes.String `json:"bodyTemplate,omitempty"`

	Parsed *ParsedDestination `json:"-"`

	Attachments map[string][]byte
}

type Source struct {
	Name      string `json:"name,omitempty"`
	Title     string `json:"title,omitempty"`
	Floor     string `json:"floor,omitempty"`
	Zone      string `json:"zone,omitempty"`
	Subsystem string `json:"subsystem,omitempty"`
}

type Root struct {
	auto.Config
	// Configuration information for how to send the email.
	Destination Destination `json:"destination,omitempty"`
	// Name of the device that implement OccupancySensor history trait and that is monitored.
	Source Source `json:"source,omitempty"`

	Now func() time.Time `json:"-"`

	ServerAddr string `json:"serverAddr,omitempty"`

	ElectricMeters []string `json:"electricMeters,omitempty"`

	WaterMeters []string `json:"waterMeters,omitempty"`
}

var (
	DefaultSendTime = jsontypes.MustParseSchedule("0 0 * * 1")
)

func ReadBytes(data []byte) (cfg Root, err error) {
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return
	}
	if cfg.Destination.Host == "" {
		err = fmt.Errorf("destination.host not specified")
		return
	}
	if len(cfg.Destination.To) == 0 {
		err = fmt.Errorf("destination.recipients is empty")
		return
	}
	// defaults
	if cfg.Destination.SendTime == nil {
		cfg.Destination.SendTime = DefaultSendTime
	}

	// validate email addresses
	parsed, err := cfg.Destination.Parse()
	if err != nil {
		return
	}
	cfg.Destination.Parsed = parsed
	for i, recipient := range cfg.Destination.To {
		_, err = mail.ParseAddress(recipient)
		if err != nil {
			err = fmt.Errorf("destination.recipients[%d] is invalid: %w", i, err)
			return
		}
	}
	return
}

func (d Destination) Parse() (*ParsedDestination, error) {
	p := &ParsedDestination{}
	var err error
	p.Addr = d.Addr()
	p.Username = d.Username
	p.Password, err = d.Password.Read()
	if err != nil {
		return nil, fmt.Errorf("destination.password: %w", err)
	}

	p.From, err = mail.ParseAddress(d.From)
	if err != nil {
		return nil, fmt.Errorf("destination.from: %w", err)
	}
	if p.Username == "" {
		p.Username = p.From.Address
	}

	for i, to := range d.To {
		a, err := mail.ParseAddress(to)
		if err != nil {
			return nil, fmt.Errorf("destination.to[%d]: %w", i, err)
		}
		p.To = append(p.To, a)
	}
	sort.Slice(p.To, func(i, j int) bool {
		return p.To[i].Address < p.To[j].Address
	})

	p.SubjectTemplate, err = d.ReadSubjectTemplate()
	if err != nil {
		return nil, fmt.Errorf("destination.subjectTemplate: %w", err)
	}
	p.BodyTemplate, err = d.ReadBodyTemplate()
	if err != nil {
		return nil, fmt.Errorf("destination.bodyTemplate: %w", err)
	}
	return p, nil
}

// Addr returns the combination of Host and Port, taking defaults into account.
// Suitable for smtp.Dial.
func (d Destination) Addr() string {
	p := d.Port
	if p == 0 {
		p = 587
	}
	return fmt.Sprintf("%s:%d", d.Host, p)
}

func (d Destination) ReadSubjectTemplate() (*template.Template, error) {
	s, err := d.SubjectTemplate.Read()
	if err != nil {
		return nil, err
	}
	if s == "" {
		s = DefaultEmailSubject
	}
	return template.New("subject").Parse(s)
}

func (d Destination) ReadBodyTemplate() (*template.Template, error) {
	s, err := d.BodyTemplate.Read()
	if err != nil {
		return nil, err
	}
	if s == "" {
		s = DefaultEmailBody
	}
	return template.New("body").
		Funcs(template.FuncMap{
			"printTime": func(t any) string {
				format := time.Stamp
				switch t := t.(type) {
				case time.Time:
					return t.Format(format)
				case *time.Time:
					return t.Format(format)
				case *timestamppb.Timestamp:
					return t.AsTime().Format(format)
				case nil:
					return ""
				default:
					return fmt.Sprintf("%v", t)
				}
			},
		}).
		Parse(s)
}

func (d *Destination) AttachFile(src string) error {
	if d.Attachments == nil {
		d.Attachments = make(map[string][]byte)
	}

	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	d.Attachments[fileName] = b
	return nil
}

const DefaultEmailSubject = `Smart Core Meter Readings {{.Now.Format "Jan 02, 2006"}}`
const DefaultEmailBody = `<html lang="en">
<head>
  <title>Once Centenary Way - Smart Core Meter Readings</title>
</head>
<body>
<section>
<h4>Please find below the total meter readings per each half floor. Attached is a detailed report which breaks these totals down into individual meter readings</h4>  
<h4>Electric Meter Readings
</h4>
<table>
<tbody>
<tr>	
	<th>Floor</th>
	<th>Zone</th>
	<th>Total Reading (kWh)</th>
</tr>
{{range .EnergySummaryReports}}
<tr>	
	<td>{{.Floor}}</td>	
	<td>{{.Zone}}</td>
	<td>{{.TotalReading}}</td>
</tr>
{{end}}
</tbody>
</table>
</section>
<section>
<h4>Water Meter Readings</h4>
<table>
<tbody>
<tr>
	<th>Floor</th>
	<th>Zone</th>
	<th>Reading (m³)</th>
</tr>
{{range .WaterSummaryReports}}
<tr>	
	<td>{{.Floor}}</td>
	<td>{{.Zone}}</td>
	<td>{{.TotalReading}}</td>
</tr>
{{end}}
</tbody>
</table>
</section>
</body>
</html>
`
