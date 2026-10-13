package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Prompt reads a single line of input from the command's input stream after
// writing the given message to its output stream.
func Prompt(w io.Writer, in io.Reader, message string) (string, error) {
	_, err := fmt.Fprint(w, message)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
