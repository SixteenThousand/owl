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
)

func TestRestrictRuneset(t *testing.T) {
	removeTCases := map[string]string{
		"Deutschland-\U0001f1e9\U0001f1ea": "Deutschland-",
		"\U0001f1e9\U0001f1ea": "_empty_",
		"The Salmon \U0001f41f Of Doubt": "The Salmon  Of Doubt",
		"Why put questions in file names?": "Why_put_questions_in_file_names",
		"This* causes ~problems~": "This_causes_~problems~",
	}
	representTCases := map[string]string{
		"Deutschland-\U0001f1e9\U0001f1ea": "Deutschland-_U1f1e9__U1f1ea_",
		"\U0001f1e9\U0001f1ea": "_empty_",
		"The Salmon \U0001f41f Of Doubt": "The Salmon _U1f41f_ Of Doubt",
		"Why put questions in file names?": "Why put questions in file names_U3f_",
		"This*-causes-~problems~": "This_U2a_-causes-~problems~",
	}
	for input, want := range removeTCases {
		if have:=restrictRuneset(input, "represent", DEFAULT_RUNEMAP); have != want {
			t.Errorf(
				"With strategy \"remove\":\n\twant <<%s>>\n\thave <<%s>>\n",
				want,
				have,
			)
		}
	}
	for input, want := range representTCases {
		if have:=restrictRuneset(input, "represent", DEFAULT_RUNEMAP); have != want {
			t.Errorf(
				"With strategy \"represent\":\n\twant <<%s>>\n\thave <<%s>>\n",
				want,
				have,
			)
		}
	}
}
