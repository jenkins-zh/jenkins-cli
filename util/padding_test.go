package util

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Padding util test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("isHan test", func() {
		var han string
		It("english test", func() {
			han = "ZSOMZAYNHG"
			Expect(isHan(han)).To(Equal(false))
		})
		It("chinese test", func() {
			han = "构建一个自由风格的软件项目"
			Expect(isHan(han)).To(Equal(true))
		})
		It("hk test", func() {
			han = "建置 Free-Style 軟體專案"
			Expect(isHan(han)).To(Equal(false))
		})
		It("japanese test", func() {
			han = "フリースタイル・プロジェクトのビルド"
			Expect(isHan(han)).To(Equal(true))
		})
		It("japanese and chinese and english test", func() {
			han = "フリースタイル・プ中ロジasdェクトのビルド"
			Expect(isHan(han)).To(Equal(false))
		})
	})

	Context("countCN test", func() {
		var han string
		It("english test", func() {
			han = "ZSOMZAYNHG"
			Expect(countCN(han)).To(Equal(0))
		})
		It("chinese test", func() {
			han = "构建一个自由风格的软件项目"
			Expect(countCN(han)).To(Equal(13))
		})
		It("hk test", func() {
			han = "建置 Free-Style 軟體專案"
			Expect(countCN(han)).To(Equal(6))
		})
		It("japanese test", func() {
			han = "フリースタイル・プロジェクトのビルド"
			Expect(countCN(han)).To(Equal(18))
		})
		It("japanese and chinese and english test", func() {
			han = "フリースタイル・プ中ロジasdェクトのビルド"
			Expect(countCN(han)).To(Equal(19))
		})
	})

	Context("Lenf test", func() {
		var han string
		It("english test", func() {
			han = "ZSOMZAYNHG"
			Expect(Lenf(han) == len(han)).To(Equal(true))
		})
		It("chinese test", func() {
			han = "构建一个自由风格的软件项目"
			Expect(Lenf(han) == len(han)).To(Equal(false))
		})
		It("hk test", func() {
			han = "建置 Free-Style 軟體專案"
			Expect(Lenf(han) == len(han)).To(Equal(false))
		})
		It("japanese test", func() {
			han = "フリースタイル・プロジェクトのビルド"
			Expect(Lenf(han) == len(han)).To(Equal(false))
		})
		It("japanese and chinese and english test", func() {
			han = "フリースタイル・プ中ロジasdェクトのビルド"
			Expect(Lenf(han) == len(han)).To(Equal(false))
		})
	})
})
