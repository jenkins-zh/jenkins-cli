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
		//request, _ := http.NewRequest("GET", fmt.Sprintf("%s/computer/api/json", computerClient.URL), nil)
		//response := &http.Response{
		//	StatusCode: 200,
		//	Request:    request,
		//	Body:       ioutil.NopCloser(bytes.NewBufferString()),
		//}
		//roundTripper.EXPECT().
		//	RoundTrip(request).Return(response, nil)

		computers, err := computerClient.List()
		Expect(err).NotTo(HaveOccurred())
		Expect(computers).NotTo(BeNil())
		Expect(len(computers.Computer)).To(Equal(1))
	})
})
