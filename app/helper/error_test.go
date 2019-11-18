package helper

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"os"
)

var _ = Describe("error helper", func() {
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
