package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("color test", func() {
	Context("should success", func() {
		It("should success", func() {
			Expect(ColorInfo("")).To(Equal(""))
			Expect(ColorStatus("")).To(Equal(""))
			Expect(ColorWarning("")).To(Equal(""))
			Expect(ColorError("")).To(Equal(""))
			Expect(ColorBold("")).To(Equal(""))
			Expect(ColorAnswer("")).To(Equal(""))
		})
	})
})
