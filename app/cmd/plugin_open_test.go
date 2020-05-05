package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var _ = Describe("plugin open test", func() {
	var (
		err error
	)

	BeforeEach(func() {
		pluginOpenOption.ExecContext = util.FakeExecCommandSuccess
		data, err := GenerateSampleConfig()
		Expect(err).To(BeNil())
		rootOptions.ConfigFile = "test.yaml"
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	JustBeforeEach(func() {
		rootCmd.SetArgs([]string{"plugin", "open"})
		_, err = rootCmd.ExecuteC()
	})

	AfterEach(func() {
		os.Remove(rootOptions.ConfigFile)
	})

	It("should success", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Context("without url", func() {
		BeforeEach(func() {
			pluginOpenOption.ExecContext = util.FakeExecCommandSuccess
			sampleConfig := getSampleConfig()
			sampleConfig.JenkinsServers[0].URL = ""
			data, err := yaml.Marshal(&sampleConfig)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("should failure", func() {
			Expect(err).To(HaveOccurred())
			Expect(fmt.Sprint(err)).To(ContainSubstring("no URL fond from"))
		})
	})
})
