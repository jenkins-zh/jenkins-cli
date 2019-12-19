package client_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("computer test", func() {
	var (
		ctrl           *gomock.Controller
		computerClient client.ComputerClient
		roundTripper   *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)

		computerClient = client.ComputerClient{}
		computerClient.RoundTripper = roundTripper
		computerClient.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("List", func() {
		client.PrepareForComputerListRequest(roundTripper, computerClient.URL, "", "")

		computers, err := computerClient.List()
		Expect(err).NotTo(HaveOccurred())
		Expect(computers).NotTo(BeNil())
		Expect(len(computers.Computer)).To(Equal(2))
	})

	It("Launch", func() {
		name := "fake-name"
		client.PrepareForLaunchComputer(roundTripper, computerClient.URL, "", "", name)

		err := computerClient.Launch(name)
		Expect(err).NotTo(HaveOccurred())
	})

	It("GetLog", func() {
		name := "fake-name"
		client.PrepareForComputerLogRequest(roundTripper, computerClient.URL, "", "", name)

		log, err := computerClient.GetLog(name)
		Expect(err).NotTo(HaveOccurred())
		Expect(log).To(Equal("fake-log"))
	})

	It("GetLog with 500", func() {
		name := "fake-name"
		client.PrepareForComputerLogRequestWithCode(roundTripper, computerClient.URL, "", "", name, 500)

		_, err := computerClient.GetLog(name)
		Expect(err).To(HaveOccurred())
	})
})
