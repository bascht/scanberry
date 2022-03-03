package scan

import (
	"reflect"
	"testing"
	"time"
)

func TestDocumentFullName(t *testing.T) {
	document := Document{
		Id:     "test-uuid",
		Name:   "Here-Be-A-Test-Name",
		Date:   time.Date(2020, 10, 12, 20, 15, 00, 0, time.UTC),
		Duplex: true,
		Events: make(chan string),
	}

	if document.FullName() != "2020-10-12-201500-Here-Be-A-Test-Name" {
		t.Error("Full name incorrect" + document.FullName())
	}

	if document.FullNameWithExtension() != "2020-10-12-201500-Here-Be-A-Test-Name.pdf" {
		t.Error("Full name with extension incorrect" + document.FullName())
	}

	if reflect.DeepEqual(document.Args(), []string{"-d", "2020-10-12-201500-Here-Be-A-Test-Name.pdf"}) {
		t.Error("Incorrect arguments")
	}
}
