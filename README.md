Read [HAR](https://en.wikipedia.org/wiki/HAR_(file_format)) ("HTTP Archive
format") archives.

Install with `go install zgo.at/har/cmd/unhar@latest`, which will put the binary
at `~/go/bin/unhar`.

You can also use it as a Go library: `har.FromFile("path.har")` will read in to
a `HAR` struct. You can then use `Extract()`.

API docs: https://godocs.io/zgo.at/har
