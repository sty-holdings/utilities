// Package main.go
/*
This will read /tmp/apt-ugrades.tmp and generate apt-get commands

RESTRICTIONS:
    None

NOTES:
    {Enter any additional notes that you believe will help the next developer.}

COPYRIGHT:
	Copyright 2022
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package main

import (
	"log"
	"os"

	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

const (
	APT_UPGRADE_TMP_FILE = "/tmp/apt-upgrades.tmp"
)

// Add types to the types.go file

var (
// Add Variables here for the file (Remember, they are global)
)

func init() {
	// Set up goes here
}

func main() {

	var (
		errorInfo pi.ErrorInfo
	)

	if file, errorInfo = os.Open(APT_UPGRADE_TMP_FILE); log.Fatal(err) {
	}

	// Don't forget to close the file!
	defer file.Close()
	// Afterwards you can perform reading operations...
}
