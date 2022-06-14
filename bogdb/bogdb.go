package bogdb

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

type Fragment struct {
	InterredAt time.Time
	Contents   []byte
}

type BogDB struct {
	rootPath string
	rand     *rand.Rand
	now      func() time.Time
}

func NewBogDB(rootPath string, randSeed int64, now func() time.Time) (*BogDB, error) {
	if now == nil {
		now = time.Now
	}
	err := os.MkdirAll(rootPath, 0750)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("could not create '%s': %w", rootPath, err)
	}
	b := &BogDB{
		rand:     rand.New(rand.NewSource(randSeed)),
		rootPath: rootPath,
		now:      now,
	}
	return b, nil
}

func (b *BogDB) timeToIntensity(t time.Time) int {
	return int(math.Log10(b.now().Sub(t).Seconds())) + 1
}

// Inter breaks up data into arbitrarily corroded fragments and stores them on the filesystem.
func (b *BogDB) Inter(data []byte) error {
	return b.interAt(data, b.now())
}

func (b *BogDB) interAt(data []byte, ts time.Time) error {
	fragName := fmt.Sprintf("%d", b.rand.Int63n(math.MaxInt64-100)+100)
	fragDir := path.Join(
		b.rootPath,
		string(fragName[0]), string(fragName[1]), string(fragName[2]))
	err := os.MkdirAll(fragDir, 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create '%s': %w", fragDir, err)
	}

	fragPath := path.Join(fragDir, string(fragName))
	f, err := os.Create(fragPath)
	if err != nil {
		return fmt.Errorf("failed to open fragment for writing '%s': %w", fragPath, err)
	}
	fmt.Fprintf(f, "%s\n", ts.Format(time.RFC3339))
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write fragment '%s': %w", fragPath, err)
	}

	return nil
}

// Exhume reads a fragment from disk, further fragmenting and corroding it. One fragment is returned, the rest re-interred. Fragments are read, deleted from disk, then re-written as needed. This absolutely can lead to lost data; this is by design.
func (b *BogDB) Exhume() (*Fragment, error) {
	dir := b.rootPath
	var absPath string
	var rawFragment []byte
	for rawFragment == nil {
		var err error
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read dir '%s': %w", dir, err)
		}
		if len(entries) == 0 {
			// TODO rmdir or continue or ...?
			return nil, nil
		}
		selected := b.rand.Intn(len(entries))
		for ix, entry := range entries {
			if ix != selected {
				continue
			}
			absPath = path.Join(dir, entry.Name())
			if entry.IsDir() {
				dir = absPath
				continue
			}
			rawFragment, err = os.ReadFile(absPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file '%s': %w", absPath, err)
			}
		}
	}

	fragment, err := parseFragment(rawFragment)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fragment '%s': %w", absPath, err)
	}

	err = os.Remove(absPath)
	if err != nil {
		return fragment, fmt.Errorf("failed to remove fragment '%s': %w", absPath, err)
	}

	fragments := b.fragment(fragment.Contents)

	fragIx := b.rand.Intn(len(fragments))
	var chosenFragment []byte
	for ix, f := range fragments {
		if ix == fragIx {
			chosenFragment = f
		} else {
			b.interAt(f, fragment.InterredAt)
		}
	}

	return &Fragment{
		InterredAt: fragment.InterredAt,
		Contents:   b.corrode(chosenFragment, b.timeToIntensity(fragment.InterredAt)),
	}, nil
}

func parseFragment(rawFragment []byte) (*Fragment, error) {
	split := strings.Split(string(rawFragment), "\n")
	ts, err := time.Parse(time.RFC3339, split[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse fragment '%s': %w", string(rawFragment), err)
	}

	return &Fragment{
		InterredAt: ts,
		Contents:   []byte(strings.Join(split[1:], "\n")),
	}, nil
}

func (b *BogDB) corrode(data []byte, intensity int) []byte {
	fmt.Printf("DBG %#v\n", "doing a corrode")
	fmt.Printf("DBG %#v\n", intensity)
	fmt.Printf("DBG %#v\n", string(data))
	rs := bytes.Runes(data)
	// numCorrodes := int(1 / len(rs) * intensity)
	numCorrodes := intensity
	fmt.Printf("DBG numC %#v\n", numCorrodes)
	for x := 0; x < numCorrodes; x++ {
		ix := b.rand.Intn(len(rs))
		r := rs[ix]
		mod := b.rand.Int31() - (math.MaxInt32 / 2)
		rs[ix] = r + mod
	}

	fmt.Printf("DBG %#v\n", string(rs))

	return []byte(string(rs))
}

func (b *BogDB) fragment(data []byte) [][]byte {
	output := [][]byte{}
	splits := b.rand.Intn(int(len(data)/4)) + 1

	startIX := 0
	endIX := 0
	for split := 0; split < splits; split++ {
		if startIX >= len(data) {
			break
		}
		endIX += b.rand.Intn(len(data))
		if endIX >= len(data) {
			endIX = len(data)
		}

		newFrag := []byte{}
		for ix := startIX; ix < endIX; ix++ {
			newFrag = append(newFrag, data[ix])
		}
		startIX = endIX
		output = append(output, newFrag)
	}

	return output
}
