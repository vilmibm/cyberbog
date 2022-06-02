package bogdb

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

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

func timeToIntensity(t time.Time) int64 {
	// TODO
	return 0
}

// Inter breaks up data into arbitrarily corroded fragments and stores them on the filesystem.
func (b *BogDB) Inter(data []byte) error {
	fragments := b.fragment(data)
	var err error
	for _, fragment := range fragments {
		now := b.now()
		fragName := fmt.Sprintf("%d", b.rand.Uint64())
		fragDir := path.Join(
			b.rootPath,
			string(fragName[0]), string(fragName[1]), string(fragName[2]))
		err = os.MkdirAll(fragDir, 0750)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("could not create '%s': %w", fragDir, err)
		}

		fragPath := path.Join(fragDir, string(fragName))
		f, err := os.Create(fragPath)
		if err != nil {
			return fmt.Errorf("failed to open fragment for writing '%s': %w", fragPath, err)
		}

		fmt.Fprintf(f, "%s\n", now)
		_, err = f.Write(b.corrode(fragment, timeToIntensity(now)))
		if err != nil {
			return fmt.Errorf("failed to write fragment '%s': %w", fragPath, err)
		}
	}

	return nil
}

// Exhume reads a fragment from disk, further fragmenting and corroding it. One fragment is returned, the rest re-interred. Fragments are read, deleted from disk, then re-written as needed. This absolutely can lead to lost data; this is by design.
func (b *BogDB) Exhume() ([]byte, error) {
	// TODO
	dir := b.rootPath
	deeper := true
	var fragment []byte
	for deeper {
		var err error
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read dir '%s': %w", dir, err)
		}
		if len(entries) == 0 {
			// TODO should this be error?
			return []byte{}, nil
		}
		selected := b.rand.Intn(len(entries))
		for ix, entry := range entries {
			if ix != selected {
				continue
			}
			absPath := path.Join(dir, entry.Name())
			if entry.IsDir() {
				dir = absPath
				continue
			}
			fragment, err = os.ReadFile(absPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file '%s': %w", absPath, err)
			}
			deeper = false
			err = os.Remove(absPath)
			if err != nil {
				return fragment, fmt.Errorf("failed to remove fragment '%s': %w", absPath, err)
			}
		}
	}

	// TODO should we fragment and put some pieces back in?

	// TODO should we strip the date? I think no, that should be up to the caller of bogdb

	return fragment, nil
}

func (b *BogDB) corrode(fragment []byte, intensity int64) []byte {
	// TODO
	return fragment
}

func (b *BogDB) fragment(fragment []byte) [][]byte {
	output := [][]byte{}
	// TODO i want to divide fragment length by 10 but Go is being recalcitrant
	splits := b.rand.Intn(len(fragment)) + 1

	startIX := 0
	endIX := 0
	for split := 0; split < splits; split++ {
		if startIX >= len(fragment) {
			break
		}
		endIX += b.rand.Intn(len(fragment))
		if endIX >= len(fragment) {
			endIX = len(fragment)
		}

		newFrag := []byte{}
		for ix := startIX; ix < endIX; ix++ {
			newFrag = append(newFrag, fragment[ix])
		}
		startIX = endIX
		output = append(output, newFrag)
	}

	return output
}
