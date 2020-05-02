package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("computer launch command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buf          io.Writer
		err          error
		name         string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		computerLaunchOption.RoundTripper = roundTripper
		name = "fake"

		var data []byte
		data, err = GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("launch a default type of agent", func() {
		It("should success", func() {

			client.PrepareForLaunchComputer(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", name)

			rootCmd.SetArgs([]string{"computer", "launch", name})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})
	})

	Context("launch a jnlp agent", func() {
		var (
			fakeJar string
		)
		BeforeEach(func() {
			fakeJar = "fake-jar-content"
			computerLaunchOption.SystemCallExec = util.FakeSystemCallExecSuccess
			computerLaunchOption.LookPathContext = util.FakeLookPath

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/jnlpJars/agent.jar", nil)
			response := &http.Response{
				StatusCode: 200,
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(fakeJar)),
			}
			roundTripper.EXPECT().
				RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)

			secret := "fake-secret"
			client.PrepareForComputerAgentSecretRequest(roundTripper,
				"http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", name, secret)
		})

		It("should success", func() {
			rootCmd.SetArgs([]string{"computer", "launch", name, "--type", "jnlp", "--show-progress=false"})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
