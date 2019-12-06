package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("casc open test", func() {
	var (
		err error
	)

	BeforeEach(func() {
		cascOpenOption.ExecContext = util.FakeExecCommandSuccess
		data, err := generateSampleConfig()
		Expect(err).To(BeNil())
		rootOptions.ConfigFile = "test.yaml"
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	JustBeforeEach(func() {
		rootCmd.SetArgs([]string{"casc", "open"})
		_, err = rootCmd.ExecuteC()
	})

	AfterEach(func() {
		os.Remove(rootOptions.ConfigFile)
	})

	It("should success", func() {
		Expect(err).NotTo(HaveOccurred())
	})
})
