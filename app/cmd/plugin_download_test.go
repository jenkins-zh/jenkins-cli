package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("plugin download command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginDownloadOption.RoundTripper = roundTripper
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
		It("should success", func() {
			var err error
			var data []byte
			data, err = generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareOnePluginWithOptionalDep(roundTripper, "fake")
			client.PrepareDownloadPlugin(roundTripper)

			rootCmd.SetArgs([]string{"plugin", "download", "fake", "--show-progress=false", "--use-mirror=false"})

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(""))

			_, err = os.Stat("fake.hpi")
			defer os.Remove("fake.hpi")
			Expect(err).To(BeNil())
		})
	})
})
