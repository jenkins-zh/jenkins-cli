package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("plugin search command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginSearchOption.RoundTripper = roundTripper
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
		It("should success, empty list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _, requestCenter, _, requestAvliable, _ := client.PrepareForEmptyAvaiablePluginList(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			requestCenter.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			requestAvliable.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "search", "fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("one plugin in the list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _, requestCenter, _, requestAvliable, _ := client.PrepareForOneAvaiablePlugin(roundTripper, "http://localhost:8080/jenkins")
			requestCenter.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			requestAvliable.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "search", "fake"})
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`number name installed version usedVersion title
0      fake false             1.0         fake
`))
		})
	})
})
