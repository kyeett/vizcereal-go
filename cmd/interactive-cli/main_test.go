package main

import (
	"strings"
	"testing"
)

func TestSingleInput(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"+A", "Node 'A' added"},
		{"+B", "Node 'B' added"},

		{"-A", "Node 'A' does not exist"},

		{"A->B", "Connection removed"},
		{"A->B:0", "Connection 'A->B:0' added"},
		{"A->B:10", "Connection 'A->B:10' added"},

		{"asd", "Invalid operation"},
	}
	for _, test := range tests {
		n := Node{}
		if got := n.handleInput(test.input); got != test.want {
			t.Errorf("HandleInput(%q) = '%v', wanted %q", test.input, got, test.want)
		}
	}
}

func TestSequence(t *testing.T) {
	var tests = []struct {
		sequence []string
		want     string
	}{
		{
			sequence: []string{"+A", "+B"},
			want:     "A\nB\n",
		},
		{
			sequence: []string{"+A", "+B", "-B", "+C"},
			want:     "A\nC\n",
		},
	}
	for _, test := range tests {
		n := NewNode("edge", "global", "")

		for _, s := range test.sequence {
			n.handleInput(s)
		}

		if got := n.TextStatus(); got != test.want {
			t.Errorf("HandleInput(%q) = '%q', wanted %q", strings.Join(test.sequence, "\n"), got, test.want)
		}
	}
}
