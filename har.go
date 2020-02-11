package har

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Har archive.
type Har struct {
	File string `json:"-"` // File this was read from.

	Log struct {
		Version string `json:"version"` // "1.2"
		Creator struct {
			Name    string `json:"name"`    // "Firefox"
			Version string `json:"version"` // "72.0.2"
		} `json:"creator"`
		Browser struct {
			Name    string `json:"name"`    // "Firefox"
			Version string `json:"version"` // "72.0.2"
		} `json:"browser"`
		Pages []struct {
			StartedDateTime time.Time `json:"startedDateTime"`
			ID              string    `json:"id"` // "page_1"
			PageTimings     struct {
				OnContentLoad int `json:"onContentLoad"`
				OnLoad        int `json:"onLoad"`
			} `json:"pageTimings"`
		} `json:"pages"`

		Entries []Entry `json:"entries"`
	} `json:"log"`
}

type Entry struct {
	PageRef         string    `json:"pageRef"` // "page_1"
	StartedDateTime time.Time `json:"startedDateTime"`

	Request struct {
		BodySize    int    `json:"bodySize"`
		HeadersSize int    `json:"headersSize"`
		Method      string `json:"method"`
		URL         string `json:"url"`         // "http://localhost/..."
		HTTPVersion string `json:"httpVersion"` // "HTTP/1.1"
		Headers     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"headers"`
		Cookies []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"headers"`
		QueryString []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"queryString"`
	} `json:"request"`

	Response struct {
		HeadersSize int    `json:"headersSize"`
		BodySize    int    `json:"bodySize"`
		Status      int    `json:"status"`
		StatusText  string `json:"statusText"`  // "OK"
		HTTPVersion string `json:"httpVersion"` // "HTTP/1.1"
		RedirectURL string `json:"redirectURL"`
		Headers     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"headers"`
		Cookies []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"headers"`
		Content struct {
			Encoding string `json:"encoding"`
			Size     int    `json:"size"`
			Text     string `json:"text"`
		} `json:"content"`
	} `json:"response"`

	Cache struct {
		// TODO
	} `json:"cache"`

	Timings struct {
		Blocked int `json:"blocked"`
		DNS     int `json:"dns"`
		Connect int `json:"connect"`
		SSL     int `json:"ssl"`
		Send    int `json:"send"`
		Wait    int `json:"wait"`
		Receive int `json:"receive"`
	} `json:"timings"`

	Time            int    `json:"time"`
	SecurityState   string `json:"_securityState"`  // TODO
	ServerIPAddress string `json:"serverIPAddress"` // "::1"
	Connection      string `json:"connection"`      // "80"
}

// FromFile reads a file in to a Har struct.
func FromFile(f string) (*Har, error) {
	fp, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	data, err := ioutil.ReadAll(fp)
	fp.Close()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var h Har
	err = json.Unmarshal(data, &h)
	if err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}

	h.File = f
	return &h, nil
}

// Extract all the files.
func (h *Har) Extract(verbose bool) error {
	root := filepath.Base(h.File)
	if ext := filepath.Ext(root); ext != "" {
		root = root[:len(root)-len(ext)]
	}

	for _, e := range h.Log.Entries {
		path := root + "/" + strings.TrimPrefix(strings.TrimPrefix(e.Request.URL, "http://"), "https://")
		if verbose {
			fmt.Println("  ", path)
		}

		for _, h := range e.Response.Headers {
			if h.Key == "Content-Disposition" {
				for _, v := range strings.Split(h.Value, ";") {
					v := strings.TrimSpace(v)
					if strings.HasPrefix(v, "filename=") {
						v = filepath.Clean(strings.Trim(v, `"`))
						v = strings.ReplaceAll(v, "/", "")
						path = filepath.Dir(path) + "/" + v
					}
				}
				break
			}
		}
		if strings.HasSuffix(path, "/") {
			path += "index.html"
		}

		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}

		var resp []byte
		if e.Response.Content.Encoding == "base64" {
			resp, err = base64.StdEncoding.DecodeString(e.Response.Content.Text)
			if err != nil {
				return fmt.Errorf("base64: %w", err)
			}
		} else {
			resp = []byte(e.Response.Content.Text)
		}

		err = ioutil.WriteFile(path, resp, 0644)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	fmt.Printf("Extracted %d files\n", len(h.Log.Entries))
	return nil
}
