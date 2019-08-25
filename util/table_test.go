package util

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Table util test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("basic test", func() {
		var han string
		It("english test", func() {
			han = "ZSOMZAYNHG"
			un := toHz(han)
			Expect(un).To(Equal("0 ZSOMZAYNHG Standalone Projects"))
		})
		It("chinese test", func() {
			han = "构建一个自由风格的软件项目"
			un := toHz(han)
			Expect(un).To(Equal("0 构建一个自由风格的软件项目 Standalone Projects"))
		})
		It("hk test", func() {
			han = "建置 Free-Style 軟體專案"
			un := toHz(han)
			Expect(un).To(Equal("0 建置 Free-Style 軟體專案 Standalone Projects"))
		})
		It("japanese test", func() {
			han = "フリースタイル・プロジェクトのビルド"
			un := toHz(han)
			Expect(un).To(Equal("0 フリースタイル・プロジェクトのビルド Standalone Projects"))
		})
		It("japanese and chinese and english test", func() {
			han = "フリースタイル・プ中ロジasdェクトのビルド"
			un := toHz(han)
			Expect(un).To(Equal("0 フリースタイル・プ中ロジasdェクトのビルド Standalone Projects"))
		})
	})
})
