package models

type EmailOpts struct {
	From             string
	SenderName       string
	To               []string
	Subject          string
	Body             string
	Attachments      []*Attachment
	HtmlTemplateBody string
	Tags             []*Mapper
	RawMessage       []byte
}

type Attachment struct {
	FileName    string
	Data        []byte
	Format      string
	ContentType string
}

type FileManageOpts struct {
	FileName    string
	Data        []byte
	Path        string
	ContentType string
	OwnedByMe   bool
}
