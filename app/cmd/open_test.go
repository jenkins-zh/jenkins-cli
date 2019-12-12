package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("test open", func() {
	var (
		err         error
		jenkinsName string
	)

	BeforeEach(func() {
		data, err := generateSampleConfig()
		Expect(err).To(BeNil())
		rootOptions.ConfigFile = "test.yaml"
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
		openOption.ExecContext = util.FakeExecCommandSuccess
		jenkinsName = "fake"
	})

	JustBeforeEach(func() {
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs([]string{"open", jenkinsName})
		_, err = rootCmd.ExecuteC()
	})

	It("open a not exists Jenkins", func() {
		Expect(err).To(HaveOccurred())
		Expect(fmt.Sprint(err)).To(ContainSubstring("no URL found with Jenkins " + jenkinsName))
	})

	Context("give a right config", func() {
		BeforeEach(func() {
			jenkinsName = "yourServer"
		})

		It("should success", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.Remove(rootOptions.ConfigFile)
		})
	})
})
