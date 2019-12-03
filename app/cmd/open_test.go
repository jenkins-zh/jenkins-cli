package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("test open", func() {
	var (
		err error
	)

	BeforeEach(func() {
		openOption.ExecContext = util.FakeExecCommandSuccess
	})

	JustBeforeEach(func() {
		rootCmd.SetArgs([]string{"open", "yourServer"})
		_, err = rootCmd.ExecuteC()
	})

	It("open a not exists Jenkins", func() {
		Expect(err).To(HaveOccurred())
		Expect(fmt.Sprint(err)).To(ContainSubstring("no URL found with Jenkins yourServer"))
	})

	Context("give a right config", func() {
		BeforeEach(func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			rootOptions.ConfigFile = "test.yaml"
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("should success", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.Remove(rootOptions.ConfigFile)
		})
	})
})
