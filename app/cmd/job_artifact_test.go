package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("job artifact command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buildID      int
		jobName      string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		buildID = 1
		jobName = "fakeJob"

		jobArtifactOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("lack of arguments", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"job", "artifact"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("requires at least 1 arg(s), only received 0"))
		})

		It("should success", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareGetArtifacts(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobName, buildID)

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"job", "artifact", jobName, fmt.Sprintf("%d", buildID)})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`Name  Path  Size
a.log a.log 0
`))
		})

		It("should success, zero build id", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
			buildID = 0

			client.PrepareGetArtifacts(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", jobName, buildID)

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"job", "artifact", jobName, fmt.Sprintf("%d", buildID)})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`Name  Path  Size
a.log a.log 0
`))
		})

		It("should success, invalid build id", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"job", "artifact", jobName, "invalid"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid syntax"))
		})
	})
})
