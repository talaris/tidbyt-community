// Package kickstarter provides details for the Kickstarter applet.
package kickstarter

import (
	_ "embed"

	"tidbyt.dev/community/apps/manifest"
)

//go:embed kickstarter.star
var source []byte

// New creates a new instance of the Kickstarter applet.
func New() manifest.Manifest {
	return manifest.Manifest{
		ID:          "kickstarter",
		Name:        "Kickstarter",
		Author:      "talaris",
		Summary:     "Track Campaign",
		Desc:        "Trakcs the progress of a Kickstarter project.",
		FileName:    "kickstarter.star",
		PackageName: "kickstarter",
		Source:  source,
	}
}
