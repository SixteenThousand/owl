# Owl

*The moon is bright. The owls are ready for the hunt.*

Owl is a command line tool for removing unwanted characters from file names.
Say you want to backup some of your files on to a memory stick. Most likely 
that memory stick uses `FAT32` or `exFAT` as its file system, as this is the 
most compatible with the most operating systems. But if any of your files 
have special characters like `:` or `?` in their names, this won't work (see 
<https://learn.microsoft.com/en-gb/windows/win32/fileio/exfat-specification> 
for more details).

Alternatively, say you have a file that needs to be on an older computer 
that will not render some characters (like multi-character emoji) 
correctly<sup>1</sup>.

And what if you have either of these problems with potentially many files in 
a directory? Renaming every file would be tedious, and even the best bulk 
renaming tool will struggle to deal with all the possible invalid 
characters, at least if you don't want to manually write a very complex 
[regular expression](https://en.wikipedia.org/wiki/Regular_expression) every 
time you encounter this problem.

Owl will hunt down any files with these issues and rename them as 
appropriate - see `man.md` for usage.


---

1. Or you just don't like having emoji in a file name.
