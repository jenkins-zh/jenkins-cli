package util

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gosuri/uiprogress"
)

const (
	CONTENT_TYPE = "Content-Type"
	APP_FORM     = "application/x-www-form-urlencoded"
)

type HTTPDownloader struct {
	TargetFilePath string
	URL            string
	ShowProgress   bool

	UserName string
	Password string

	Debug bool
}

// DownloadFile download a file with the progress
func (h *HTTPDownloader) DownloadFile() error {
	filepath, url, showProgress := h.TargetFilePath, h.URL, h.ShowProgress
	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if h.UserName != "" && h.Password != "" {
		req.SetBasicAuth(h.UserName, h.Password)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if h.Debug {
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				ioutil.WriteFile("debug-download.html", data, 0664)
			}
		}
		return fmt.Errorf("Invalidate status code: %d", resp.StatusCode)
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
	writer.Init()

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
		i.bar.PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("%s: ", i.Title)
		})
	}
}

func (i *ProgressIndicator) Write(p []byte) (n int, err error) {
	n, err = i.Writer.Write(p)
	i.setBar(n)
	return
}

func (i *ProgressIndicator) Read(p []byte) (n int, err error) {
	n, err = i.Reader.Read(p)
	i.setBar(n)
	return
}

func (i *ProgressIndicator) setBar(n int) {
	i.count += float64(n)
	i.bar.Set((int)(i.count * 100 / i.Total))
}
