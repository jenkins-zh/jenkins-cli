package app

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("version test", func() {
	Context("basic function", func() {
		It("shoud success", func() {
			Expect(GetVersion()).To(Equal(""))
			Expect(GetCommit()).To(Equal(""))
			Expect(GetCombinedVersion()).NotTo(Equal(""))
		})
	})
})
