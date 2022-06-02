package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vilmibm/cyberbog/bogdb"
)

func _main(b *bogdb.BogDB, in io.Reader, out io.Writer) error {
	// TODO do not assume something on STDIN
	input, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	if len(input) != 0 {
		err = b.Inter(input)
		if err != nil {
			return fmt.Errorf("failed to inter: %w", err)
		}
		return nil
	}

	output, err := b.Exhume()
	if err != nil {
		return fmt.Errorf("failed to exhume: %w", err)
	}

	fmt.Fprintln(out, output)

	return nil
}

func main() {
	// TODO take a rootpath arg
	b, err := bogdb.NewBogDB("/tmp/bog", time.Now().UnixNano(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to make a bogdb: %s", err.Error())
	}

	err = _main(b, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
