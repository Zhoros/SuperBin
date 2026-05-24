/*
This file is part of GigaPaste.

GigaPaste is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

GigaPaste is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with GigaPaste. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Setting struct {
	FileSizeLimit                  int64  `json:"FileSizeLimitMB"`
	TextSizeLimit                  int64  `json:"TextSizeLimitMB"`
	StreamSizeLimit                int64  `json:"StreamSizeLimitKB"`
	StreamThrottle                 int64  `json:"StreamThrottleMS"`
	Pbkdf2Iteraions                int    `json:"Pbkdf2Iteraions"`
	CmdUploadDefaultDurationMinute int64  `json:"CmdUploadDefaultDurationMinute"`
	EnablePassword                 bool   `json:"enablePassword"`
	Password                       string `json:"password"`
}

var Global Setting

func InitSettings() {

	file, err := os.Open("./data/settings.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Global)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Reject weak or default password when password protection is enabled
	if Global.EnablePassword {
		if Global.Password == "" {
			fmt.Fprintf(os.Stderr, "fatal: enablePassword is true but password is empty. Set a strong password in data/settings.json\n")
			os.Exit(1)
		}
		if Global.Password == "password" {
			fmt.Fprintf(os.Stderr, "fatal: password is set to the default value 'password'. Please change it to a strong password in data/settings.json\n")
			os.Exit(1)
		}
		if len(Global.Password) < 8 {
			fmt.Fprintf(os.Stderr, "fatal: password is too short (%d characters). Use at least 8 characters.\n", len(Global.Password))
			os.Exit(1)
		}
	}

}
