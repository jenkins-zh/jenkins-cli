package i18n

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("test LoadTranslations", func() {
	var (
		root        string
		getLangFunc func() string
		err         error
	)

	JustBeforeEach(func() {
		err = LoadTranslations(root, getLangFunc)
	})

	AfterEach(func() {
		root = ""
		getLangFunc = nil
	})

	It("default param", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Context("unknown language", func() {
		BeforeEach(func() {
			getLangFunc = func() string {
				return "fake"
			}
		})

		It("should not have error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("given invalid environment of language", func() {
		var (
			osEnvErr error
		)

		BeforeEach(func() {
			osEnvErr = os.Setenv("LC_ALL", "zh_CN")
		})

		It("should not have error", func() {
			Expect(osEnvErr).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("given valid environment of language", func() {
		var (
			osEnvErr error
		)

		BeforeEach(func() {
			root = "jcli"
			osEnvErr = os.Setenv("LC_ALL", "zh_CN.utf-8")
		})

		It("should not have error", func() {
			Expect(osEnvErr).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("test i18n function T", func() {
	var (
		text   string
		args   []int
		result string
	)

	JustBeforeEach(func() {
		result = T(text, args...)
	})

	It("simple case, without args", func() {
		Expect(result).To(Equal(text))
	})

	Context("with args", func() {
		BeforeEach(func() {
			text = "fake %d"
			args = []int{1}
		})

		It("should success", func() {
			Expect(result).To(Equal(fmt.Sprintf(text, 1)))
		})
	})
})
