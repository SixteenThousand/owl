/**
 * Owl - file renaming tool
 * Copyright (C) 2025 User SixteenThousand of github.com
 * Email: thomsixteenthousand@gmail.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"testing"
	"strings"
)

func sliceToString(slice []string) string {
	var b strings.Builder
	b.WriteString("[\n")
	for _, s := range slice {
		b.WriteString("  ")
		b.WriteString(s)
		b.WriteRune('\n')
	}
	b.WriteRune(']')
	return b.String()
}

func TestRestrictRuneset(t *testing.T) {
	removeTCases := map[string]string{
		"::?\\": "_EMPTY_",
		"The Salmon \U0001f41f Of Doubt": "The Salmon \U0001f41f Of Doubt",
		"Why put questions in file names?": "Why put questions in file names",
		"This* causes ~problems~": "This causes ~problems~",
		"There are* *alot ** of c*veats * here": "There are alot  of cveats  here",
	}
	representTCases := map[string]string{
		"::?\\": "_U3A__U3A__U3F__U5C_",
		"Why put questions in file names?": "Why put questions in file names_U3F_",
		"The Salmon \U0001f41f Of Doubt": "The Salmon \U0001f41f Of Doubt",
		"This*-causes-~problems~": "This_U2A_-causes-~problems~",
	}
	for input, want := range removeTCases {
		have := restrictRuneset(input, "remove")
		if ; have != want {
			t.Errorf(
				"With strategy \"remove\":\n\twant <<%s>>\n\thave <<%s>>\n",
				want,
				have,
			)
		}
	}
	for input, want := range representTCases {
		if have:=restrictRuneset(input, "represent"); have != want {
			t.Errorf(
				"With strategy \"represent\":\n\twant <<%s>>\n\thave <<%s>>\n",
				want,
				have,
			)
		}
	}
}
