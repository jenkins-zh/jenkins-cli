package util

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gosuri/uiprogress"
)

const (
	// ContentType is for the http header of content type
	ContentType = "Content-Type"
	// ApplicationForm is for the form submit
	ApplicationForm = "application/x-www-form-urlencoded"
)

// HTTPDownloader is the downloader for http request
type HTTPDownloader struct {
	TargetFilePath     string
	URL                string
	ShowProgress       bool
	InsecureSkipVerify bool

	UserName string
	Password string

	Proxy     string
	ProxyAuth string

	Debug        bool
	RoundTripper http.RoundTripper
}

// SetProxy set the proxy for a http
func SetProxy(proxy, proxyAuth string, tr *http.Transport) (err error) {
	if proxy == "" {
		return
	}

	var proxyURL *url.URL
	if proxyURL, err = url.Parse(proxy); err != nil {
		return
	}

	tr.Proxy = http.ProxyURL(proxyURL)

	if proxyAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxyAuth))
		tr.ProxyConnectHeader = http.Header{}
		tr.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
	}
	return
}

// DownloadFile download a file with the progress
func (h *HTTPDownloader) DownloadFile() error {
	filepath, downloadURL, showProgress := h.TargetFilePath, h.URL, h.ShowProgress
	// Get the data
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}

	if h.UserName != "" && h.Password != "" {
		req.SetBasicAuth(h.UserName, h.Password)
	}
	var tr http.RoundTripper
	if h.RoundTripper != nil {
		tr = h.RoundTripper
	} else {
		trp := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: h.InsecureSkipVerify},
		}
		tr = trp
		if err = SetProxy(h.Proxy, h.ProxyAuth, trp); err != nil {
			return err
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if h.Debug {
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				ioutil.WriteFile("debug-download.html", data, 0664)
			}
		}
		return fmt.Errorf("invalidate status code: %d", resp.StatusCode)
	}

	writer := &ProgressIndicator{
		Title: "Downloading",
	}
	if showProgress {
		if total, ok := resp.Header["Content-Length"]; ok && len(total) > 0 {
			fileLength, err := strconv.ParseInt(total[0], 10, 64)
			if err == nil {
				writer.Total = float64(fileLength)
			}
		}
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	writer.Writer = out

	if showProgress {
		writer.Init()
	}

	// Write the body to file
	_, err = io.Copy(writer, resp.Body)
	return err
}

// ProgressIndicator hold the progress of io operation
type ProgressIndicator struct {
	Writer io.Writer
	Reader io.Reader
	Title  string

	// bytes.Buffer
	Total float64
	count float64
	bar   *uiprogress.Bar
}

// Init set the default value for progress indicator
func (i *ProgressIndicator) Init() {
	uiprogress.Start()             // start rendering
	i.bar = uiprogress.AddBar(100) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	i.bar.AppendCompleted()
	// i.bar.PrependElapsed()

	if i.Title != "" {
		i.bar.PrependFunc(func(_ *uiprogress.Bar) string {
			return fmt.Sprintf("%s: ", i.Title)
		})
	}
}

// Write writes the progress
func (i *ProgressIndicator) Write(p []byte) (n int, err error) {
	n, err = i.Writer.Write(p)
	i.setBar(n)
	return
}

// Read reads the progress
func (i *ProgressIndicator) Read(p []byte) (n int, err error) {
	n, err = i.Reader.Read(p)
	i.setBar(n)
	return
}

func (i *ProgressIndicator) setBar(n int) {
	i.count += float64(n)

	if i.bar != nil {
		i.bar.Set((int)(i.count * 100 / i.Total))
	}
}
