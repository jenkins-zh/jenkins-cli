package cmd

import (
	"encoding/xml"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"os"
)

var _ = Describe("center list command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		centerListOption.RoundTripper = roundTripper
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
		It("no jenkins version information in the list", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			resp, _ := http.Get(LTSURL)
			bytes, _ := ioutil.ReadAll(resp.Body)
			var centerListOption CenterListOption
			xml.Unmarshal(bytes, &centerListOption)
			theNewestJenkinsVersion := centerListOption.Channel.Items[0].Title[8:]
			temp, _ := printChangelog(LTSURL, theNewestJenkinsVersion)

			Expect(temp == "You already have the latest version of Jenkins installed!").To(BeTrue())

		})
	})
})
