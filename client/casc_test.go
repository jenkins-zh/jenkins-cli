package client_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		cascManager  client.CASCManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		cascManager = client.CASCManager{}
		cascManager.RoundTripper = roundTripper
		cascManager.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("normal cases", func() {
		client.PrepareForSASCReload(roundTripper, cascManager.URL, "", "")
		client.PrepareForSASCApply(roundTripper, cascManager.URL, "", "")
		client.PrepareForSASCExport(roundTripper, cascManager.URL, "", "")
		client.PrepareForSASCSchema(roundTripper, cascManager.URL, "", "")

		reloadErr := cascManager.Reload()
		applyErr := cascManager.Apply()
		config, exportErr := cascManager.Export()
		schema, schemaErr := cascManager.Schema()

		Expect(reloadErr).NotTo(HaveOccurred())
		Expect(applyErr).NotTo(HaveOccurred())
		Expect(exportErr).NotTo(HaveOccurred())
		Expect(schemaErr).NotTo(HaveOccurred())

		Expect(config).To(Equal("sample"))
		Expect(schema).To(Equal("sample"))
	})

	Context("with error code", func() {
		BeforeEach(func() {
			client.PrepareForSASCExportWithCode(roundTripper, cascManager.URL, "", "", 500)
			client.PrepareForSASCSchemaWithCode(roundTripper, cascManager.URL, "", "", 500)
		})

		It("get error", func() {
			_, exportErr := cascManager.Export()
			_, schemaErr := cascManager.Schema()

			Expect(exportErr).To(HaveOccurred())
			Expect(schemaErr).To(HaveOccurred())
		})
	})
})
