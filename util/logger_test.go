package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("logger test", func() {
	Context("InitLogger", func() {
		It("basic test", func() {
			logger, err := InitLogger("warn")
			Expect(err).To(BeNil())
			Expect(logger).NotTo(BeNil())
		})
	})
})
