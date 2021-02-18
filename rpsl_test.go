package rpsl

import (
	"io"
	"strings"
	"testing"
)

func TestReadRoute(t *testing.T) {
	var rpsl = `route:       128.223.0.0/16
descr:       UONet
descr:       University of Oregon
descr:       Computing Center
descr:       Eugene, OR 97403-1212
descr:       USA
origin:      AS3582
mnt-by:      MAINT-AS3582
changed:     meyer@ns.uoregon.edu 19960222
source:      RADB`

	object, err := NewReader(strings.NewReader(rpsl)).Read()
	if err != nil {
		t.Fatalf("Read: %s", err)
	}
	if object.Class != "route" {
		t.Errorf("expected route, got %v", object.Class)
	}
	if l := len(object.Values["descr"]); l != 5 {
		t.Errorf("expected 5 descr lines, got %d", l)
	}
	if origin := object.Get("origin"); origin != "AS3582" {
		t.Errorf("expected origin `AS3582`, got %q", origin)
	}
}

func TestReadObjects(t *testing.T) {
	var rpsl = `# A Tale Of Two Records
aut-num: AS123
as-name: Foo Bar
descr:   Test

route:  127.0.0.0/8
descr:  Test route
origin: AS123`

	reader := NewReader(strings.NewReader(rpsl))

	if object, err := reader.Read(); err != nil {
		t.Errorf("Read aut-num: %s", err)
	} else if object == nil {
		t.Errorf("No aut-num returned")
	} else if object.Class != "aut-num" {
		t.Errorf("Expected class of `aut-num`, got %q", object.Class)
	}

	if object, err := reader.Read(); err != nil {
		t.Errorf("Read route: %s", err)
	} else if object == nil {
		t.Errorf("No route returned")
	} else if object.Class != "route" {
		t.Errorf("Expected class of `route`, got %q", object.Class)
	}

	if object, err := reader.Read(); err != io.EOF {
		t.Errorf("expected EOF")
	} else if object != nil {
		t.Errorf("expected nil Object")
	}
}

func TestWeirdComments(t *testing.T) {
	var rpsl = `route:         209.120.192.0/24
descr:         Yipes Communications Inc
origin:        AS6517
remarks:       MIA-VisionLab-NET
notify:        Peering@yipes.com
mnt-by:        MAINT-AS6517
changed:       dlim@yipes.com 20011011
source:        LEVEL3
               #delete:       juhlson@yipes.com no longer yipes customer`

	reader := NewReader(strings.NewReader(rpsl))

	if object, err := reader.Read(); err != nil {
		t.Errorf("Read route: %s", err)
	} else if object == nil {
		t.Errorf("No route returned")
	} else if object.Class != "route" {
		t.Errorf("Expected class of `aut-num`, got %v", object.Class)
	} else if len(object.Values["source"]) != 1 ||
		object.Values["source"][0] != "LEVEL3" {
		t.Errorf("Expected 'LEVEL3', got %v", object.Values["source"])
	}
}
