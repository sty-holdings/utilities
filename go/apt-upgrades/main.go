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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

//goland:noinspection ALL
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
		errorInfo       pi.ErrorInfo
		tAptUpgradefile *os.File
		tNodes          []string
	)

	if tAptUpgradefile, errorInfo.Error = os.Open(APT_UPGRADE_TMP_FILE); errorInfo.Error != nil {
		log.Fatal(errorInfo.Error)
	}
	defer tAptUpgradefile.Close()

	scanner := bufio.NewScanner(tAptUpgradefile)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), ctv.FORWARD_SLASH) {
			tNodes = strings.Split(scanner.Text(), ctv.FORWARD_SLASH)
			fmt.Printf("sudo apt-get upgrade %v -y \n", tNodes[0])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}
