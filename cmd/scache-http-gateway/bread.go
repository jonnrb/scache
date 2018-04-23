package main

import (
	"errors"
	"fmt"
	"html"
	"io"
)

var notFullyWritten = errors.New("not fully written")

func writeChain(w io.Writer, items ...func(io.Writer) error) error {
	for _, out := range items {
		if err := out(w); err != nil {
			return err
		}
	}
	return nil
}

func writeAll(w io.Writer, s string) error {
	n, err := io.WriteString(w, s)
	if err == nil && n != len(s) {
		err = notFullyWritten
	}
	return err
}

func stringWriter(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		return writeAll(w, s)
	}
}

func headerWriter(titles ...string) func(io.Writer) error {
	return func(w io.Writer) error {
		err := writeAll(w, `<!doctype html>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.0/css/bootstrap.min.css" integrity="sha384-9gVQ4dYFwwWSjIDZnLEWnxCjeSWFphJiwGPXr1jddIhOegiu1FwO5qRGvFXOdJZ4" crossorigin="anonymous">
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.0/js/bootstrap.min.js" integrity="sha384-uefMccjFJAIv6A+rW+L4AHf99KvxDjWSu1z9VI8SKNVmz4sk7buKt/6v9KI65qnm" crossorigin="anonymous"></script>

<div class="container h-100">
`)
		if err != nil {
			return err
		}
		for i, t := range titles {
			h := fmt.Sprintf("<h%d>%s</h%d>\n", i+1, html.EscapeString(t), i+1)
			if err := writeAll(w, h); err != nil {
				return err
			}
		}
		return nil
	}
}

func writeFooter(w io.Writer) error {
	return writeAll(w, `</div>
`)
}

func startCode(w io.Writer) error {
	return writeAll(w, `<pre><code>
`)
}

func endCode(w io.Writer) error {
	return writeAll(w, `</code></pre>
`)
}

func startList(w io.Writer) error {
	return writeAll(w, `<ul>
`)
}

func writeListLink(w io.Writer, title, href string) error {
	return writeAll(w, fmt.Sprintf(
		"<li><a href=\"%s\">%s</a>",
		html.EscapeString(href),
		html.EscapeString(title),
	))
}

func endList(w io.Writer) error {
	return writeAll(w, `</ul>
`)
}
