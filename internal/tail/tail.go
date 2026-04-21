package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Options configures tail behaviour.
type Options struct {
	// N is the number of lines to return from the end of the file.
	N int
}

// DefaultOptions returns sensible defaults for tail.
func DefaultOptions() Options {
	return Options{N: 20}
}

// ReadLastN returns the last n lines from r.
// It reads the entire stream into a ring buffer of size n.
func ReadLastN(r io.Reader, n int) ([]string, error) {
	if n <= 0 {
		return nil, fmt.Errorf("tail: n must be > 0, got %d", n)
	}

	buf := make([]string, n)
	pos := 0
	count := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		buf[pos%n] = line
		pos++
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scan error: %w", err)
	}

	if count == 0 {
		return []string{}, nil
	}

	size := n
	if count < n {
		size = count
	}

	result := make([]string, size)
	start := pos % n
	if count < n {
		start = 0
	}
	for i := 0; i < size; i++ {
		result[i] = buf[(start+i)%n]
	}
	return result, nil
}

// ReadLastNFromFile opens the named file and returns the last n lines.
func ReadLastNFromFile(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail: open %q: %w", path, err)
	}
	defer f.Close()
	return ReadLastN(f, n)
}
