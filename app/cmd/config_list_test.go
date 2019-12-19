package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config list command", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		config = nil
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`Name       URL                           Description
yourServer http://localhost:8080/jenkins 
`))
		})

		It("with long description", func() {
			config := getSampleConfig()
			config.JenkinsServers[0].Description = "01234567890123456789"
			data, err := yaml.Marshal(&config)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`Name       URL                           Description
yourServer http://localhost:8080/jenkins 01234567890123456789
`))
		})

		It("print the list of PreHooks", func() {
			sampleConfig := getSampleConfig()
			config = &sampleConfig

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list", "--config", "PreHooks"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`Path Command
`))
		})

		It("print the list of PostHooks", func() {
			sampleConfig := getSampleConfig()
			config = &sampleConfig

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list", "--config", "PostHooks"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`Path Command
`))
		})

		It("print the list of Mirrors", func() {
			sampleConfig := getSampleConfig()
			config = &sampleConfig

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list", "--config", "Mirrors"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`Name     URL
default  http://mirrors.jenkins.io/
tsinghua https://mirrors.tuna.tsinghua.edu.cn/jenkins/
huawei   https://mirrors.huaweicloud.com/jenkins/
tencent  https://mirrors.cloud.tencent.com/jenkins/
`))
		})

		It("print the list of PluginSuites", func() {
			sampleConfig := getSampleConfig()
			config = &sampleConfig

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list", "--config", "PluginSuites"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(`Name Description
`))
		})

		It("print the list of a fake type", func() {
			sampleConfig := getSampleConfig()
			config = &sampleConfig

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"config", "list", "--config", "fake"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("unknow config"))
		})
	})
})
