package app

// Version represents the Jenkins CLI build version.
type Version struct {
	// Major and minor version.
	Number float32

	// Increment this for bug releases
	PatchLevel int

	// JCLI Suffix is the suffix used in the Jenkins CLI version string.
	// It will be blank for release versions.
	Suffix string
}
