package cmd

import (
	"encoding/xml"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

//LtsURL is the URL of stable Jenkins RSS
const LtsURL = "https://www.jenkins.io/changelog-stable/rss.xml"

//WeeklyYURL is the URL of stable Jenkins RSS
const WeeklyYURL = "https://www.jenkins.io/changelog/rss.xml"

//WidthOfDescription is the width of the description column
const WidthOfDescription = 60

//NumberOfLinesOfDescription is the number of lines to be printed in description column
const NumberOfLinesOfDescription = 10

//AsciiOfLineFeed is the ASCII of line feed
const AsciiOfLineFeed = 10

//AsciiOfSpace is the ASCII of space
const AsciiOfSpace = 32

//CenterListOption as options for Jenkins RSS
type CenterListOption struct {
	Channel      Channel `xml:"channel"`
	RoundTripper http.RoundTripper
}

//Channel as part of CenterListOption
type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

//Item as a option for information of newly-released Jenkins
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

var centerListOption CenterListOption

func init() {
	centerCmd.AddCommand(centerListCmd)
}

var centerListCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("Print the information of recent-released Jenkins"),
	Long:  i18n.T("Print the information of recent-released Jenkins"),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		jenkins := getCurrentJenkinsFromOptionsOrDie()
		jclient := &client.JenkinsStatusClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerOption.RoundTripper,
			},
		}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth
		status, error := jclient.Get()
		if error != nil {
			return error
		}
		jenkinsVersion := status.Version
		changeLog, err := getChangelog(LtsURL, jenkinsVersion, getVersionData)
		cmd.Println(changeLog)
		//err = printChangelog(WEEKLYURL,cmd)
		return err
	},
}

func getVersionData(rss string) ([]Item, string, error) {
	resp, err := http.Get(rss)
	if err != nil {
		return nil, "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	var centerListOption CenterListOption
	xml.Unmarshal(bytes, &centerListOption)
	resp.Body.Close()
	return centerListOption.Channel.Items, centerListOption.Channel.Title, nil
}

func getChangelog(rss string, version string, getData func(rss string) ([]Item, string, error)) (changelog string, err error) {
	items, title, error := getData(rss)
	if error != nil {
		return "", error
	}
	t := table.NewWriter()
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(table.Row{title, title, title, title}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"Index", "Title", "Description", "PubDate"})
	var temp string
	var isTheLatestVersion = 1
	for index, item := range items {
		if compareVersionTitle(version, item.Title[8:]) >= 0 {
			break
		}
		isTheLatestVersion = 0
		temp = trimXMLSymbols(item.Description)
		temp = regulateWidthAndLines(temp, WidthOfDescription, NumberOfLinesOfDescription)
		t.AppendRow([]interface{}{index + 1, item.Title, temp, item.PubDate[:17]})
		t.AppendSeparator()
	}
	if isTheLatestVersion == 1 {
		return "You already have the latest version of Jenkins installed!", nil
	}
	return t.Render(), nil
}
func regulateWidthAndLines(content string, width int, numberOfLines int) string {
	var count = 0
	myContent := []uint8(content)
	var i int
	for i = 0; i < len(myContent); i++ {
		if myContent[i] == AsciiOfLineFeed {
			count++
		}
		if count == numberOfLines {
			break
		}
	}
	myContent = myContent[:i]
	indexInLine := 0
	for index, char := range myContent {
		indexInLine++
		if indexInLine%width == 0 {
			if char == AsciiOfSpace {
				myContent[index] = AsciiOfLineFeed
				indexInLine = 0
			} else {
				indexInLine = 0
				for ; myContent[index] != AsciiOfSpace; index-- {
					indexInLine++
				}
				myContent[index] = AsciiOfLineFeed
			}
		}
		if char == AsciiOfLineFeed {
			indexInLine = 0
		}
	}
	return string(myContent)
}

func compareVersionTitle(versionOne string, versionTwo string) int {
	versionOneString := strings.Split(versionOne, ".")
	versionTwoString := strings.Split(versionTwo, ".")
	if strings.Compare(versionOne, versionTwo) == 0 {
		return 0
	}
	if versionOneString[0] > versionTwoString[0] {
		return 1
	} else if versionOneString[0] < versionTwoString[0] {
		return -1
	}
	if versionOneString[1] != "" && versionTwoString[1] != "" && versionOneString[1] > versionTwoString[1] {
		return 1
	} else if versionOneString[1] < versionTwoString[1] {
		return -1
	}
	if versionOneString[2] != "" && versionTwoString[2] != "" && versionOneString[2] > versionTwoString[2] {
		return 1
	}
	return -1
}

func trimXMLSymbols(temp string) string {
	temp = strings.TrimSpace(temp)
	xmlSymbols := []string{"<br>", "<br/>", "<br />", "<ul>", "<li>", "</ul>", "<strong>", "</strong>", "<code>", "</code>", "<em>", "</em>"}
	for _, xmlSymbol := range xmlSymbols {
		temp = strings.Replace(temp, xmlSymbol, "", -1)
	}
	temp = strings.Replace(temp, "</li>", "\n", -1)
	return temp
}
