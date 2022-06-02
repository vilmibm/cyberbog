package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vilmibm/cyberbog/bogdb"
)

type opts struct {
	In    io.Reader
	Out   io.Writer
	Verb  string
	BogDB *bogdb.BogDB
}

func _main(o opts) error {
	if o.Verb == "inter" {
		input, err := io.ReadAll(o.In)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		if len(input) != 0 {
			err = o.BogDB.Inter(input)
			if err != nil {
				return fmt.Errorf("failed to inter: %w", err)
			}
		}

		return nil
	}

	output, err := o.BogDB.Exhume()
	if err != nil {
		return fmt.Errorf("failed to exhume: %w", err)
	}

	fmt.Fprintln(o.Out, string(output))

	return nil
}

func main() {
	pathFlag := flag.String("path", "/tmp/bog", "path to bog")

	flag.Parse()

	verbArg := flag.Arg(0)
	if verbArg != "exhume" && verbArg != "inter" {
		fmt.Fprintln(os.Stderr, "verb should be one of inter or exhume")
		os.Exit(1)
	}

	b, err := bogdb.NewBogDB(*pathFlag, time.Now().UnixNano(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to make a bogdb: %s\n", err.Error())
	}

	o := opts{
		In:    os.Stdin,
		Out:   os.Stdout,
		Verb:  verbArg,
		BogDB: b,
	}

	err = _main(o)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
}
