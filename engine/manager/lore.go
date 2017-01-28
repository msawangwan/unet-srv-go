package manager

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/msawangwan/unet-srv-go/engine/prng"
)

type Lore struct {
	Description []string `json:"description"`
}

type LoreGenerator struct {
	cache map[string]*Lore

	*prng.Instance

	sync.Mutex
}

func NewLoreGenerator(fname string, rand *prng.Instance) (*LoreGenerator, error) {
	if !strings.HasSuffix(fname, ".json") {
		return nil, errors.New("not a json file (info text file)")
	}

	f, e := os.Open(fname)
	if e != nil {
		return nil, e
	}

	var (
		l *Lore
		r *prng.Instance
	)

	e = json.NewDecoder(f).Decode(&l)
	if e != nil {
		return nil, e
	}

	if rand == nil {
		r = prng.New(0)
	} else {
		r = rand
	}

	lg := &LoreGenerator{
		cache:    make(map[string]*Lore),
		Instance: r,
	}

	lg.Lock()
	defer lg.Unlock()

	lg.cache["singles"] = l

	return lg, nil
}

func (lg *LoreGenerator) Random() (*string, error) {
	lg.Lock()
	defer lg.Unlock()

	singles := lg.cache["singles"]

	desc := singles.Description[lg.InRangeInt(0, len(singles.Description))]

	return &desc, nil
}
