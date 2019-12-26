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
		userClient   UserClient

		username string
		password string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		userClient = UserClient{}
		userClient.RoundTripper = roundTripper
		userClient.URL = "http://localhost"

		username = "admin"
		password = "token"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Get", func() {
		It("should success", func() {
			userClient.UserName = username
			userClient.Token = password

			PrepareGetUser(roundTripper, userClient.URL, username, password)

			user, err := userClient.Get()
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())
			Expect(user.FullName).To(Equal(username))
		})
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

	Context("Delete", func() {
		It("should success", func() {
			userName := "fakeName"
			PrepareForDeleteUser(roundTripper, userClient.URL, userName, "", "")

			err := userClient.Delete(userName)
			Expect(err).To(BeNil())
		})
	})

	Context("Create", func() {
		It("should success", func() {
			targetUserName := "fakeName"
			userClient.UserName = username
			userClient.Token = password

			PrepareCreateUser(roundTripper, userClient.URL, username, password, targetUserName)

			result, err := userClient.Create(targetUserName, "fakePass")
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Username).To(Equal(targetUserName))
			Expect(result.Password1).To(Equal(result.Password2))
			Expect(result.Password1).NotTo(Equal(""))
		})
	})

	Context("CreateToken", func() {
		It("should success, given token name", func() {
			newTokenName := "fakeName"
			userClient.UserName = username
			userClient.Token = password

			PrepareCreateToken(roundTripper, userClient.URL, username, password, newTokenName, username)

			token, err := userClient.CreateToken("", newTokenName)
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeNil())
			Expect(token.Status).To(Equal("ok"))
		})
	})
})
