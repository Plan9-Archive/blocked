package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"goplan9.googlecode.com/hg/plan9/acme"
)

// Assumes a 64 byte line
func block2body(path string, w *acme.Win) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file %s: %s\n", flag.Arg(0), err)
		os.Exit(1)
	}

	var line [64]byte
	for {
		n, err := f.Read(line[:])
		if err == os.EOF {
			break
		}
		b := strings.TrimRight(string(line[:n]), " ")
		w.Printf("body", "%s\n", b)
	}
}


func body2block(w *acme.Win, path string) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	defer f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open file %s: %s\n", flag.Arg(0), err)
		os.Exit(1)
	}

	b, err := w.ReadAll("body")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read acme data file: %s\n", err)
		os.Exit(1)
	}

	lines := strings.FieldsFunc(string(b), func (c int) bool { return c == '\n' })
	for _, line := range lines {
		f.WriteString(line + strings.Repeat(" ", 64-len(line)))
	}
}
	
func main() {
	flag.Parse()

	w, err := acme.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating window: %s\n", err)
		os.Exit(1)
	}

	cwd, _ := os.Getwd()
	w.Name(cwd + "/" + flag.Arg(0))
	block2body(flag.Arg(0), w)

	c := w.EventChan()
	for {
		ev, ok := <-c
		if !ok {
			break
		}

		if ev.C1 == 'M' && ev.C2 == 'x' && string(ev.Text[:3]) == "Put" {
			body2block(w, flag.Arg(0))
		} else {
			w.WriteEvent(ev)
		}
	}
}
