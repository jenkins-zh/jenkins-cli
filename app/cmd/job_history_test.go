package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("job history command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		err          error
		jenkinsRoot  string
		username     string
		token        string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jobHistoryOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		jenkinsRoot = "http://localhost:8080/jenkins"
		username = "admin"
		token = "111e3a2f0231198855dceaff96f20540a9"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		err = os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobName := "fakeJob"

			client.PrepareForGetJob(roundTripper, jenkinsRoot, jobName, username, token)
			client.PrepareForGetBuild(roundTripper, jenkinsRoot, jobName, 1, username, token)
			client.PrepareForGetBuild(roundTripper, jenkinsRoot, jobName, 2, username, token)

			rootCmd.SetArgs([]string{"job", "history", jobName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`DisplayName Building Result
fake        false    
fake        false    
`))
		})
	})
})

var _ = Describe("ColorResult test", func() {
	It("should success", func() {
		Expect(ColorResult("unknown")).To(ContainSubstring("unknown"))
		Expect(ColorResult("SUCCESS")).To(ContainSubstring("SUCCESS"))
		Expect(ColorResult("FAILURE")).To(ContainSubstring("FAILURE"))
	})
})
