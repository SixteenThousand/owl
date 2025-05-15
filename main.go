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

// Unrelated to the standard library's "context".
type context = struct {
	FileSpecs []string
	Directories []string
	DryRun bool
	Strategy string
	DoHelp bool
	DoVersion bool
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

/// MAIN FUNCTIONS
func isFatValid(r rune) bool {
	for _, runeRange := range FAT_RUNESET {
		if runeRange[0] <= r && r <= runeRange[1] {
			return true
		}
	}
	return false
}

func restrictRuneset(s, strategy string) string {
	result := s
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
				toValidSubs[r] = fmt.Sprintf("_U%X_", r)
			}
		}
	}
	for old, new := range toValidSubs {
		result = strings.ReplaceAll(result, string(old), new)
	}
	if len(result) == 0 {
		return "_EMPTY_"
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

// TODO: Collect all invalid args into one error instead of failing fast.
func parseCLIArgs(args []string) (context, error) {
	// Set defaults
	result := context{
		FileSpecs: []string{},
		Directories: []string{},
		DryRun: false,
		Strategy: "fat",
		DoHelp: false,
		DoVersion: false,
	}
	index := 1
	isFlag := func(arg string) bool {
		return arg[0] == '-'
	}
	for index < len(args) {
		switch arg := args[index]; arg {
		case "-d", "--directory":
			result.Directories = append(result.Directories, arg)
			index++
		case "-n", "--dry-run":
			result.DryRun = true
		case "-s", "--strategy":
			result.Strategy = args[index]
			index++
		case "-h", "--help":
			result.DoHelp = true
		case "-v", "--version":
			result.DoVersion = true
		default:
			if isFlag(arg) {
				return result, errors.New(fmt.Sprintf(
					"Invalid flag <<%s>>",
					arg,
				))
			} else {
				result.FileSpecs = append(result.FileSpecs, arg)
			}
		}
		index++
	}
	return result, nil
}

// TODO: Finish writing options short help. Mention man page where relevant.
func printHelp() {
	fmt.Print(`Owl - a hunter of bad characters in filenames

 Usage:
   owl [options] FILES

 Rename FILES such that all characters that invalid in FAT file systems 
 (?,\,*,etc.) are removed.
 
 Options:
   -s,--strategy
   -h,--help
   -v,--version
   -d,--directory DIRECTORY

`);
}


func main() {
	ctx, err := parseCLIArgs(os.Args)
	kaput(err)
	if ctx.DoHelp {
		printHelp()
	} else if ctx.DoVersion {
		fmt.Printf(
			"Owl - a hunter of bad characters in file names\nversion %s\n",
			OwlVersion,
		)
	} else {
		for _, file := range ctx.FileSpecs {
			if _, err := os.Stat(file); err != nil {
				warn("File <<%s>> does not exist!", file)
				continue
			}
			oldName := fpath.Base(file)
			dirName := fpath.Dir(file)
			if ctx.DryRun {
				fmt.Printf(
					"%s -> <<%s>>\n",
					file,
					restrictRuneset(oldName, "represent"),
				)
			} else {
				os.Rename(file,fpath.Join(
					dirName,
					restrictRuneset(oldName, "represent"),
				))
			}
		}
	}
}
