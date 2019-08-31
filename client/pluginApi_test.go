package client

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginApi test", func() {
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
		jclient := PluginAPI{}
		It("test return count", func() {
			Expect(jclient.SearchPlugins("ccm").Total >= 1).To(BeTrue())

			Expect(jclient.SearchPlugins("pipeline").Total >= 100).To(BeTrue())
		})
	})
})
