package cmd

import (
	"bytes"
	"fmt"
	"github.com/Netflix/go-expect"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("job build command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		jobName      string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		jobName = "fakeJob"
		jobBuildOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		ResetJobBuildOption()
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

			request, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/jenkins/job/%s/build", jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 201,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/crumbIssuer/api/json", nil)
			requestCrumb.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			rootCmd.SetArgs([]string{"job", "build", jobName, "-b", "true"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})

		It("with --param-entry and invalid --param", func() {
			var err error
			rootCmd.SetArgs([]string{"job", "build", jobName, "--param", "fake-param", "--param-entry", "key=value"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
		})

		It("with --param-entry", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForBuildWithParams(roundTripper, "http://localhost:8080/jenkins", jobName,
				"admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"job", "build", jobName, "--param-entry", "name=value", "-b", "true", "--param", ""})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func TestBuildJob(t *testing.T) {
	RunEditCommandTest(t, EditCommandTest{
		Args: []string{"job", "build", "fake", "-b=false"},
		ConfirmProcedure: func(c *expect.Console) {
			c.ExpectString("Are you sure to build job fake")
			c.SendLine("y")
			//c.ExpectEOF()
		},
		Procedure: func(c *expect.Console) {
			c.ExpectString("Edit your pipeline script")
			c.SendLine("")
			go c.ExpectEOF()
			time.Sleep(time.Millisecond)
			c.Send(`VGdi[{"Description":"","name":"name","Type":"StringParameterDefinition","value":"value","DefaultParameterValue":{"Description":"","Value":null}}]`)
			c.Send("\x1b")
			c.SendLine(":wq!")
		},
		CommonOption: &jobBuildOption.CommonOption,
		BatchOption:  &jobBuildOption.BatchOption,
	})
}

func RunEditCommandTest(t *testing.T, test EditCommandTest) {
	RunTest(t, func(stdio terminal.Stdio) (err error) {
		var data []byte
		rootOptions.ConfigFile = "test.yaml"
		data, err = generateSampleConfig()
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)

		ctrl := gomock.NewController(t)
		roundTripper := mhttp.NewMockRoundTripper(ctrl)

		var (
			url     = "http://localhost:8080/jenkins"
			jobName = "fake"
			user    = "admin"
			token   = "111e3a2f0231198855dceaff96f20540a9"
		)

		request, _ := http.NewRequest("GET", fmt.Sprintf("%s/job/%s/api/json",
			url, jobName), nil)
		request.SetBasicAuth(user, token)
		response := &http.Response{
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			Request:    request,
			Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"name":"fake",
"property" : [
    {
      "_class" : "hudson.model.ParametersDefinitionProperty",
      "parameterDefinitions" : [
        {
          "_class" : "hudson.model.StringParameterDefinition",
          "defaultParameterValue" : {
            "_class" : "hudson.model.StringParameterValue",
            "name" : "name",
            "value" : "value"
          },
          "description" : "",
          "name" : "name",
          "type" : "StringParameterDefinition"
        }
      ]
    }
]}
				`)),
		}
		roundTripper.EXPECT().
			RoundTrip(request).Return(response, nil)

		client.PrepareForBuildWithParams(roundTripper, url, jobName, user, token)

		jobBuildOption.RoundTripper = roundTripper
		test.BatchOption.Stdio = stdio
		test.CommonOption.Stdio = stdio
		rootCmd.SetArgs(test.Args)
		_, err = rootCmd.ExecuteC()
		return
	}, test.ConfirmProcedure, test.Procedure)
}
