package pkg

import (
	"strings"
	"unicode"

	"github.com/pkg/errors"
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

	// Three characters used to represent the ellipsis in a string.
	threeChars = 3
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

	// preContext and postContext are string variables used to hold portions of the original
	// string before and after the ellipsis, respectively. These variables are typically used
	// within the `Ellipsize` function to store the preserved segments of the string before
	// and after the truncation process.
	//
	// - preContext stores the substring from the start of the string up to the `from` index.
	// - postContext stores the substring from the `to` index onward to the end of the string.
	//
	// These variables help in constructing the final ellipsized string that combines
	// the preserved parts with an ellipsis in the middle.
	preContext, postContext string
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
				if cmdlineRune[i+1] == qr {
					builder.WriteRune(qr)
					i = i + 1
					continue NextRune
				}
			}
		}
		// If the current rune is a quote rune, we need to start a quoted string.
		for _, qr := range quoteRunes {
			if head == qr {
				lastQuoteRune = qr
				lastQuotePos = i
				continue NextRune
			}
		}
		// If the current rune is a space, we have reached the end of a token.
		if unicode.IsSpace(head) && lastQuoteRune == 0 {
			tokenizedLines = append(tokenizedLines, builder.String())
			builder.Reset()
			continue
		}
		builder.WriteRune(head)
	}
	if lastQuoteRune > 0 {
		context := Ellipsize(
			lastQuotePos-errUnterminatedQuote,
			lastQuotePos+errUnterminatedQuote+1,
			cmdline,
		)
		return nil, errors.Errorf("Unterminated quote at position %d: %s", lastQuotePos, context)
	}
	if builder.Len() > 0 {
		tokenizedLines = append(tokenizedLines, builder.String())
	}
	return tokenizedLines, nil
}

// Ellipsize shortens a string by replacing the middle part with an ellipsis ("...").
// This function takes a string `str` and truncates it such that the start of the string
// is preserved up to the `from` index, and the end of the string is preserved from the `to` index.
// The middle portion of the string between `from` and `to` is replaced with an ellipsis.
// If the `from` and `to` indices do not leave enough room for the ellipsis, the original string may be returned.
//
// Parameters:
//   - from: The index at which to start preserving the string from the beginning.
//   - to: The index at which to start preserving the string from the end.
//   - val: The original string to be ellipsized.
//
// Returns:
//   - A new string where the middle portion between `from` and `to` is replaced with an ellipsis,
//     or the original string if the indices do not allow for proper truncation.
func Ellipsize(from, to int, val string) string {
	preContextIndex := from
	if preContextIndex <= threeChars {
		preContextIndex = 0
		preContext = ""
	} else {
		preContext = "..."
	}
	postContextIndex := to
	if postContextIndex >= (len(val) - threeChars) {
		postContextIndex = len(val)
		postContext = ""
	} else {
		postContext = "..."
	}
	return preContext + val[preContextIndex:postContextIndex] + postContext
}
