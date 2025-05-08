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
	"errors"
    "fmt"
    "os"
	fpath "path/filepath"
	"strings"
)


/// TYPES
/**
 * A runeset represents a list of Unicode Code Points (Hereafter referred to
 * as "Code Points") which are considered valid by Owl. Each pair in the
 * list should be two Code Points, which represent the lower & upper bounds
 * (in that order) of a range of Code Points which are valid. Ranges are
 * inclusive.
 */
type runeset = [][2]rune

/**
 * A runemap represents a list of rune swaps that should be preformed before
 * dealing with any invalid runes. The first number in each entry is the
 * Code Point of the rune to be replaced, the second the Code Point of its
 * replacement.
 */
type runemap = [][2]rune

// Unrelated to the standard library's "context".
type context = struct {
	command string
	hasDryRun bool
	fileSpecs []string
}


/// CONSTANTS
// The largest Unicode code point.
// See link below for more details.
// https://www.unicode.org/versions/Unicode16.0.0/core-spec/chapter-2/#G25564
const MAX_CODE_POINT rune = 0x10ffff

var OwlVersion string

// The runeset of valid runes for file names in FAT32 & exFAT file systems.
var FAT_RUNESET = runeset{
	{ 0x20, 0x21 },
	{ 0x23, 0x29 },
	{ 0x2b, 0x2e },
	{ 0x30, 0x39 },
	{ 0x3b, 0x3b },
	{ 0x3d, 0x3d },
	{ 0x40, 0x5b },
	{ 0x5d, 0x7b },
	{ 0x7d, MAX_CODE_POINT},
}

var DEFAULT_RUNEMAP = runemap{
	{ 0x20, 0x5f }, // space -> underscore (a.k.a low line)
}


/// MAIN FUNCTIONS
func isFatValid(r rune) bool {
	for _, runeRange := range FAT_RUNESET {
		if runeRange[0] <= r && r <= runeRange[1] {
			return true
		}
	}
	return false
}

func restrictRuneset(s, strategy string, userSubs runemap) string {
	result := s
	for _, sub := range userSubs {
		result = strings.ReplaceAll(result, string(sub[0]), string(sub[1]))
	}
	toValidSubs := make(map[rune]string)
	if strategy == "remove" {
		for _, r := range result {
			if !isFatValid(r) {
				toValidSubs[r] = ""
			}
		}
	} else {
		for _, r := range result {
			if !isFatValid(r) {
				toValidSubs[r] = fmt.Sprintf("_U%d_", r)
			}
		}
	}
	for old, new := range toValidSubs {
		result = strings.ReplaceAll(result, string(old), new)
	}
	return result
}

func warn(msgFmt, detail string) {
	msg := fmt.Sprintf(msgFmt, detail)
	fmt.Fprintf(os.Stderr, "\x1b[33m%s\x1b[0m\n", msg)
}

func kaput(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", err.Error())
		os.Exit(1)
	}
}

func parseCLIArgs(args []string) (context, error) {
	var result context
	index := 1
	isFlag := func(arg string) bool {
		return arg[0] == '-'
	}
	if !isFlag(args[index]) {
		result.command = args[index]
		index++
	}
	for index < len(args) {
		switch arg := args[index]; arg {
		case "-n", "--dry-run":
			result.hasDryRun = true
		default:
			if isFlag(arg) {
				return result, errors.New(fmt.Sprintf(
					"Invalid flag <<%s>>",
					arg,
				))
			} else {
				result.fileSpecs = append(result.fileSpecs, arg)
			}
		}
		index++
	}
	return result, nil
}

func cleanCommand(ctx context) {
	for _, file := range ctx.fileSpecs {
		if _, err := os.Stat(file); err != nil {
			warn("File <<%s>> does not exist!", file)
			continue
		}
		oldName := fpath.Base(file)
		dirName := fpath.Dir(file)
		if ctx.hasDryRun {
			fmt.Printf(
				"%s -> <<%s>>\n",
				file,
				restrictRuneset(oldName, "represent", DEFAULT_RUNEMAP),
			)
		} else {
			os.Rename(file,fpath.Join(
				dirName,
				restrictRuneset(oldName, "represent", DEFAULT_RUNEMAP),
			))
		}
	}
}

func helpCommand() {
	fmt.Println(`TBD`)
}

func versionCommand() {
	fmt.Printf("Owl file renaming tool, version %s\n", OwlVersion)
}

func main() {
	ctx, err := parseCLIArgs(os.Args)
	kaput(err)
	switch ctx.command {
	case "clean":
		cleanCommand(ctx)
	case "help":
		helpCommand()
	case "version":
		versionCommand()
	default:
		helpCommand()
	}
}
