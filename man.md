# owl(1)

*The moon is bright. The owls are ready for the hunt.*

Owl is a file renaming tool; it removes character from file names which are 
not compatible with `FAT` file systems (typically `FAT32` or `exFAT`).

```
owl COMMAND OPTIONS FILE1 FILE2 ...
```

Rename `FILE`s given at the command line such that all `FAT`-incompatible 
characters are removed. By default it replaces invalid characters with 
`_U{Unicode Code Point}_`, but this can be changed, see the `--strategy` 
flag below.


## Options:
### -s, \-\-strategy STRATEGY
Change what happens to invalid characters. Choices are:

- "remove": just remove the characters. If left with an empty string, rename 
  to "\_empty_".
- "represent": replace each character with "\_Unum_", where "num" is the 
  Unicode Code Point of the character.

### -n, \-\-dry-run
Causes Owl to just print a representation of what would be done without 
actually renaming any files.

### -v, \-\-version
Show version information.

### -h, \-\-help
Show a small help message.


## Edge Cases & Other Tidbits
## What doesn't FAT allow?
Mainly `*<>\|/:?'`

### What about whitespace?
Some of these (like tab, line feed, and carriage return) are not valid under 
FAT and will be removed.

If a sequence of invalid characters surrounded by spaces is in one of the 
file names, and the `remove` strategy is being used, then you will be left 
with multiple whitespace characters in a row. For example,
```
Really questionable ?? filename
```
will become
```
Really questionable  filename
```
with a double space before `filename`.

### What about invalid UTF-8/UTF-16?
This is removed and replaced with "\_INVALID_". This replacement is done 
before anything else.

