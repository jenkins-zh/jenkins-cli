package client

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plugin api test", func() {

	var (
		ctrl      		*gomock.Controller
		pluginApi 		PluginAPI
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("NewPlugins", func() {
		It("New Plugins list", func() {
			response, err := pluginApi.NewPlugins()
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
		})
	})

})

