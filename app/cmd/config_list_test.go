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
			Expect(buf.String()).To(Equal(`number name        url                           description
0      *yourServer http://localhost:8080/jenkins 
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
			Expect(buf.String()).To(Equal(`number name        url                           description
0      *yourServer http://localhost:8080/jenkins 012345678901234
`))
		})
	})
})
