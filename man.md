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
