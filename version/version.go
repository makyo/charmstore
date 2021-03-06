// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package version

type Version struct {
	GitCommit string
	Version   string
}

var VersionInfo = unknownVersion

var unknownVersion = Version{
	GitCommit: "unknown git commit",
	Version:   "unknown version",
}
