package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job search command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		rootURL      string
		username     string
		token        string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobParamOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		rootURL = "http://localhost:8080/jenkins"
		username = "admin"
		token = "111e3a2f0231198855dceaff96f20540a9"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases, need RoundTripper", func() {
		It("without parameters", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fake"
			client.PrepareForGetJob(roundTripper, rootURL, jobName, username, token)
			rootCmd.SetArgs([]string{"job", "param", jobName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("with one parameter", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fake"
			client.PrepareForGetJobWithParams(roundTripper, rootURL, jobName, username, token)
			rootCmd.SetArgs([]string{"job", "param", jobName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("[{\"Description\":\"\",\"name\":\"name\",\"Type\":\"StringParameterDefinition\",\"value\":\"\",\"DefaultParameterValue\":{\"Description\":\"\",\"Value\":\"jake\"}}]\n"))
		})

		It("with one parameter, output with indent", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fake"
			client.PrepareForGetJobWithParams(roundTripper, rootURL, jobName, username, token)
			rootCmd.SetArgs([]string{"job", "param", jobName, "--indent=true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("[\n {\n  \"Description\": \"\",\n  \"name\": \"name\",\n  \"Type\": \"StringParameterDefinition\",\n  \"value\": \"\",\n  \"DefaultParameterValue\": {\n   \"Description\": \"\",\n   \"Value\": \"jake\"\n  }\n }\n]\n"))
		})
	})
})
