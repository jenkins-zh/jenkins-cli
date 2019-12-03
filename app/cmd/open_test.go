package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test open browser", func() {
	It("should success", func() {
		err := Open("fake://url", FakeExecCommandSuccess)
		Expect(err).To(BeNil())
	})
})
