package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "github.com/AlecAivazis/survey/v2/core"
	// "github.com/AlecAivazis/survey/v2/terminal"
)

var _ = Describe("job edit command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		jenkinsRoot  string
		username     string
		token        string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobEditOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		jenkinsRoot = "http://localhost:8080/jenkins"
		username = "admin"
		token = "111e3a2f0231198855dceaff96f20540a9"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("edit with script param", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "test"
			client.PrepareForUpdatePipelineJob(roundTripper, jenkinsRoot, "sample", username, token)

			rootCmd.SetArgs([]string{"job", "edit", jobName, "--script", "sample"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("edit with file param", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			tempFile, err := ioutil.TempFile("", "example")
			Expect(err).To(BeNil())
			defer os.Remove(tempFile.Name())
			err = ioutil.WriteFile(tempFile.Name(), []byte("sample"), 0644)

			jobName := "test"
			client.PrepareForUpdatePipelineJob(roundTripper, jenkinsRoot, "sample", username, token)

			rootCmd.SetArgs([]string{"job", "edit", jobName, "--filename", tempFile.Name(), "--script", ""})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})
	})
})
