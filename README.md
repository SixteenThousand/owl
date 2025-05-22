# Owl

*The moon is bright. The owls are ready for the hunt.*

Owl is a command line tool for removing unwanted characters from file names.
Say you want to backup some of your files on to a memory stick. Most likely 
that memory stick uses `FAT32` or `exFAT` as its file system, as this is the 
most compatible with the most operating systems. But if any of your files 
have special characters like `:` or `?` in their names, this won't work 
[1](#links).

Alternatively, say you have a file that needs to be on an older computer 
that will not render some characters (like multi-character emoji) correctly 
(or you just don't like having emoji in a file name).

And what if you have either of these problems with potentially many files in 
a directory? Renaming every file would be tedious, and even the best bulk 
renaming tool will struggle to deal with all the possible invalid 
characters, at least if you don't want to manually write a very complex 
regular expression every time you encounter this problem.

Owl will hunt down any files with these issues and rename them as 
appropriate - see `man.md` for usage.


---

## Links

1. A Stack Exchange question which outlines some of the problems in this area: <https://unix.stackexchange.com/questions/779052/removing-hidden-control-characters-in-filenames>
2. Another, similar stack exchange question: <https://stackoverflow.com/questions/1976007/what-characters-are-forbidden-in-windows-and-linux-directory-names>
3. FAT Specification: 
   <https://learn.microsoft.com/en-gb/windows/win32/fileio/exfat-specification> 
   (see section 7.7.3)
4. More file system specifications: <https://learn.microsoft.com/en-us/troubleshoot/windows-client/backup-and-storage/fat-hpfs-and-ntfs-file-systems?source=recommendations>
5. mmv: <https://github.com/rrthomas/mmv>
6. PathShortener: <https://github.com/ElectricRCAircraftGuy/eRCaGuy_PathShortener>
7. fuseblk-filename-fixer: <https://github.com/DDR0/fuseblk-filename-fixer>
