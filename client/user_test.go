package client

import (
	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user test", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		userClient  UserClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		userClient = UserClient{}
		userClient.RoundTripper = roundTripper
		userClient.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("EditDesc", func() {
		It("should success", func() {
			userName := "fakeName"
			description := "fakeDesc"
			PrepareForEditUserDesc(roundTripper, userClient.URL, userName, description, "", "")

			userClient.UserName = userName
			err := userClient.EditDesc(description)
			Expect(err).To(BeNil())
		})
	})
})
