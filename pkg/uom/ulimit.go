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
	"strconv"
	"strings"
)

// Ulimit is a human friendly version of Rlimit.
type Ulimit struct {
	Name string
	Hard int64
	Soft int64
}

// Rlimit specifies the resource limits, such as max open files.
type Rlimit struct {
	Type int    `json:"type,omitempty"`
	Hard uint64 `json:"hard,omitempty"`
	Soft uint64 `json:"soft,omitempty"`
}

const (
	// magic numbers for making the syscall
	// some of these are defined in the syscall package, but not all.
	// Also since Windows client doesn't get access to the syscall package, need to
	//	define these here
	rlimitAs         = 9
	rlimitCore       = 4
	rlimitCPU        = 0
	rlimitData       = 2
	rlimitFsize      = 1
	rlimitLocks      = 10
	rlimitMemlock    = 8
	rlimitMsgqueue   = 12
	rlimitNice       = 13
	rlimitNofile     = 7
	rlimitNproc      = 6
	rlimitRss        = 5
	rlimitRtprio     = 14
	rlimitRttime     = 15
	rlimitSigpending = 11
	rlimitStack      = 3
)

var ulimitNameMapping = map[string]int{
	//"as":         rlimitAs, // Disabled since this doesn't seem usable with the way Bhojpur.NET Platform initializes a Labni.
	"core":       rlimitCore,
	"cpu":        rlimitCPU,
	"data":       rlimitData,
	"fsize":      rlimitFsize,
	"locks":      rlimitLocks,
	"memlock":    rlimitMemlock,
	"msgqueue":   rlimitMsgqueue,
	"nice":       rlimitNice,
	"nofile":     rlimitNofile,
	"nproc":      rlimitNproc,
	"rss":        rlimitRss,
	"rtprio":     rlimitRtprio,
	"rttime":     rlimitRttime,
	"sigpending": rlimitSigpending,
	"stack":      rlimitStack,
}

// ParseUlimit parses and returns a Ulimit from the specified string.
func ParseUlimit(val string) (*Ulimit, error) {
	parts := strings.SplitN(val, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ulimit argument: %s", val)
	}

	if _, exists := ulimitNameMapping[parts[0]]; !exists {
		return nil, fmt.Errorf("invalid ulimit type: %s", parts[0])
	}

	var (
		soft int64
		hard = &soft // default to soft in case no hard was set
		temp int64
		err  error
	)
	switch limitVals := strings.Split(parts[1], ":"); len(limitVals) {
	case 2:
		temp, err = strconv.ParseInt(limitVals[1], 10, 64)
		if err != nil {
			return nil, err
		}
		hard = &temp
		fallthrough
	case 1:
		soft, err = strconv.ParseInt(limitVals[0], 10, 64)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("too many limit value arguments - %s, can only have up to two, `soft[:hard]`", parts[1])
	}

	if *hard != -1 {
		if soft == -1 {
			return nil, fmt.Errorf("ulimit soft limit must be less than or equal to hard limit: soft: -1 (unlimited), hard: %d", *hard)
		}
		if soft > *hard {
			return nil, fmt.Errorf("ulimit soft limit must be less than or equal to hard limit: %d > %d", soft, *hard)
		}
	}

	return &Ulimit{Name: parts[0], Soft: soft, Hard: *hard}, nil
}

// GetRlimit returns the RLimit corresponding to Ulimit.
func (u *Ulimit) GetRlimit() (*Rlimit, error) {
	t, exists := ulimitNameMapping[u.Name]
	if !exists {
		return nil, fmt.Errorf("invalid ulimit name %s", u.Name)
	}

	return &Rlimit{Type: t, Soft: uint64(u.Soft), Hard: uint64(u.Hard)}, nil
}

func (u *Ulimit) String() string {
	return fmt.Sprintf("%s=%d:%d", u.Name, u.Soft, u.Hard)
}
