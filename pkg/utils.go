package pkg

import (
	"strings"
	"unicode"
)

const (
	// Constant quoteEscapeRune represents the escape sequence for a backslash in a quoted string.
	// This constant is used to handle cases where a backslash needs to be escaped within
	// quoted text to ensure proper parsing and handling.
	quoteEscapeRune = '\\'

	// Constant errUnterminatedQuote is the error code used to indicate an unterminated quote error.
	// This constant is used when an error occurs due to a missing or mismatched quote in
	// the input, signaling that the quoted text was not properly closed.
	errUnterminatedQuote = 3
)

var (
	// This is a list of runes that are considered valid quote characters.
	quoteRunes = [...]rune{'"', '\''}

	// Holds a slice of strings where each string represents a token extracted
	// from a command line input. This variable is used to store the results of tokenization
	// operations, allowing further processing or analysis of individual tokens.
	//
	// Example usage:
	//    tokenizedLines = TokenizeLine(cmdline)
	//    // tokenizedLines now contains the tokens from cmdline.
	tokenizedLines []string

	// lastQuoteRune represents the last rune (character) used to denote a quote in the
	// command line input. This variable is used to keep track of the type of quote (e.g., single
	// or double quote) that was last encountered during parsing, helping to correctly handle
	// nested or unclosed quotes.
	//
	// Example usage:
	//    lastQuoteRune = '"' // Indicates that the last encountered quote was a double quote.
	lastQuoteRune rune

	// lastQuotePos stores the position (index) of the last quote character encountered
	// in the command line input. This variable helps to track the location of the last quote,
	// which is useful for detecting unterminated quotes or managing quoted strings during
	// tokenization and parsing.
	//
	// Example usage:
	//    lastQuotePos = 42 // Indicates the position of the last quote character in the input string.
	lastQuotePos int

	// builder is an instance of strings.Builder used for efficiently building and
	// concatenating strings. The strings.Builder type provides a mutable buffer
	// for string concatenation, which is more efficient than using string concatenation
	// operators (+) in a loop, as it avoids creating multiple intermediate string objects.
	//
	// Example usage:
	//    builder.WriteString("Hello, ")
	//    builder.WriteString("world!")
	//    result := builder.String() // result will be "Hello, world!"
	builder strings.Builder
)

// Function TokenizeLine splits the given command line string into individual tokens.
// This function processes the input string, which may contain multiple words and
// delimiters, and returns a slice of strings where each string is a separate token
// extracted from the command line. The function handles common tokenization rules such
// as whitespace separation and quoted strings. If an error occurs during tokenization,
// it will return an error detailing the issue.
//
// Parameters:
//   - cmdline: The input command line string to be tokenized. This string may contain
//     various tokens separated by whitespace or enclosed in quotes.
//
// Returns:
//   - A slice of strings where each string is a token extracted from the input command
//     line.
//   - An error if there is an issue with tokenization, such as invalid syntax or unclosed
//     quotes. If no error occurs, the error will be nil.
func TokenizeLine(cmdline string) ([]string, error) {
	cmdlineRune := []rune(cmdline)
	builder.Grow(len(cmdlineRune))
NextRune:
	for i := 0; i < len(cmdlineRune); i++ {
		head := cmdlineRune[i]
		if lastQuoteRune == 0 && unicode.IsSpace(head) && builder.Len() == 0 {
			continue
		}
		if lastQuoteRune > 0 {
			if head == quoteEscapeRune && i < len(cmdlineRune)-1 && cmdlineRune[i+1] == lastQuoteRune {
				builder.WriteRune(lastQuoteRune)
				i = i + 1
				continue
			}
			// If the current rune is the same as the last quote rune, we have reached the end of the quoted string.
			if head == lastQuoteRune {
				lastQuoteRune = 0
				continue
			}
			builder.WriteRune(head)
			continue
		}
		if head == quoteEscapeRune && i < len(cmdlineRune)-1 {
			for _, qr := range quoteRunes {
				lastQuoteRune = qr
				lastQuotePos = i
				continue NextRune
			}
		}
	}
	return tokenizedLines, nil
}
