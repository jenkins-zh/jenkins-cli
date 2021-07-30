package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("doc command test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		config = nil

		rootOptions.ConfigFile = "test.yaml"

		data, err := GenerateSampleConfig()
		Expect(err).To(BeNil())
		err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		config = nil
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("should success", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)

			tmpdir := os.TempDir()
			defer os.RemoveAll(tmpdir)

			rootCmd.SetArgs([]string{"doc", tmpdir})
			_, err := rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(Equal(""))

			_, err = os.Stat(filepath.Join(tmpdir, "jcli_doc.md"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
