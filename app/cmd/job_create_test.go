package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("job create command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buf          io.Writer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		jobCreateOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		var (
			jobPayload client.CreateJobPayload
			err        error
		)

		BeforeEach(func() {
			jobPayload = client.CreateJobPayload{
				Name: "jobName",
				Mode: "jobType",
			}

			var data []byte
			data, err = GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("create a job by the normal way", func() {
			client.PrepareForCreatePipelineJob(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobPayload)

			rootCmd.SetArgs([]string{"job", "create", jobPayload.Name, "--type", jobPayload.Mode})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})

		It("create a job by copy way", func() {
			jobPayload.From = "another-one"
			jobPayload.Mode = "copy"
			client.PrepareForCreatePipelineJob(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobPayload)

			rootCmd.SetArgs([]string{"job", "create", jobPayload.Name, "--copy", jobPayload.From})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})
	})
})
