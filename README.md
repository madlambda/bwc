# bwc - bitwise calculator

Because I'm dumb.

# language

The language was designed in a such a way that
makes easy to copy-paste bitwises from code
in most programming languages (at least from
C generation). Take a look in the code below:

```go
func stripe(val uint32) uint64 {
	X := uint64(val)
	X = (X | (X << 16)) & 0x0000ffff0000ffff
	X = (X | (X << 8)) & 0x00ff00ff00ff00ff
	X = (X | (X << 4)) & 0x0f0f0f0f0f0f0f0f
	X = (X | (X << 2)) & 0x3333333333333333
	X = (X | (X << 1)) & 0x5555555555555555
	return X
}
```

To understand this code you should define X and
then copy the lines or parts of the bitwise you 
are interested in. For example:

```
$ bwc
# Define the variable
bwc> X = 0xff
bin: 11111111
hex: ff

# Paste the code
bwc> X = (X | (X << 16)) & 0x0000ffff0000ffff

# eval X
bwc> X
bin: 11111111
hex: ff
bwc>
```