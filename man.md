# owl(1)

*The moon is bright. The owls are ready for the hunt.*

Owl is a file renaming tool; it removes character from file names which are 
not compatible with `FAT` file systems (typically `FAT32` or `exFAT`).

```
owl OPTIONS FILE1 FILE2 ...
```

Rename `FILE`s given at the command line such that all `FAT`-incompatible 
characters are removed. By default it replaces invalid characters with 
`_U{Unicode Code Point}_`, but this can be changed, see the `--strategy` 
flag below.


## OPTIONS:
### -r, \-\-recurse DIRECTORY
Recursively search DIRECTORY for files/directories to rename. Search 
includes DIRECTORY itself.

### -s, \-\-strategy STRATEGY
Change what happens to invalid characters. Choices are:

- "remove": just remove the characters. If left with an empty string, rename 
  to "\_empty_".
- "represent": replace each character with "\_Unum_", where "num" is the 
  Unicode Code Point of the character.

### -e,\-\-valid-set fat|posix|shell
Select which characters are considered "invalid" or "bad". Options are:
 - fat: Characters which are valid in file names on FAT file systems. This 
   is the default.
 - posix: The POSIX Portable Filename Character Set `A-Za-z0-9.-_`
 - shell: Characters that *should* not need to be quoted when used in a 
   shell, i.e. all characters apart from control characters, whitespace, and 
   most punctuation.

### -p,\-\-portable
An alias for "\-\-valid-set posix".

### -t,\-\-truncate LENGTH
Truncate all given file names to at most `LENGTH` bytes. Note that if you 
have file names with larger Unicode code points (like CJK characters), you 
may get unexpected results with this option, though it will not cut the file 
name part way through a character.

### -c,\-\-replace TARGET:REPLACEMENT1,REPLACEMENT2,...
Remove every instance of the string `TARGET` in every file name and replace 
the first instance with `REPLACEMENT1`, the second with `REPLACEMENT2`, and 
so on. If there are N `REPLACEMENT`'s and more than N instances of `TARGET` 
in a file name, all instances after the Nth are replaced by the last 
replacement. Replacements can be empty. See EXAMPLES below for this one.

Remember to quote `TARGET:REPLACEMENT1,...` so that your shell doesn't 
misinterpret the colon.

### -n, \-\-dry-run
Causes Owl to just print a representation of what would be done without 
actually renaming any files. Checks for name collisions with existing files 
or other renaming operations are still done.

### -v, \-\-version
Show version information.

### -h, \-\-help
Show a small help message.


## EXAMPLES
Suppose we have a file called
```
PodEp 37: Why are there so many spaces in this file name?
```
in the current working directory. We will refer to it as FILE Then:

- The command
  ```
  owl --replace ' :,-,_' FILE
  ```
  will rename FILE to  
  `PodEp37:-Why_are_there_so_many_spaces_in_this_file_name_U3F_`
- The command
  ```
  owl --truncate 48 --replace ' :,-,_' FILE
  ```
  will rename FILE to
  `PodEp37:-Why_are_there_so_many_spaces_in_this_fi`
- The command
  ```
  owl --portable --strategy remove --truncate 48 FILE
  ```
  will rename FILE to
  `PodEp37:Whyaretheresomanyspacesinthisfilename`

Now suppose we have files "IMPORTANT FILE?" and "Important File".
- If you run
  ```
  owl --strategy remove "IMPORTANT FILE?"
  ```
  owl will just return an error, saying that the new name would collide with 
  "important file".
- If you run
  ```
  owl --strategy remove --valid-set posix "IMPORTANT FILE?" "Important File"
  ```
  then owl will will rename one file, but refuse to rename the other, as it 
  would conflict with the now-renamed first file.
  Note you can use the `--dry-run` flag to check for collisions like this 
  before renaming anything.


## EDGE CASES & OTHER TIDBITS
### What order does owl truncate, remove characters, etc.?
Owl does things in the following order:

1. Remove invalid UTF-8
2. Replace stuff (see `--replace` flag)
3. Deal with invalid characters (see `--strategy` flag)
4. Truncate file name

The order in which you use the flags does not change this.

### What if the new name Owl chooses for a file already exists?
Owl checks before each renaming operation to see if the new name already 
exists, and skips that file if it does. The check is also case-insensitive.
See [EXAMPLES](#examples) for more details.


### What doesn't FAT allow?
Mainly `*<>\|/:?'`

### What about whitespace?
Some of these (like tab, line feed, and carriage return) are not valid under 
FAT and will be removed.

If a sequence of invalid characters surrounded by spaces is in one of the 
file names, and the `remove` strategy is being used, then you will be left 
with multiple whitespace characters in a row. For example,
```
"Really questionable ?? filename"
```
will become
```
"Really questionable  filename"
```
with a double space before "filename".

### What about invalid UTF-8?
This is removed and replaced with "\_INVALID_". This replacement is done 
before anything else.

### What about UTF-16/Windows?
Owl has only been tested so far on Linux with a UTF-8 locale. UTF-16 support 
on Windows is likely possible, but has not been tested. Please contact me if 
you would like Windows support (see below).

### What about other restrictions on file names?
Windows (and by extension FAT file systems, which come from Microsoft) may 
have other restrictions on file names, like specific disallowed names like 
"CON" or names which ends with a ".". Owl does not check for such things; 
you should look for those using a tool like `find`.

For completeness, the disallowed names are
```
CON,AUX,COM1,COM2COM3,COM4,LPT1,LPT2,LPT3,PRN,NUL
```

### What about symbolic links (a.k.a. symlinks)?
Owl will rename the *link*, but not the file that the link points to, even 
if the link is broken. If you are having issues with broken symbolic links 
and are on a POSIX/UNIX-like system, you can use
```
find -L DIRECTORY -type l
```
to find all broken links in DIRECTORY.


## AUTHORS & COPYRIGHT
Written by user SixteenThousand of github.com (email 
thomsixteenthousand@gmail.com) and licensed under the GNU General Public 
License version 3 (GPL v3). See <https://www.gnu.org/licenses/> for more 
information.
