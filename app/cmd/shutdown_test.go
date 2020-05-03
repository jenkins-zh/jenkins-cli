package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/app/cmd/common"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("shutdown command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		configFile   string
		err          error
	)

	BeforeEach(func() {
		rootOptions := GetRootOptions()
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootOptions.CommonOption = &common.CommonOption{
			RoundTripper: roundTripper,
		}
		rootOptions.Jenkins = ""
		configFile = path.Join(os.TempDir(), "fake.yaml")
		rootOptions.ConfigFile = configFile

		var data []byte
		data, err = GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.Remove(configFile)
		GetRootOptions().ConfigFile = ""
		ctrl.Finish()
	})

	Context("with batch mode", func() {
		It("should success", func() {
			client.PrepareForShutdown(roundTripper, "http://localhost:8080/jenkins", "admin", "111e3a2f0231198855dceaff96f20540a9", true)

			GetRootCommand().SetArgs([]string{"shutdown", "-b"})

			buf := new(bytes.Buffer)
			GetRootCommand().SetOutput(buf)
			_, err = GetRootCommand().ExecuteC()
			Expect(err).To(BeNil())
		})
	})
})
