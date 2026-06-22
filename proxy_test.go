package main

import (
	"testing"

	host "github.com/rickypc/native-messaging-host"
)

func TestIsEndOfStream(t *testing.T) {
	cases := []struct {
		name string
		h    host.H
		want bool
	}{
		{"end", host.H{"data": "end"}, true},
		{"other", host.H{"data": "more"}, false},
		{"empty", host.H{}, false},
		{"wrong-key", host.H{"foo": "end"}, false},
	}
	for _, c := range cases {
		h := c.h
		if got := isEndOfStream(&h); got != c.want {
			t.Errorf("%s: isEndOfStream = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestIsEndOfStream_Nil(t *testing.T) {
	if isEndOfStream(nil) {
		t.Fatal("nil response must not be end-of-stream")
	}
}
