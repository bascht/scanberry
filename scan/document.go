package scan

import (
	"time"
)

type Document struct {
	Id     string
	Name   string
	Duplex bool
	Date   time.Time
	Events chan Event
}

func (document Document) Args() []string {
	cmd := []string{}

	if document.Duplex {
		cmd = append(cmd, " -d")
	}
	cmd = append(cmd, document.FullName())

	return cmd
}

func (document Document) FullName() string {

	return document.Date.Format("2006-01-02-150405") + "-" + document.Name
}

func (document Document) FullNameWithExtension() string {
	return document.FullName() + ".pdf"
}
