package util

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Password util test", func() {
	var (
		ctrl   *gomock.Controller
		length int
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		length = 3
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("password length", func() {
			pass := GeneratePassword(length)
			Expect(len(pass)).To(Equal(length))
		})

		It("Different length", func() {
			length = 6
			pass := GeneratePassword(length)
			Expect(len(pass)).To(Equal(length))
		})

		It("Negative length", func() {
			length = -1
			pass := GeneratePassword(length)
			Expect(len(pass)).To(Equal(0))
		})

		It("Zero length", func() {
			length = 0
			pass := GeneratePassword(length)
			Expect(len(pass)).To(Equal(length))
		})
	})
})
