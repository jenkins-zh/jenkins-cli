package cmd

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"bytes"
)

var _ = Describe("doc command test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		config = nil
	})

	AfterEach(func() {
		config = nil
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("lack of arguments", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			docCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
				cmd.Print("help")
			})

			rootCmd.SetArgs([]string{"doc"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("help"))
		})

		It("should success", func(){
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			rootCmd.SetArgs([]string{"doc", "/tmp"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(""))
		})
	})
})
