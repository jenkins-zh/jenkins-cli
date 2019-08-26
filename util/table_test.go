package util

import (
	"bytes"

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
		It("english test", func() {
			var buffer bytes.Buffer
			table := CreateTable(&buffer)
			table.AddRow("number", "name", "type")
			table.AddRow("0", "Freestyle project", "Standalone Projects")
			table.AddRow("1", "Maven project", "Standalone Projects")
			table.AddRow("2", "Pipeline", "Standalone Projects")
			table.AddRow("3", "External Job", "Standalone Projects")
			table.AddRow("4", "Multi-configuration project", "Standalone Projects")
			table.AddRow("0", "Bitbucket Team/Project", "Nested Projects")
			table.AddRow("1", "Folder", "Nested Projects")
			table.AddRow("2", "GitHub Organization", "Nested Projects")
			table.AddRow("3", "Multibranch Pipeline", "Nested Projects")
			table.Render()
			comp := string(`number name                        type
0      Freestyle project           Standalone Projects
1      Maven project               Standalone Projects
2      Pipeline                    Standalone Projects
3      External Job                Standalone Projects
4      Multi-configuration project Standalone Projects
0      Bitbucket Team/Project      Nested Projects
1      Folder                      Nested Projects
2      GitHub Organization         Nested Projects
3      Multibranch Pipeline        Nested Projects`)
			comp = comp + "\n"
			Expect(buffer.String()).To(Equal(comp))
		})

		It("chinese test", func() {
			var buffer bytes.Buffer
			table := CreateTable(&buffer)
			table.AddRow("number", "name", "type")
			table.AddRow("0", "构建一个自由风格的软件项目", "Standalone Projects")
			table.AddRow("1", "构建一个maven项目", "Standalone Projects")
			table.AddRow("2", "流水线", "Standalone Projects")
			table.AddRow("3", "External Job", "Standalone Projects")
			table.AddRow("4", "构建一个多配置项目", "Standalone Projects")
			table.AddRow("0", "Bitbucket Team/Project", "Nested Projects")
			table.AddRow("1", "文件夹", "Nested Projects")
			table.AddRow("2", "GitHub 组织", "Nested Projects")
			table.AddRow("3", "多分支流水线", "Nested Projects")
			table.Render()
			comp := string(`number name                       type
0      构建一个自由风格的软件项目 Standalone Projects
1      构建一个maven项目          Standalone Projects
2      流水线                     Standalone Projects
3      External Job               Standalone Projects
4      构建一个多配置项目         Standalone Projects
0      Bitbucket Team/Project     Nested Projects
1      文件夹                     Nested Projects
2      GitHub 组织                Nested Projects
3      多分支流水线               Nested Projects`)
			comp = comp + "\n"
			Expect(buffer.String()).To(Equal(comp))
		})

		It("japanese test", func() {
			var buffer bytes.Buffer
			table := CreateTable(&buffer)
			table.AddRow("number", "name", "type")
			table.AddRow("0", "フリースタイル・プロジェクトのビルド", "Standalone Projects")
			table.AddRow("1", "Mavenプロジェクトのビルド", "Standalone Projects")
			table.AddRow("2", "パイプライン", "Standalone Projects")
			table.AddRow("3", "外部ジョブ", "Standalone Projects")
			table.AddRow("4", "マルチ構成プロジェクトのビルド", "Standalone Projects")
			table.AddRow("0", "Bitbucket Team/Project", "Nested Projects")
			table.AddRow("1", "フォルダ", "Nested Projects")
			table.AddRow("2", "GitHub Organization", "Nested Projects")
			table.AddRow("3", "Multibranch Pipeline", "Nested Projects")
			table.Render()
			comp := string(`number name                                 type
0      フリースタイル・プロジェクトのビルド Standalone Projects
1      Mavenプロジェクトのビルド            Standalone Projects
2      パイプライン                         Standalone Projects
3      外部ジョブ                           Standalone Projects
4      マルチ構成プロジェクトのビルド       Standalone Projects
0      Bitbucket Team/Project               Nested Projects
1      フォルダ                             Nested Projects
2      GitHub Organization                  Nested Projects
3      Multibranch Pipeline                 Nested Projects`)
			comp = comp + "\n"
			Expect(buffer.String()).To(Equal(comp))
		})
	})
})
