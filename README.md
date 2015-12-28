# A SIMPLE FILE TYPE COUNTER TOOL

The tool is a very simple command line tool written in Go.
The purpose of the tool is to count all files within a specific directory, by file types (extensions).
This program can be easily extended to count (and sort) by filesize, timestamps, etc.

## Program Usage

```
go run filetypecount.go -dir=YourDirectory
```

where `YourDirectory` is a directory path you want to check.

## Sample Usage

For a Go installation folder (Windows OS), the output would be as follow:
```
go run filetypecount.go -dir=c:\go
scanning: c:\go
...
.go             => 5304
(no ext)                => 549
.s              => 347
.a              => 333
.sample         => 144
.png            => 132
.article                => 106
.golden         => 103
.c              => 91
.src            => 74
.html           => 63
.expected               => 51
...

Total files: 8012, unique file types: 130
```
