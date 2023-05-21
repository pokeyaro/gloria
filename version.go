// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

var (
	// Title is the name of the project.
	Title string = "Gloria"

	// Version is the version number of the application.
	Version string = "1.0.0"

	// BuildTime is the build time of the application.
	// Better practice: go build -ldflags "-X main.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')" main.go
	BuildTime string = ""

	// GitHash is the commit hash of the application's Git repository.
	GitHash string = ""
)
