package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job edit command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		jenkinsRoot  string
		username     string
		token        string

		script string
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

		script = "sample"
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
			client.PrepareForUpdatePipelineJob(roundTripper, jenkinsRoot, script, username, token)

			rootCmd.SetArgs([]string{"job", "edit", jobName, "--script", script})

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
			err = ioutil.WriteFile(tempFile.Name(), []byte(script), 0644)

			jobName := "test"
			client.PrepareForUpdatePipelineJob(roundTripper, jenkinsRoot, script, username, token)

			rootCmd.SetArgs([]string{"job", "edit", jobName, "--filename", tempFile.Name(), "--script", ""})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("edit with url param", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "test"
			client.PrepareForUpdatePipelineJob(roundTripper, jenkinsRoot, "sample", username, token)

			remoteJenkinsfileURL := "http://test"
			remoteJenkinsfileReq, _ := http.NewRequest("GET", remoteJenkinsfileURL, nil)
			remoteJenkinsfileResponse := &http.Response{
				StatusCode: 200,
				Request:    remoteJenkinsfileReq,
				Body:       ioutil.NopCloser(bytes.NewBufferString(script)),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(remoteJenkinsfileReq)).Return(remoteJenkinsfileResponse, nil)

			rootCmd.SetArgs([]string{"job", "edit", jobName, "--filename", "", "--script", "", "--url", remoteJenkinsfileURL})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})
	})
})
