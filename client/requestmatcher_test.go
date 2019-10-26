package client

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("user test", func() {
	Context("matchHeader", func() {
		var (
			left http.Header
			right http.Header
		)

		BeforeEach(func() {
			left = http.Header{}
			right = http.Header{}
		})

		It("two empty headers", func() {
			Expect(matchHeader(left, right)).To(Equal(true))
		})

		It("two same header with data", func() {
			left.Add("a", "a")
			right.Add("a", "a")

			Expect(matchHeader(left, right)).To(Equal(true))
		})

		It("different length of headers", func() {
			right.Add("a", "a")

			Expect(matchHeader(left, right)).To(Equal(false))
		})

		It("different value of headers", func() {
			right.Add("a", "a")
			left.Add("a", "b")

			Expect(matchHeader(left, right)).To(Equal(false))
		})

		It("different key of headers", func() {
			right.Add("a", "a")
			left.Add("b", "a")

			Expect(matchHeader(left, right)).To(Equal(false))
		})
	})
})