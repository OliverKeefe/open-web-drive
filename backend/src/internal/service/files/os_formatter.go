package files

import "time"

// MacOSXMetaData contains metadata for files on Mac OS X.
type MacOSXMetaData struct{}

// WindowsMetaData contains metadata for files on Windows.
type WindowsMetaData struct{}

// GNULinuxMetaData contains metadata for files on Linux.
type GNULinuxMetaData struct {
	File     string
	Size     uint64
	FileType string
	Links    uint64
	Mode     string
	User     string
	Group    string
	Access   time.Time
	Modify   time.Time
	Change   time.Time
	Birth    time.Time
}

// AndroidMetaData contains metadata for files on Android OS.
type AndroidMetaData struct{}

// iOSMetaData contains metadata for files on iOS devices.
type iOSMetaData struct{}

// iPadOSMetaData contains metadata for files on iPad OS devices.
type iPadOSMetaData struct{}
