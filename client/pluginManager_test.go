package client

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginManager test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("basic function test", func() {
		It("get install plugin query string", func() {
			names := make([]string, 0)
			Expect(getPluginsInstallQuery(names)).To(Equal(""))

			names = append(names, "abc")
			Expect(getPluginsInstallQuery(names)).To(Equal("plugin.abc="))

			names = append(names, "def")
			Expect(getPluginsInstallQuery(names)).To(Equal("plugin.abc=&plugin.def="))

			names = append(names, "")
			Expect(getPluginsInstallQuery(names)).To(Equal("plugin.abc=&plugin.def="))
		})
	})
})
