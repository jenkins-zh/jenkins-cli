package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("job stop command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobStopOption.RoundTripper = roundTripper
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
		It("should success, with batch mode", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fakeJob"
			buildID := 1
			request, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/jenkins/job/%s/%d/stop", jobName, buildID), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
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

			rootCmd.SetArgs([]string{"job", "stop", jobName, "1", "-b"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("stop the last build, with batch mode", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fakeJob"
			request, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/jenkins/job/%s/lastBuild/stop", jobName), nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
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

			rootCmd.SetArgs([]string{"job", "stop", jobName, "-b"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("stop the last build, with batch mode", func() {
			jobName := "fakeJob"
			rootCmd.SetArgs([]string{"job", "stop", jobName, "not-number", "-b"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err := rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())

			Expect(buf.String()).To(ContainSubstring("Error: strconv.Atoi: parsing"))
		})
	})
})
