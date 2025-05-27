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
	"strconv"
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

/**
 * A basic search-and-replace operation. Replaces each instance of "Target" 
 * with each member of "Subs" in order, with all instances beyond 
 * length(Subs) being replaced by the last member of Subs.
 */
type replacement struct {
	Target string
	Subs []string
}

// Unrelated to the standard library's "context".
type context struct {
	FileList []string
	RecurseDirs []string
	DryRun bool
	Strategy string
	DoHelp bool
	DoVersion bool
	Rset runeset
	TruncLen int
	Replacements []replacement
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

var POSIX_PORTABLE_RUNESET = runeset{
	{ 0x2d, 0x2d },
	{ 0x2e, 0x2e },
	{ 0x5f, 0x5f },
	{ 0x30, 0x39 },
	{ 0x41, 0x5a },
	{ 0x61, 0x7a },
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

/**
 * Parses the targeting options into a list of files to rename (stored in a 
 * map as keys, for storing their new names as values later) and a list of 
 * all files in the same directory as a target file, with their name in 
 * lower case, for checking name collisions.
 */
// TODO: Write tests for this function
func (ctx *context) parseFileList() ([]string, map[string]bool, error) {
	errMsgs := []string{}
	nearbyFiles := make(map[string]bool)
	targets := []string{}
	addPath := func(pathList *[]string, path string) {
		loc, isDup := slices.BinarySearchFunc(*pathList, path, comparePaths)
		if !isDup {
			*pathList = slices.Insert(*pathList, loc, path)
		}
	}
	addRecursively := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			// This can only happen with the path passed to WalkDir
			errMsgs = append(
				errMsgs, fmt.Sprintf("Directory <<%s>> not searchable", path))
			return nil
		}
		lowerCasedPath := fpath.Join(fpath.Dir(path), strings.ToLower(entry.Name()))
		nearbyFiles[lowerCasedPath] = true
		addPath(&targets, path)
		return nil
	}
	for _, dir := range ctx.RecurseDirs {
		dir, err := fpath.Abs(dir)
		if err != nil {
			return targets, nearbyFiles, ErrNoPwd
		}
		fpath.WalkDir(dir, addRecursively)
	}
	for _, path := range ctx.FileList {
		if _, err := os.Stat(path); err != nil {
			msg := fmt.Sprintf( "File <<%s>> does not exist", path)
			errMsgs = append(errMsgs, msg)
		}
		path, err := fpath.Abs(path)
		if err != nil {
			return targets, nearbyFiles, ErrNoPwd
		}
		addPath(&targets, path)
	}
	if len(errMsgs) == 0 {
		return targets, nearbyFiles, nil
	} else {
		return targets, nearbyFiles, errors.New(strings.Join(errMsgs, "\n"))
	}
}

func inRuneset(r rune, rset runeset) bool {
	for _, runeRange := range rset {
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
			if !inRuneset(r, ctx.Rset) {
				toValidSubs[r] = ""
			}
		}
	} else {
		for _, r := range result {
			if !inRuneset(r, ctx.Rset) {
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

/**
 * Truncate to the number of bytes given by the context, but don't break in 
 * the middle of a rune.
 */
func (ctx *context) truncate(name string) string {
	if ctx.TruncLen < 1 {
		return name
	}
	for byteIndex, _ := range name {
		if byteIndex > ctx.TruncLen {
			return name[:byteIndex]
		}
	}
	return name
}

func (ctx *context) searchAndReplace(name string) string {
	for _, rep := range ctx.Replacements {
		for i:=0; i<len(rep.Subs)-1; i++ {
			name = strings.Replace(name, rep.Target, rep.Subs[i], 1)
		}
		name = strings.ReplaceAll(name, rep.Target, rep.Subs[len(rep.Subs)-1])
	}
	return name
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
		Rset: FAT_RUNESET,
		TruncLen: 0,
		Replacements: []replacement{},
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
		case "-p", "--portable":
			result.Rset = POSIX_PORTABLE_RUNESET
		case "-c", "--replace":
			index++
			target, subs, ok := strings.Cut(args[index], ":")
			if !ok {
				return result, errors.New(fmt.Sprintf(
					"Incorrect syntax (<<%s>>) for replacement: please use '{TARGET}:{REPLACEMNT1},{REPLACEMNT2},...'",
					args[index],
				))
			}
			result.Replacements = append(
				result.Replacements,
				replacement{ target, strings.Split(subs, ",") },
			)
		case "-t", "--truncate":
			index++
			var err error
			result.TruncLen, err = strconv.Atoi(args[index])
			if err != nil || result.TruncLen < 1 {
				return result, errors.New(fmt.Sprintf(
					"Invalid truncation length: <<%s>>",
					args[index],
				))
			}
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
   -n,--dry-run
   -p,--portable

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
		targets, nearbyFiles, err := ctx.parseFileList()
		if err == ErrNoPwd {
			kaput(err)
		}
		if err != nil {
			warn(err.Error())
		}
		numRenamed := 0
		for _, file := range targets {
			// Calculate the new name
			oldName := fpath.Base(file)
			dirName := fpath.Dir(file)
			newName := ctx.truncate(ctx.restrictRuneset(ctx.searchAndReplace(
				strings.ToValidUTF8(oldName, "_INVALID_"),
			)))
			newPath := fpath.Join(dirName, newName)
			// Check that we want to rename this file
			if newPath == file {
				continue
			}
			lowerCasedPath := fpath.Join(dirName, strings.ToLower(newName))
			if nearbyFiles[lowerCasedPath] {
				warn(fmt.Sprintf(
					"Path <<%s>> would be renamed to\n  <<%s>>,\nwhich collides with\n  <<%s>>\nSkipping...",
					file,
					newPath,
					lowerCasedPath,
				))
				continue
			}
			nearbyFiles[lowerCasedPath] = true
			// Do the rename, or just print what would happen
			numRenamed++
			if ctx.DryRun {
				format := "%s -> <<%s>>\n"
				if numRenamed % 2 == 0 {
					format = "\x1b[2m%s -> <<%s>>\x1b[0m\n"
				}
				fmt.Printf(
					format,
					file,
					newPath,
				)
			} else {
				os.Rename(file, newPath)
			}
		}
		fmt.Printf("%d files renamed!\n", numRenamed)
	}
}
