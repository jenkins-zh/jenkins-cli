package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("collect test", func() {
	Context("MaxAndMin", func() {
		It("normal case, should success", func() {
			data := []float64{0.1, 0.2, 0.3, 0.4}
			max, min := MaxAndMin(data)
			Expect(max).To(Equal(0.4))
			Expect(min).To(Equal(0.1))
			return
		})

		It("empty collect, should success", func() {
			data := []float64{}
			max, min := MaxAndMin(data)
			Expect(max).To(Equal(0.0))
			Expect(min).To(Equal(0.0))
			return
		})

		It("only one item, should success", func() {
			data := []float64{0.3}
			max, min := MaxAndMin(data)
			Expect(max).To(Equal(0.3))
			Expect(min).To(Equal(0.3))
			return
		})
	})

	Context("PrintCollectTrend", func() {
		It("should success", func() {
			data := []float64{1512, 3472, 4385, 3981}
			buf := PrintCollectTrend(data)
			Expect(buf).NotTo(Equal(""))
			return
		})
	})
})
