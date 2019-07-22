package app

var (
	version string
	commit  string
)

func GetVersion() string {
	return version
}

func GetCommit() string {
	return commit
}
