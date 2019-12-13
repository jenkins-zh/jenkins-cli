package cmd

import (
	"bytes"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
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
				PreHooks: []CommndHook{CommndHook{
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
				PreHooks: []CommndHook{CommndHook{
					Path:    "test",
					Command: successCmd,
				}, CommndHook{
					Path:    "test",
					Command: "echo 2",
				}, CommndHook{
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
				PreHooks: []CommndHook{CommndHook{
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
