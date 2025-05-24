# Infrequently Asked Questions
No one but has asked these questions, but I'm giving answers anyway, mainly 
so I don't forget them in future.

## Why didn't you use Punycode?
[Punycode](https://en.wikipedia.org/wiki/Punycode) is a web standard for 
converting general Unicode text into ASCII, for use in domain names, and for 
nothing else. As the Wikipedia page says, it maps "München" to "Mnchen-3ya", 
which I personally find harder to understand than "M_Ufc_nchen". It also 
would add a dash to the end of pure 
[ASCII](https://en.wikipedia.org/wiki/ASCII) text, so "London" would become 
"London-", and its output does not handle characters that are bad for 
shells, like "?".

## Why didn't you use `base64/32/16`?
These are internet encoding standards used to convert any binary data to 
ASCII text, which sounds perfect for the job. However, their outputs are not 
human readable, and Owl is designed so that you will at least know roughly 
what the file was originally called.
See <https://www.rfc-editor.org/rfc/rfc4648.html#section-10> for examples.

## Why not just `mmv`? Or `zsh`?
[mmv](https://github.com/rrthomas/mmv) is command line tool to do a kind of 
"glob search and replace" with filenames. [zsh](https://zsh.sourceforge.io/) 
is a shell, Bash or the PowerShell, but it has a tool to do a similar thing 
called `zmv`. These require writing a regular expression of glob pattern to 
represent the changes you wish to make, which is fiddly work when what you 
want is to exclude potentially thousands of characters, not to mention 
potential issues with hidden files (a.k.a. "dot files") and shell globbing.

### Other alternatives
- <https://github.com/DDR0/fuseblk-filename-fixer>: This actually does 
  pretty much what Owl does, but is more specialised to its author's needs.

## Tree
This is used for testing. It should be available on any Linux system from 
your distribution's repositories. If not, the source code can be downloaded 
from the GitLab repository at 
<https://gitlab.com/OldManProgrammer/unix-tree> or from the home page of the 
project at <https://oldmanprogrammer.net/source.php?dir=projects/tree>.
