# bwc - bitwise calculator

Because bitwises are trick.

# The tool

The language was designed in a such a way that
makes easy to copy-paste bitwises from code
in most programming languages (at least from
C descendents). Take a look in the code below:

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

# paste the rest of the statements
bwc> X = (X | (X << 8)) & 0x00ff00ff00ff00ff
bin: 11111111
hex: ff
bwc> X = (X | (X << 4)) & 0x0f0f0f0f0f0f0f0f
bin: 111100001111
hex: f0f
bwc> X = (X | (X << 2)) & 0x3333333333333333
bin: 11001100110011
hex: 3333
bwc> X = (X | (X << 1)) & 0x5555555555555555
bin: 101010101010101
hex: 5555

# hmm, bit interleaving =)
```

# The language

```bnf
letter		= "a".."z" | "A".."Z";
alphanum	= letter | decdigit;

decdigit 	= "0".."9";
hexdigit 	= decdigit | "a".."f";
bindigit 	= "0" | "1";

decimal	= decdigit { decdigit };
hexadecimal	= "0x" hexdigit { hexdigit };
binary		= "0b" bindigit { bindigit };

number		= decimal | hexadecimal | binary;
ident		= letter {alphanum};
binaryop	= "&" | "|" | "^" | "<<" | ">>";
unaryop		= "~";
mathexpr	= [ "(" ] unaryexpr | binaryexpr [ ")" ];
operand		= expr | binaryexpr | unaryexpr | number;
binaryexpr	= operand binaryop operand;
unaryexpr	= unaryop operand;
expr		= number | ident | mathexpr;

assignment	= ident "=" expr;

grammar		= assignment | expr;
```