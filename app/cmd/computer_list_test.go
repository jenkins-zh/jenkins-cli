package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-client/pkg/computer"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-client/pkg/mock/mhttp"
)

var _ = Describe("computer list command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buf          *bytes.Buffer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		computerListOption.RoundTripper = roundTripper
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		var (
			err error
		)

		BeforeEach(func() {
			var data []byte
			data, err = GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("should success", func() {
			computer.PrepareForComputerListRequest(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"computer", "list"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})

		It("with a fake jenkins as option", func() {
			rootCmd.SetArgs([]string{"computer", "list", "--jenkins", "fake"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("cannot found the configuration"))
		})
	})
})
