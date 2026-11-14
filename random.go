/*
This file is part of GigaPaste.

GigaPaste is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

GigaPaste is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with GigaPaste. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"crypto/rand"
	"database/sql"
	"os"
	"strconv"
	"time"
)

// we want to avoid ambiguous characters like i, I, l, 1, etc
const charset = "abcdefghkmnpqrstwxyzABCDEFGHJKLMNPQRSTWXYZ23456789"

func GenRandFileName(basePath string, extension string) string {

	for {
		fileName := strconv.FormatInt(time.Now().UnixMilli(), 10) + genRandString(5) + extension
		filePath := basePath + fileName

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath
		}
	}
}

func GenRandPath(length int, db *sql.DB) string {

	for {

		randPath := genRandString(length)
		var id string

		db.QueryRow("SELECT id FROM data WHERE id = ?", randPath).Scan(&id)

		if id == "" {
			return randPath
		}

	}

}

func genRandString(length int) string {

	// Use a CSPRNG: these strings are unguessable access tokens for stored
	// content. Reject the biased tail so the mapping onto charset stays uniform
	// (len(charset) does not divide 256 evenly).
	limit := byte(256 - (256 % len(charset)))

	result := make([]byte, length)
	buf := make([]byte, 1)
	for i := range result {
		for {
			if _, err := rand.Read(buf); err != nil {
				continue
			}
			if buf[0] < limit {
				result[i] = charset[int(buf[0])%len(charset)]
				break
			}
		}
	}
	return string(result)
}
