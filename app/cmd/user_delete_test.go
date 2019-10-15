package cmd

import (
	"bytes"
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("user delete command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
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
		BeforeEach(func() {
			roundTripper = mhttp.NewMockRoundTripper(ctrl)
			userDeleteOption.RoundTripper = roundTripper
		})

		It("lack of arguments", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			userDeleteCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
				cmd.Print("help")
			})

			rootCmd.SetArgs([]string{"user", "delete"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("help"))
		})

		It("should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			targetUserName := "fakename"
			client.PrepareForDeleteUser(roundTripper, "http://localhost:8080/jenkins", targetUserName, "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"user", "delete", targetUserName, "-b", "true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("with status code 500", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			targetUserName := "fakename"
			response := client.PrepareForDeleteUser(roundTripper, "http://localhost:8080/jenkins", targetUserName, "admin", "111e3a2f0231198855dceaff96f20540a9")
			response.StatusCode = 500

			rootCmd.SetArgs([]string{"user", "delete", targetUserName, "-b", "true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("unexpected status code: 500\n"))
		})
	})
})
