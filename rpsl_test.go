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

func TestComments(t *testing.T) {
	var rpsl = `#
# The contents of this file are subject to 
# AFRINIC Database Terms and Conditions
#
# http://www.afrinic.net/en/services
#`

	reader := NewReader(strings.NewReader(rpsl))

	object, err := reader.Read()
	if err != io.EOF {
		t.Errorf("Read comment: %s", err)
	} else if object != nil {
		t.Errorf("Expected nil object, got: %v", object)
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
               #delete:       juhlson@yipes.com no longer yipes customer

inetnum:        80.6.88.112 - 80.6.88.127
netname:        ASPIRE-HOUSING-LTD
descr:          NEWCASTLE UNDER LYME HOUSING
country:        GB
admin-c:        DUMY-RIPE
tech-c:         DUMY-RIPE
status:         ASSIGNED PA
mnt-by:         AS5089-MNT
created:        2003-04-17T12:25:21Z
last-modified:  2012-03-01T14:13:18Z
source:         RIPE #
remarks:        ****************************
remarks:        * THIS OBJECT IS MODIFIED

as-set:         AS-COFRACTAL
descr:          Cofractal, Inc.
remarks:        Customer ASN(s) for #2595510-56fa56ea-2a2b-4a79-b38b-cc8c24ad71d9
members:        AS17080
`

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

	if object, err := reader.Read(); err != nil {
		t.Errorf("Read inetnum: %s", err)
	} else if object == nil {
		t.Errorf("No inetnum returned")
	} else if object.Class != "inetnum" {
		t.Errorf("Expected class of `inetnum`, got %v", object.Class)
	} else if len(object.Values["source"]) != 1 ||
		object.Values["source"][0] != "RIPE " { //yep
		t.Errorf("Expected 'RIPE ', got %v", object.Values["source"])
	}

	if object, err := reader.Read(); err != nil {
		t.Errorf("Read as-set: %s", err)
	} else if object == nil {
		t.Errorf("No as-set returned")
	} else if object.Values["remarks"][0] != "Customer ASN(s) for #2595510-56fa56ea-2a2b-4a79-b38b-cc8c24ad71d9" {
		t.Errorf("ugh comments, got %v", object.Values["remarks"][0])
	}
}
