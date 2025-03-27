**findroot is a simple command-line tool that extracts root domains from a list of subdomains.**

**Features**
Automatically reads from tld.txt in the current directory by default.

Extracts root domains from subdomains.

Supports input from standard input (stdin).

Allows specifying an output file with -o.

**Usage**
Extract root domains from a file:

findroot
By default, it processes tld.txt and extracts root domains.

Using a custom input file with output redirection:

**cat input.txt | findroot -o rootdomains.txt**

Example
**Input (input.txt):**

sub.example.com  
test.example.com  
another.test.co.uk  

**cat input.txt | findroot -o rootdomains.txt**

**Output (rootdomains.txt):**

example.com  
test.co.uk

**Installation**

**git clone https://github.com/GDATTACKER-RESEARCHER/findroot && cd findroot && go build**
