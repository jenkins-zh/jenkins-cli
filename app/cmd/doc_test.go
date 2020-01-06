package cmd

import (
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
	})

	AfterEach(func() {
		config = nil
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("should success", func() {
			tmpdir := os.TempDir()
			defer os.RemoveAll(tmpdir)

			rootCmd.SetArgs([]string{"doc", tmpdir})
			_, err := rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(filepath.Join(tmpdir, "jcli_doc.md"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
