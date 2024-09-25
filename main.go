package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

type TrieNode struct {
	Count    int
	Children map[string]*TrieNode
}

func main() {
	var inputFile, outputFile string
	var threshold int

	flag.StringVar(&inputFile, "i", "", "Input file (default is stdin)")
	flag.StringVar(&outputFile, "o", "", "Output file (default is stdout)")
	flag.IntVar(&threshold, "c", 2, "Minimum occurrence count")
	flag.Parse()

	// Reader: Input from file or stdin
	var reader *bufio.Reader
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	// Trie initialization
	root := &TrieNode{Children: make(map[string]*TrieNode)}

	// Buffer for large lines
	buf := make([]byte, 0, 4096)

	// Reading and processing input
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		if isPrefix {
			buf = append(buf, line...) // Handle long lines
			continue
		} else {
			buf = append(buf, line...)
			processLine(buf, root) // Process the line
			buf = buf[:0]          // Reset buffer
		}
		if err == io.EOF {
			break
		}
	}

	// Writer: Output to file or stdout
	var writer *bufio.Writer
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}

	// Traverse and write results
	traverseTrie(root, []string{}, threshold, writer)
	writer.Flush()
}

// processLine processes a line and inserts into the trie
func processLine(line []byte, root *TrieNode) {
	line = trimSpaces(line)
	if len(line) == 0 {
		return
	}

	// Lowercase conversion in place
	for i := 0; i < len(line); i++ {
		line[i] = byte(unicode.ToLower(rune(line[i])))
	}

	// Split line into labels (without allocations)
	labels := splitLabels(line)

	// Reverse labels for suffix trie
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}

	// Insert into trie
	node := root
	for _, labelBytes := range labels {
		label := string(labelBytes)
		if node.Children == nil {
			node.Children = make(map[string]*TrieNode)
		}
		if _, ok := node.Children[label]; !ok {
			node.Children[label] = &TrieNode{}
		}
		node = node.Children[label]
		node.Count++
	}
}

// trimSpaces trims leading and trailing spaces from a byte slice
func trimSpaces(b []byte) []byte {
	start, end := 0, len(b)-1
	for start <= end && (b[start] == ' ' || b[start] == '\t' || b[start] == '\r' || b[start] == '\n') {
		start++
	}
	for end >= start && (b[end] == ' ' || b[end] == '\t' || b[end] == '\r' || b[end] == '\n') {
		end--
	}
	return b[start : end+1]
}

// splitLabels splits the domain into labels without allocations
func splitLabels(line []byte) [][]byte {
	labels := make([][]byte, 0, 10)
	start := 0
	for i := 0; i <= len(line); i++ {
		if i == len(line) || line[i] == '.' {
			labels = append(labels, line[start:i])
			start = i + 1
		}
	}
	return labels
}

// traverseTrie traverses the trie and writes results that meet the threshold
func traverseTrie(node *TrieNode, labels []string, threshold int, writer *bufio.Writer) {
	if node.Count >= threshold && len(labels) > 0 {
		suffix := labelsToDomain(labels)
		writer.WriteString(suffix + "\n")
	}
	for label, child := range node.Children {
		traverseTrie(child, append(labels, label), threshold, writer)
	}
}

// labelsToDomain constructs the domain string from reversed labels
func labelsToDomain(labels []string) string {
	domainParts := make([]string, len(labels))
	for i := 0; i < len(labels); i++ {
		domainParts[i] = labels[len(labels)-1-i]
	}
	return joinLabels(domainParts)
}

// joinLabels joins labels into a domain string without extra allocation
func joinLabels(labels []string) string {
	totalLen := 0
	for _, label := range labels {
		totalLen += len(label) + 1 // +1 for '.'
	}
	totalLen-- // Remove the last '.'

	buf := make([]byte, totalLen)
	pos := 0
	for i, label := range labels {
		copy(buf[pos:], label)
		pos += len(label)
		if i < len(labels)-1 {
			buf[pos] = '.'
			pos++
		}
	}
	return string(buf)
}
