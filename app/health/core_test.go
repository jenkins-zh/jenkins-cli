package health

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Opt only for test
type Opt struct{}

// Check only for test
func (o *Opt) Check() error {
	return nil
}

var _ = Describe("test command health check interface", func() {
	var (
		register CheckRegister
	)

	BeforeEach(func() {
		register = CheckRegister{
			Member: make(map[string]CommandHealth, 0),
		}
	})

	It("basic test", func() {
		Expect(register.Member).NotTo(BeNil())
		Expect(len(register.Member)).To(Equal(0))
	})

	Context("register a fake one", func() {
		It("should success", func() {
			register.Register("fake", &Opt{})
			Expect(len(register.Member)).To(Equal(1))
		})
	})
})
