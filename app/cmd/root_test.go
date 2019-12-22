package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var _ = Describe("Root cmd test", func() {
	var (
		ctrl        *gomock.Controller
		fakeRootCmd *cobra.Command
		successCmd  string
		errorCmd    string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		fakeRootCmd = &cobra.Command{Use: "root"}
		successCmd = "echo 1"
		errorCmd = "exit 1"
		config = nil
	})

	AfterEach(func() {
		config = nil
		ctrl.Finish()
	})

	Context("invalid logger level", func() {
		It("cause errors", func() {
			rootCmd.SetArgs([]string{"--logger-level", "fake"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())

			rootCmd.SetArgs([]string{"--logger-level", "warn"})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("PreHook test", func() {
		It("only with root cmd", func() {
			path := getCmdPath(fakeRootCmd)
			Expect(path).To(Equal(""))
		})

		It("one sub cmd", func() {
			subCmd := &cobra.Command{
				Use: "sub",
			}
			fakeRootCmd.AddCommand(subCmd)

			path := getCmdPath(subCmd)
			Expect(path).To(Equal("sub"))
		})

		It("two sub cmds", func() {
			sub1Cmd := &cobra.Command{
				Use: "sub1",
			}
			sub2Cmd := &cobra.Command{
				Use: "sub2",
			}
			fakeRootCmd.AddCommand(sub1Cmd)
			sub1Cmd.AddCommand(sub2Cmd)

			path := getCmdPath(sub2Cmd)
			Expect(path).To(Equal("sub1.sub2"))
		})
	})

	Context("execute cmd test", func() {
		It("basic command", func() {
			var buf bytes.Buffer
			err := execute(successCmd, &buf)

			Expect(buf.String()).To(Equal("1\n"))
			Expect(err).To(BeNil())
		})

		It("error command", func() {
			var buf bytes.Buffer
			err := execute(errorCmd, &buf)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("execute pre cmd", func() {
		It("should error", func() {
			err := executePreCmd(nil, nil, nil)
			Expect(err).To(HaveOccurred())
		})

		It("basic use case with one preHook, should success", func() {
			config = &Config{
				PreHooks: []CommandHook{CommandHook{
					Path:    "test",
					Command: successCmd,
				}},
			}

			rootCmd := &cobra.Command{}
			subCmd := &cobra.Command{
				Use: "test",
			}
			rootCmd.AddCommand(subCmd)

			var buf bytes.Buffer
			err := executePreCmd(subCmd, nil, &buf)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("1\n"))
		})

		It("basic use case with many preHooks, should success", func() {
			config = &Config{
				PreHooks: []CommandHook{CommandHook{
					Path:    "test",
					Command: successCmd,
				}, CommandHook{
					Path:    "test",
					Command: "echo 2",
				}, CommandHook{
					Path:    "fake",
					Command: successCmd,
				}},
			}

			rootCmd := &cobra.Command{}
			subCmd := &cobra.Command{
				Use: "test",
			}
			rootCmd.AddCommand(subCmd)

			var buf bytes.Buffer
			err := executePreCmd(subCmd, nil, &buf)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("1\n2\n"))
		})

		It("basic use case without preHooks, should success", func() {
			config = &Config{}

			rootCmd := &cobra.Command{}
			subCmd := &cobra.Command{
				Use: "test",
			}
			rootCmd.AddCommand(subCmd)

			var buf bytes.Buffer
			err := executePreCmd(subCmd, nil, &buf)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal(""))
		})

		It("basic use case with error command, should success", func() {
			config = &Config{
				PreHooks: []CommandHook{CommandHook{
					Path:    "test",
					Command: errorCmd,
				}},
			}

			rootCmd := &cobra.Command{}
			subCmd := &cobra.Command{
				Use: "test",
			}
			rootCmd.AddCommand(subCmd)

			var buf bytes.Buffer
			err := executePreCmd(subCmd, nil, &buf)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("basic root command test", func() {
		var (
			buf *bytes.Buffer
		)

		BeforeEach(func() {
			rootOptions = RootOptions{}
			buf = new(bytes.Buffer)
			rootCmd.SetOut(buf)
		})

		It("should contain substring Version:", func() {
			rootCmd.SetArgs([]string{"--version"})
			_, err := rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("Version:"))
		})

		It("with a fake jenkins as option", func() {
			rootCmd.SetArgs([]string{"--jenkins", "fake"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("cannot found the configuration:"))
		})

		It("with an exists jenkins as option", func() {
			configFile, err := ioutil.TempFile("/tmp", ".yaml")
			Expect(err).NotTo(HaveOccurred())

			defer os.Remove(configFile.Name())

			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(configFile.Name(), data, 0664)
			Expect(err).To(BeNil())

			rootCmd.SetArgs([]string{"--jenkins", "yourServer", "--configFile", configFile.Name()})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("Current Jenkins is:"))
		})
	})

	Context("use connection from options", func() {
		var (
			err error
		)

		BeforeEach(func() {
			rootOptions = RootOptions{}
		})

		AfterEach(func() {
			rootOptions = RootOptions{}
		})

		It("fake jenkins, but with URL from option", func() {
			rootCmd.SetArgs([]string{"--configFile", "fake", "--url", "fake-url",
				"--username", "fake-user", "--token", "fake-token"})
			_, err = rootCmd.ExecuteC()
			Expect(err).NotTo(HaveOccurred())

			jenkins := getCurrentJenkinsFromOptions()
			Expect(jenkins.URL).To(Equal("fake-url"))
			Expect(jenkins.UserName).To(Equal("fake-user"))
			Expect(jenkins.Token).To(Equal("fake-token"))
		})
	})
})

// FakeOpt only for test
type FakeOpt struct{}

// Check fake, only for test
func (o *FakeOpt) Check() error {
	return nil
}

var _ = Describe("RunDiagnose test", func() {
	It("should success", func() {
		opt := RootOptions{Doctor: true}

		rootCmd := &cobra.Command{
			Use: "fake",
		}
		healthCheckRegister.Register(getCmdPath(rootCmd), &FakeOpt{})
		err := opt.RunDiagnose(rootCmd)
		Expect(err).NotTo(HaveOccurred())
	})
})
