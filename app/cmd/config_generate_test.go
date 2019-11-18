package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config generate command", func() {
	var (
		ctrl     *gomock.Controller
		buf      *bytes.Buffer
		cmdArray []string
		cmdErr   error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		cmdArray = []string{"config", "generate", "-i=false"}

		buf = new(bytes.Buffer)
		rootCmd.SetOutput(buf)
	})

	JustBeforeEach(func() {
		rootCmd.SetArgs(cmdArray)
		_, cmdErr = rootCmd.ExecuteC()
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should success", func() {
			Expect(cmdErr).To(BeNil())
			Expect(buf.String()).To(Equal(`current: yourServer
language: ""
jenkins_servers:
- name: yourServer
  url: http://localhost:8080/jenkins
  username: admin
  token: 111e3a2f0231198855dceaff96f20540a9
  proxy: ""
  proxyAuth: ""
  description: ""
preHooks: []
postHooks: []
pluginSuites: []
mirrors:
- name: default
  url: http://mirrors.jenkins.io/
- name: tsinghua
  url: https://mirrors.tuna.tsinghua.edu.cn/jenkins/
- name: huawei
  url: https://mirrors.huaweicloud.com/jenkins/
- name: tencent
  url: https://mirrors.cloud.tencent.com/jenkins/
# Language context is accept-language for HTTP header, It contains zh-CN/zh-TW/en/en-US/ja and so on
# Goto 'http://localhost:8080/jenkins/me/configure', then you can generate your token.
`))
		})
	})
})
