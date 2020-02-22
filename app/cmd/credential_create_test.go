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

var _ = Describe("credential create command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		buf          *bytes.Buffer
		store        string
		id           string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		rootCmd.SetArgs([]string{})
		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		credentialCreateOption.RoundTripper = roundTripper

		store = "system"
		id = "fake-id"
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
			data, err = generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("lack of the necessary parameters", func() {
			rootCmd.SetArgs([]string{"credential", "create", "--store="})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
		})

		It("unknown type", func() {
			rootCmd.SetArgs([]string{"credential", "create", "--type", "fake-type", "--store", store})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
		})

		It("should success with user name and password", func() {
			credential := client.UsernamePasswordCredential{
				Credential: client.Credential{Scope: "GLOBAL"},
			}

			client.PrepareForCreateUsernamePasswordCredential(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", store, credential)

			rootCmd.SetArgs([]string{"credential", "create", "--type", "basic", store, id})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should success with secret", func() {
			credential := client.StringCredentials{
				Credential: client.Credential{Scope: "GLOBAL"},
			}

			client.PrepareForCreateSecretCredential(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", store, credential)

			rootCmd.SetArgs([]string{"credential", "create", "--type", "secret", store, id})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
