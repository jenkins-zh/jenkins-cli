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
		ctrl    *gomock.Controller
		rootCmd *cobra.Command
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		rootCmd = &cobra.Command{Use: "root"}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("PreHook test", func() {
		It("only with root cmd", func() {
			path := getCmdPath(rootCmd)
			Expect(path).To(Equal(""))
		})

		It("one sub cmd", func() {
			subCmd := &cobra.Command{
				Use: "sub",
			}
			rootCmd.AddCommand(subCmd)

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
			rootCmd.AddCommand(sub1Cmd)
			sub1Cmd.AddCommand(sub2Cmd)

			path := getCmdPath(sub2Cmd)
			Expect(path).To(Equal("sub1.sub2"))
		})
	})

	Context("execute cmd test", func() {
		It("basic command", func() {
			var buf bytes.Buffer
			err := execute("echo 1", &buf)

			Expect(buf.String()).To(Equal("1\n"))
			Expect(err).To(BeNil())
		})

		It("error command", func() {
			var buf bytes.Buffer
			err := execute("exit 1", &buf)

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
				PreHooks: []PreHook{PreHook{
					Path:    "test",
					Command: "echo 1",
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
				PreHooks: []PreHook{PreHook{
					Path:    "test",
					Command: "echo 1",
				}, PreHook{
					Path:    "test",
					Command: "echo 2",
				}, PreHook{
					Path:    "fake",
					Command: "echo 1",
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
				PreHooks: []PreHook{PreHook{
					Path:    "test",
					Command: "exit 1",
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
