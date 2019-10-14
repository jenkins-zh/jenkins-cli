package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("user token command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		userTokenOption.RoundTripper = roundTripper
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

	Context("with http requests", func() {
		It("lack of arguments", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			userTokenCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
				cmd.Print("help")
			})

			rootCmd.SetArgs([]string{"user", "token"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("help"))
		})

		It("should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			tokenName := "fakename"
			client.PrepareCreateToken(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", tokenName)

			rootCmd.SetArgs([]string{"user", "token", "-g", "-n", tokenName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("{\n  \"Status\": \"ok\",\n  \"Data\": {\n    \"TokenName\": \"\",\n    \"TokenUUID\": \"\",\n    \"TokenValue\": \"\"\n  }\n}\n"))
		})

		It("with status code 500", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			tokenName := "fakename"
			response := client.PrepareCreateToken(roundTripper, "http://localhost:8080/jenkins",
				"admin", "111e3a2f0231198855dceaff96f20540a9", tokenName)
			response.StatusCode = 500

			rootCmd.SetArgs([]string{"user", "token", "-g", "-n", tokenName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("unexpected status code: 500\n"))
		})
	})
})
