package uom

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func ExampleParseUlimit() {
	fmt.Println(ParseUlimit("nofile=512:1024"))
	fmt.Println(ParseUlimit("nofile=1024"))
	fmt.Println(ParseUlimit("cpu=2:4"))
	fmt.Println(ParseUlimit("cpu=6"))
}

func TestParseUlimitValid(t *testing.T) {
	u1 := &Ulimit{"nofile", 1024, 512}
	if u2, _ := ParseUlimit("nofile=512:1024"); *u1 != *u2 {
		t.Fatalf("expected %q, but got %q", u1, u2)
	}
}

func TestParseUlimitInvalidLimitType(t *testing.T) {
	if _, err := ParseUlimit("notarealtype=1024:1024"); err == nil {
		t.Fatalf("expected error on invalid ulimit type")
	}
}

func TestParseUlimitBadFormat(t *testing.T) {
	if _, err := ParseUlimit("nofile:1024:1024"); err == nil {
		t.Fatal("expected error on bad syntax")
	}

	if _, err := ParseUlimit("nofile"); err == nil {
		t.Fatal("expected error on bad syntax")
	}

	if _, err := ParseUlimit("nofile="); err == nil {
		t.Fatal("expected error on bad syntax")
	}
	if _, err := ParseUlimit("nofile=:"); err == nil {
		t.Fatal("expected error on bad syntax")
	}
	if _, err := ParseUlimit("nofile=:1024"); err == nil {
		t.Fatal("expected error on bad syntax")
	}
}

func TestParseUlimitHardLessThanSoft(t *testing.T) {
	if _, err := ParseUlimit("nofile=1024:1"); err == nil {
		t.Fatal("expected error on hard limit less than soft limit")
	}
	if _, err := ParseUlimit("nofile=-1:1024"); err == nil {
		t.Fatal("expected error on hard limit less than soft limit")
	}
}

func TestParseUlimitUnlimited(t *testing.T) {
	tt := []struct {
		in       string
		expected Ulimit
	}{
		{
			in:       "nofile=-1",
			expected: Ulimit{Name: "nofile", Soft: -1, Hard: -1},
		},
		{
			in:       "nofile=1024",
			expected: Ulimit{Name: "nofile", Soft: 1024, Hard: 1024},
		},
		{
			in:       "nofile=1024:-1",
			expected: Ulimit{Name: "nofile", Soft: 1024, Hard: -1},
		},
		{
			in:       "nofile=-1:-1",
			expected: Ulimit{Name: "nofile", Soft: -1, Hard: -1},
		},
	}

	for _, tc := range tt {
		t.Run(tc.in, func(t *testing.T) {
			u, err := ParseUlimit(tc.in)
			if err != nil {
				t.Fatalf("unexpected error when setting unlimited hard limit: %v", err)
			}
			if u.Name != tc.expected.Name {
				t.Fatalf("unexpected name. expected %s, got %s", tc.expected.Name, u.Name)
			}
			if u.Soft != tc.expected.Soft {
				t.Fatalf("unexpected soft limit. expected %d, got %d", tc.expected.Soft, u.Soft)
			}
			if u.Hard != tc.expected.Hard {
				t.Fatalf("unexpected hard limit. expected %d, got %d", tc.expected.Hard, u.Hard)
			}

		})
	}
}

func TestParseUlimitInvalidValueType(t *testing.T) {
	if _, err := ParseUlimit("nofile=asdf"); err == nil {
		t.Fatal("expected error on bad value type, but got no error")
	} else if _, ok := err.(*strconv.NumError); !ok {
		t.Fatalf("expected error on bad value type, but got `%s`", err)
	}

	if _, err := ParseUlimit("nofile=1024:asdf"); err == nil {
		t.Fatal("expected error on bad value type, but got no error")
	} else if _, ok := err.(*strconv.NumError); !ok {
		t.Fatalf("expected error on bad value type, but got `%s`", err)
	}
}

func TestParseUlimitTooManyValueArgs(t *testing.T) {
	if _, err := ParseUlimit("nofile=1024:1:50"); err == nil {
		t.Fatalf("expected error on more than two value arguments")
	}
}

func TestUlimitStringOutput(t *testing.T) {
	u := &Ulimit{"nofile", 1024, 512}
	if s := u.String(); s != "nofile=512:1024" {
		t.Fatal("expected String to return nofile=512:1024, but got", s)
	}
}

func TestGetRlimit(t *testing.T) {
	tt := []struct {
		ulimit Ulimit
		rlimit Rlimit
	}{
		{Ulimit{"core", 10, 12}, Rlimit{rlimitCore, 10, 12}},
		{Ulimit{"cpu", 1, 10}, Rlimit{rlimitCPU, 1, 10}},
		{Ulimit{"data", 5, 0}, Rlimit{rlimitData, 5, 0}},
		{Ulimit{"fsize", 2, 2}, Rlimit{rlimitFsize, 2, 2}},
		{Ulimit{"locks", 0, 0}, Rlimit{rlimitLocks, 0, 0}},
		{Ulimit{"memlock", 10, 10}, Rlimit{rlimitMemlock, 10, 10}},
		{Ulimit{"msgqueue", 9, 1}, Rlimit{rlimitMsgqueue, 9, 1}},
		{Ulimit{"nice", 9, 9}, Rlimit{rlimitNice, 9, 9}},
		{Ulimit{"nofile", 4, 100}, Rlimit{rlimitNofile, 4, 100}},
		{Ulimit{"nproc", 5, 5}, Rlimit{rlimitNproc, 5, 5}},
		{Ulimit{"rss", 0, 5}, Rlimit{rlimitRss, 0, 5}},
		{Ulimit{"rtprio", 100, 65}, Rlimit{rlimitRtprio, 100, 65}},
		{Ulimit{"rttime", 55, 102}, Rlimit{rlimitRttime, 55, 102}},
		{Ulimit{"sigpending", 14, 20}, Rlimit{rlimitSigpending, 14, 20}},
		{Ulimit{"stack", 1, 1}, Rlimit{rlimitStack, 1, 1}},
		{Ulimit{"stack", -1, -1}, Rlimit{rlimitStack, math.MaxUint64, math.MaxUint64}},
	}

	for _, te := range tt {
		res, err := te.ulimit.GetRlimit()
		if err != nil {
			t.Errorf("expected not to fail: %s", err)
		}
		if res.Type != te.rlimit.Type {
			t.Errorf("expected Type to be %d but got %d",
				te.rlimit.Type, res.Type)
		}
		if res.Soft != te.rlimit.Soft {
			t.Errorf("expected Soft to be %d but got %d",
				te.rlimit.Soft, res.Soft)
		}
		if res.Hard != te.rlimit.Hard {
			t.Errorf("expected Hard to be %d but got %d",
				te.rlimit.Hard, res.Hard)
		}

	}
}

func TestGetRlimitBadUlimitName(t *testing.T) {
	name := "bla"
	uLimit := Ulimit{name, 0, 0}
	if _, err := uLimit.GetRlimit(); err == nil {
		t.Error("expected error on bad Ulimit name")
	}
}
