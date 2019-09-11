package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("queue test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		queueClient  QueueClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		queueClient = QueueClient{}
		queueClient.RoundTripper = roundTripper
		queueClient.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("get queue", func() {
		It("should success", func() {
			PrepareGetQueue(roundTripper, queueClient.URL, "", "")

			queue, err := queueClient.Get()
			Expect(err).To(BeNil())
			Expect(queue).NotTo(BeNil())
			Expect(len(queue.Items)).To(Equal(1))
			Expect(queue.Items[0].ID).To(Equal(62))
		})
	})

	Context("cancel", func() {
		It("should success", func() {
			PrepareCancelQueue(roundTripper, queueClient.URL, "", "")

			err := queueClient.Cancel(1)
			Expect(err).To(BeNil())
		})
	})
})
