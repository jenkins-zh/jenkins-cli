package center

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/app/cmd"
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("center command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		cmd.rootCmd.SetArgs([]string{})
		cmd.rootOptions.Jenkins = ""
		cmd.rootOptions.ConfigFile = "test.yaml"

		centerOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		cmd.rootCmd.SetArgs([]string{})
		os.Remove(cmd.rootOptions.ConfigFile)
		cmd.rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should success", func() {
			data, err := cmd.GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(cmd.rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			requestCrumb, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/api/json", nil)
			requestCrumb.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"version":"0"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(requestCrumb)).Return(responseCrumb, nil)

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/updateCenter/api/json?pretty=false&depth=1", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"RestartRequiredForCompletion":true}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

			cmd.rootCmd.SetArgs([]string{"center"})
			_, err = cmd.rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})
	})
})
