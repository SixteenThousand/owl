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
	"io/fs"
    "os"
	fpath "path/filepath"
	"slices"
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
type runeset [][2]rune

// Unrelated to the standard library's "context".
type context struct {
	FileList []string
	RecurseDirs []string
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
var (
	ErrNoPwd = errors.New("Could not get current working directory!")
)

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
/**
 * Compares paths so that all files come before the directories that contain 
 * them. Files in the same directory come in alphanumeric
 * order.
 * Note that paths a,b MUST be absolute.
 */
func comparePaths(a, b string) int {
	sep := string(fpath.Separator)
	aNumComponents := len(strings.Split(a, sep))
	bNumComponents := len(strings.Split(b, sep))
	if aNumComponents == bNumComponents {
	    return strings.Compare(fpath.Base(a), fpath.Base(b))
	}
	return bNumComponents - aNumComponents
}

func (ctx *context) parseFileList() ([]string, error) {
	errMsgs := []string{}
	result := slices.Clone(ctx.FileList)
	slices.SortFunc(result, comparePaths)
	for _, path := range result {
		if _, err := os.Stat(path); err != nil {
			msg := fmt.Sprintf( "File <<%s>> does not exist", path)
			errMsgs = append(errMsgs, msg)
		}
	}
	addFiles := func(fullPath string, fileInfo fs.DirEntry, err error) error {
		if err != nil {
			msg := fmt.Sprintf(
				"Directory <<%s>> does not exist or is not searchable", 
				fullPath,
			)
			errMsgs = append(errMsgs,msg)
			return err
		}
		loc, isDup := slices.BinarySearchFunc(result, fullPath, comparePaths)
		if !isDup {
			result = slices.Insert(result, loc, fullPath)
		}
		return nil
	}
	for _, dir := range ctx.RecurseDirs {
		dir, err := fpath.Abs(dir)
		if err != nil {
			return result, ErrNoPwd
		}
		fpath.WalkDir(dir, addFiles)
	}
	if len(errMsgs) == 0 {
		return result, nil
	} else {
		return result, errors.New(strings.Join(errMsgs, "\n"))
	}
}

func isFatValid(r rune) bool {
	for _, runeRange := range FAT_RUNESET {
		if runeRange[0] <= r && r <= runeRange[1] {
			return true
		}
	}
	return false
}

func (ctx *context) restrictRuneset(s string) string {
	result := s
	toValidSubs := make(map[rune]string)
	if ctx.Strategy == "remove" {
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

func warn(msg string) {
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
		FileList: []string{},
		RecurseDirs: []string{},
		DryRun: false,
		Strategy: "represent",
		DoHelp: false,
		DoVersion: false,
	}
	index := 1
	isFlag := func(arg string) bool {
		return arg[0] == '-'
	}
	for index < len(args) {
		switch arg := args[index]; arg {
		case "-r", "--recurse":
			index++
			result.RecurseDirs = append(result.RecurseDirs, args[index])
		case "-n", "--dry-run":
			result.DryRun = true
		case "-s", "--strategy":
			index++
			result.Strategy = args[index]
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
				result.FileList = append(result.FileList, arg)
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

 Rename FILES such that all characters that are invalid in FAT file systems 
 (?,\,*,etc.) are removed.
 
 Options:
   -s,--strategy
   -h,--help
   -v,--version
   -r,--recurse DIRECTORY

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
		ctx.FileList, err = ctx.parseFileList()
		if err == ErrNoPwd {
			kaput(err)
		}
		if err != nil {
			warn(err.Error())
		}
		counter := 0
		for _, file := range ctx.FileList {
			oldName := fpath.Base(file)
			dirName := fpath.Dir(file)
			newName := ctx.restrictRuneset(
				strings.ToValidUTF8(oldName, "_INVALID_"),
			)
			newPath := fpath.Join(dirName, newName)
			if newPath == file {
				continue
			}
			if ctx.DryRun {
				format := "%s -> <<%s>>\n"
				if counter % 2 == 1 {
					format = "\x1b[2m%s -> <<%s>>\x1b[0m\n"
				}
				counter++
				fmt.Printf(
					format,
					file,
					newPath,
				)
			} else {
				os.Rename(file, newPath)
			}
		}
	}
}
