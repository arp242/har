[![GoDoc](https://godoc.org/github.com/arp242/har?status.svg)](https://pkg.go.dev/arp242.net/har)

Read [HAR](https://en.wikipedia.org/wiki/HAR_(file_format)) ("HTTP Archive
format") archives.

Use `har.FromFile("path.har")` to read in to a `HAR` struct. You can then use
`Extract()`.

See `./cmd/unhar` for an example.

Install with `go get arp242.net/har/cmd/unhar`, which will put the binary at
`~/go/bin/unhar`.
