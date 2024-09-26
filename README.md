# Subtrie

## Overview

The **Subtrie** tool is designed to process domain-like strings (such as hostnames) and efficiently count their occurrences using a trie (prefix tree). This program is particularly useful for extracting domain suffixes that appear frequently in a large dataset. It allows users to specify a minimum occurrence count (`-c`) to filter domain suffixes based on their frequency.

The tool supports reading from both files and standard input, and can output the result to a file or standard output.

## Features

- **Trie Data Structure**: Efficiently stores and counts domain suffixes.
- **Minimum Occurrence Count**: Set a threshold (`-c`) to only output domain suffixes that meet or exceed this count.
- **Input/Output Options**: 
  - Read from a specified file or from standard input.
  - Write results to a file or to standard output.
- **Handles Large Inputs**: Uses buffered reading and processing to handle large files efficiently.

## Installation

```bash
go install github.com/xhzeem/subtrie@latest
```

## Usage

The tool can be run from the command line. Below are the usage options:

```bash
subtrie [options]
```

### Command-line Options

- `-i` (optional): Specify the input file containing domain strings. If omitted, input is taken from `stdin`.
- `-o` (optional): Specify the output file to write results to. If omitted, output is written to `stdout`.
- `-c` (optional): Specify the minimum count for domain suffixes to be included in the output. Default is `2`.

### Examples
Suppose we have a file `domains.txt` containing the following domains:

```
a.mail.example.com
b.mail.example.com
mail.example.com
www.example.com
a.blog.example.com
b.blog.example.com
c.blog.example.com
blog.example.com
example.com
test.sub.blog.example.com
deep.sub.sub.blog.example.com
```

Now, let's run the tool with a threshold of `2` to count suffixes that appear at least twice.

#### Command
```bash
./subtrie -i domains.txt -c 2
```

#### Expected Output
```
example.com
blog.example.com
mail.example.com
```

### Explanation:
- **example.com**: Appears 11 times as a suffix across all the domains (e.g., in `www.example.com`, `a.mail.example.com`, etc.).
- **blog.example.com**: Appears 6 times (in `a.blog.example.com`, `b.blog.example.com`, `c.blog.example.com`, etc.).
- **mail.example.com**: Appears 3 times (in `a.mail.example.com`, `b.mail.example.com`, `mail.example.com`).

Suffixes like `sub.blog.example.com` and `deep.sub.sub.blog.example.com` only appear once, so they are not included in the output as their count doesn't meet the threshold of `2`.

This example showcases the ability of the tool to handle multi-level subdomains and group them by their suffixes, counting occurrences based on the reversed structure of the domains.
