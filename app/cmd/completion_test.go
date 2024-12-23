package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("completion command", func() {
	var (
		ctrl    *gomock.Controller
		cmdArgs []string
		buf     *bytes.Buffer
		err     error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		var data []byte
		data, err = GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	JustBeforeEach(func() {
		rootCmd.SetArgs(cmdArgs)

		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
		_, err = rootCmd.ExecuteC()
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("with default option value", func() {
		BeforeEach(func() {
			cmdArgs = []string{"completion"}
		})

		It("should success", func() {
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("bash completion for jcli"))
		})
	})

	Context("generate zsh completion", func() {
		BeforeEach(func() {
			cmdArgs = []string{"completion", "--type", "zsh"}
		})

		It("should success", func() {
			Expect(err).To(BeNil())
			Expect(buf.String()).NotTo(Equal(""))
		})
	})

	Context("generate powerShell completion", func() {
		BeforeEach(func() {
			cmdArgs = []string{"completion", "--type", "powerShell"}
		})

		It("should success", func() {
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("using namespace System.Management.Automation"))
		})
	})

	Context("generate fish completion", func() {
		BeforeEach(func() {
			cmdArgs = []string{"completion", "--type", "fish"}
		})

		It("should success", func() {
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("fish completion for jcli"))
		})
	})

	Context("generate unknown shell type completion", func() {
		BeforeEach(func() {
			cmdArgs = []string{"completion", "--type", "fake"}
		})

		It("error occurred", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown shell type"))
		})
	})
})
