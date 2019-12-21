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
		name           string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)

		computerClient = client.ComputerClient{}
		computerClient.RoundTripper = roundTripper
		computerClient.URL = "http://localhost"
		name = "fake-name"
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
		client.PrepareForLaunchComputer(roundTripper, computerClient.URL, "", "", name)

		err := computerClient.Launch(name)
		Expect(err).NotTo(HaveOccurred())
	})

	It("GetLog", func() {
		client.PrepareForComputerLogRequest(roundTripper, computerClient.URL, "", "", name)

		log, err := computerClient.GetLog(name)
		Expect(err).NotTo(HaveOccurred())
		Expect(log).To(Equal("fake-log"))
	})

	It("GetLog with 500", func() {
		client.PrepareForComputerLogRequestWithCode(roundTripper, computerClient.URL, "", "", name, 500)

		_, err := computerClient.GetLog(name)
		Expect(err).To(HaveOccurred())
	})

	It("Delete an agent", func() {
		client.PrepareForComputerDeleteRequest(roundTripper, computerClient.URL, "", "", name)

		err := computerClient.Delete(name)
		Expect(err).NotTo(HaveOccurred())
	})

	It("GetSecret of an agent", func() {
		secret := "fake-secret"
		client.PrepareForComputerAgentSecretRequest(roundTripper,
			computerClient.URL, "", "", name, secret)

		result, err := computerClient.GetSecret(name)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(secret))
	})

	It("Create an agent", func() {
		client.PrepareForComputerCreateRequest(roundTripper, computerClient.URL, "", "", name)

		err := computerClient.Create(name)
		Expect(err).NotTo(HaveOccurred())
	})
})
