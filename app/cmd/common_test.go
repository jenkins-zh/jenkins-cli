package cmd_test

import (
	"bytes"
	"github.com/Netflix/go-expect"
	"github.com/jenkins-zh/jenkins-cli/app/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	cmd.RunPromptTest(t, cmd.PromptTest{
		Message:    "essage",
		MsgConfirm: &cmd.BatchOption{},
		Procedure: func(c *expect.Console) {
			c.ExpectString("message")
			c.SendLine("y")
			c.ExpectEOF()
		},
		Expected: true,
	})
}
