module github.com/jenkins-zh/jenkins-cli

go 1.12

replace github.com/AlecAivazis/survey v1.8.5 => gopkg.in/AlecAivazis/survey.v1 v1.8.5

require (
	github.com/AlecAivazis/survey/v2 v2.0.1
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/atotto/clipboard v0.1.2
	github.com/golang/mock v1.3.1
	github.com/gosuri/uiprogress v0.0.1
	github.com/linuxsuren/jenkins-cli v0.0.17
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.6.0
	github.com/spf13/cobra v0.0.5
	gopkg.in/yaml.v2 v2.2.2
)
