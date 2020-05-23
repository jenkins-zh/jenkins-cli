package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("center identity command", func() {
	var (
		ctrl           *gomock.Controller
		roundTripper   *mhttp.MockRoundTripper
		targetFilePath string

		err error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		centerIdentityOption.RoundTripper = roundTripper
		targetFilePath = "jenkins.war"

		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		err = os.Remove(targetFilePath)
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should not error", func() {
			var data []byte
			data, err = GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForGetIdentity(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9")

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs([]string{"center", "identity"})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal(`{
 "Fingerprint": "fingerprint",
 "PublicKey": "publicKey",
 "SystemMessage": "systemMessage"
}
`))
		})
	})
})
