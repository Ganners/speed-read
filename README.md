Speed Reading in the Terminal
=============================

This is a speed reader which will present some plain text to you word by word
by flashing individual words in the center of the terminal. Super simple and I
believe quite effective. Should be better with an increased font size.

Example usage:
--------------

```bash
go get github.com/ganners/speed-read
DEMO_BOOK=$GOPATH/src/github.com/ganners/speed-read/italian-villas.txt
cat $DEMO_BOOK | speed-read -wpm=600 -lines=$(tput lines) -cols=$(tput cols)
```

For some free plain-text books, check out http://www.gutenberg.org/. An example
from there is included in the repository.
