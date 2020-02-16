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
  insecureSkipVerify: true
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

//func TestConfigGenerate(t *testing.T) {
//	RunEditCommandTest(t, EditCommandTest{
//		ConfirmProcedure: func(c *expect.Console) {
//			c.ExpectString("Cannot found your config file, do you want to edit it?")
//			c.SendLine("y")
//			//c.ExpectEOF()
//		},
//		Procedure: func(c *expect.Console) {
//			c.ExpectString("Edit your config file")
//			c.SendLine("")
//			go c.ExpectEOF()
//			time.Sleep(time.Millisecond)
//			c.Send(`ifake-config`)
//			c.Send("\x1b")
//			c.SendLine(":wq!")
//		},
//		Test: func(stdio terminal.Stdio) (err error) {
//			configFile := path.Join(os.TempDir(), "fake.yaml")
//			defer os.Remove(configFile)
//			configGenerateOption.BatchOption.Stdio = stdio
//			configGenerateOption.CommonOption.Stdio = stdio
//			rootCmd.SetArgs([]string{"config", "generate", "--interactive", "--copy=false", "--configFile=" + configFile})
//			_, err = rootCmd.ExecuteC()
//			return
//		},
//	})
//}
