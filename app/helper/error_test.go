package helper

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
)

var _ = Describe("error helper for StandardErrorMessage", func() {
	var (
		err error
		msg string
	)

	JustBeforeEach(func() {
		msg, _ = StandardErrorMessage(err)
	})

	It("err is nil", func() {
		Expect(msg).To(Equal(""))
	})

	Context("invalid host error", func() {
		BeforeEach(func() {
			err = url.InvalidHostError("host")
		})

		It("fake host", func() {
			Expect(msg).To(Equal("invalid character \"host\" in host name"))
		})
	})

	Context("url error", func() {
		BeforeEach(func() {
			err = &url.Error{
				Op:  "create",
				URL: "google.com",
				Err: fmt.Errorf("unknow"),
			}
		})

		It("url error", func() {
			Expect(msg).To(Equal("Unable to connect to the server: unknow"))
		})

		Context("connection refused", func() {
			BeforeEach(func() {
				err = &url.Error{
					Op:  "create",
					URL: "http://google.com",
					Err: fmt.Errorf("connection refused"),
				}
			})

			It("url error", func() {
				Expect(msg).To(Equal("The connection to the server google.com was refused - did you specify the right host or port?"))
			})
		})
	})

	Context("os path error", func() {
		BeforeEach(func() {
			err = &os.PathError{
				Op:   "create",
				Path: "fake-path",
				Err:  fmt.Errorf("unknow"),
			}
		})

		It("fake host", func() {
			Expect(msg).To(Equal("error: create fake-path: unknow"))
		})
	})
})

// MemoryPrinter only for test
type MemoryPrinter struct {
	Buffer *bytes.Buffer
}

// PrintErr print the error into memory
func (p *MemoryPrinter) PrintErr(i ...interface{}) {
	p.Buffer.WriteString(fmt.Sprintf("%v", i))
}

// Println print the object
func (p *MemoryPrinter) Println(i ...interface{}) {
	p.Buffer.WriteString(fmt.Sprintf("%v\n", i))
}

// Printf print against a format
func (p *MemoryPrinter) Printf(format string, i ...interface{}) {
	p.Buffer.WriteString(fmt.Sprintf(format, i...))
}

var _ = Describe("CheckErr", func() {
	var (
		printer *MemoryPrinter
		err     error
	)

	BeforeEach(func() {
		printer = &MemoryPrinter{
			Buffer: new(bytes.Buffer),
		}
	})

	JustBeforeEach(func() {
		CheckErr(printer, err)
	})

	It("error is nil", func() {
		Expect(printer.Buffer.String()).To(Equal(""))
	})

	Context("error from fmt", func() {
		BeforeEach(func() {
			err = fmt.Errorf("fake error")
		})

		It("fake error", func() {
			Expect(printer.Buffer.String()).To(Equal("[error: fake error]"))
		})
	})
})
