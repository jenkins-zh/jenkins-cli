package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/jenkins-zh/jenkins-cli/app/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

var _ = Describe("test OutputOption", func() {
	var (
		outputOption cmd.OutputOption
		fakeFoos     []FakeFoo
	)

	BeforeEach(func() {
		outputOption = cmd.OutputOption{}

		fakeFoos = []FakeFoo{{
			Name: "fake",
		}, {
			Name: "foo-1",
		}}
	})

	Context("ListFilter test", func() {
		var (
			result interface{}
		)

		JustBeforeEach(func() {
			result = outputOption.ListFilter(fakeFoos)
		})

		It("without filter", func() {
			Expect(result).To(Equal(fakeFoos))
		})

		Context("with filter", func() {
			BeforeEach(func() {
				outputOption = cmd.OutputOption{
					Filter: []string{"Name=fake"},
				}
			})

			It("should success", func() {
				Expect(result).NotTo(Equal(fakeFoos))
			})
		})
	})

	Context("OutputV2 test", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			err = outputOption.OutputV2(fakeFoos)
		})

		It("without io writer", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("writer"))
		})

		Context("without columns", func() {
			BeforeEach(func() {
				outputOption.Writer = new(bytes.Buffer)
			})

			It("get no columns error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("columns"))
			})
		})

		Context("with io writer", func() {
			var (
				buf *bytes.Buffer
			)

			BeforeEach(func() {
				buf = new(bytes.Buffer)
				outputOption.Writer = buf
				outputOption.Columns = "Name"
			})

			It("with the default output format", func() {
				Expect(buf.String()).To(Equal(`Name
fake
foo-1
`))
			})

			Context("with json format", func() {
				BeforeEach(func() {
					outputOption.Format = cmd.JSONOutputFormat
				})

				It("should get a json text", func() {
					Expect(buf.String()).To(Equal(`[
  {
    "Name": "fake"
  },
  {
    "Name": "foo-1"
  }
]`))
				})
			})

			Context("with yaml format", func() {
				BeforeEach(func() {
					outputOption.Format = cmd.YAMLOutputFormat
				})

				It("should get a yaml text", func() {
					Expect(buf.String()).To(Equal(`- name: fake
- name: foo-1
`))
				})
			})

			Context("with a unknown format", func() {
				BeforeEach(func() {
					outputOption.Format = "fake"
				})

				It("should get an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("not support format"))
				})
			})
		})
	})

	Context("Match test", func() {
		var (
			result bool
		)

		JustBeforeEach(func() {
			result = outputOption.Match(reflect.ValueOf(fakeFoos[0]))
		})

		It("without filter", func() {
			Expect(result).To(BeTrue())
		})

		Context("ignore invalid filter", func() {
			BeforeEach(func() {
				outputOption = cmd.OutputOption{
					Filter: []string{"Name"},
				}
			})

			It("not matched", func() {
				Expect(result).To(BeTrue())
			})
		})
	})
})

// FakeFoo only for test
type FakeFoo struct {
	Name string
}

func TestHelloTest(t *testing.T) {
	RunPromptTest(t, PromptTest{
		Message:    "essage",
		MsgConfirm: &cmd.BatchOption{},
		procedure: func(c *expect.Console) {
			c.ExpectString("message")
			c.SendLine("y")
			c.ExpectEOF()
		},
		expected: true,
	})
}

type PromptTest struct {
	Message string
	//prompt     survey.Prompt
	MsgConfirm cmd.MsgConfirm
	procedure  func(*expect.Console)
	expected   interface{}
}

func RunPromptTest(t *testing.T, test PromptTest) {
	var answer interface{}
	RunTest(t, test.procedure, func(stdio terminal.Stdio) error {
		batch := &cmd.BatchOption{
			Batch: false,
			Stdio: stdio,
		}
		answer = batch.Confirm(test.Message)
		return nil
	})
	require.Equal(t, test.expected, answer)
}

func Stdio(c *expect.Console) terminal.Stdio {
	return terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
}

func RunTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio) error) {
	t.Parallel()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)

	//c, err := expect.NewConsole(expect.WithStdout(buf))
	//c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	c, _, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))

	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(c)
	}()

	err = test(Stdio(c))
	fmt.Println("Raw output: ", buf.String())
	require.Nil(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	c.Tty().Close()
	<-donec

	// Dump the terminal's screen.
	//fmt.Sprintf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
