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
	return nil, nil
}

func (b *BogDB) corrode(fragment []byte, intensity int64) []byte {
	// TODO
	return fragment
}

func (b *BogDB) fragment(fragment []byte) [][]byte {
	// TODO
	return [][]byte{fragment}
}
