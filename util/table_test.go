package util

import (
	"os"

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
		It("test", func() {
			table := CreateTable(os.Stdout)
			table.AddRow("number", "name", "type")
			table.AddRow("0 ZSOMZAYNHG Standalone Projects")
			table.AddRow("0 构建一个自由风格的软件项 Standalone Projects")
			table.AddRow("0 建置 Free-Style 軟體專案 Standalone Projects")
			table.AddRow("0 フリースタイル・プロジェクトのビルド Standalone Projects")
			table.AddRow("0 フリースタイル・プ中ロジasdェクトのビル Standalone Projects")
			table.Render()
			table1 := CreateTable(os.Stdout)
			table1.AddRow(`
			number name                                type
			0      ZSOMZAYNHG                          Standalone Projects
			1      构建一个自由风格的软件项                Standalone Projects
			2      建置 Free-Style 軟體專案              Standalone Projects
			3      フリースタイル・プロジェクトのビルド      Standalone Projects
			4      フリースタイル・プ中ロジasdェクトのビル   Standalone Projects
			`)
			Expect(table.Out).To(Equal(table1.Out))
		})
	})
})
