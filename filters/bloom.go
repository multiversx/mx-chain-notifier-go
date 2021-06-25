package filters

import (
	"hash"
	"hash/fnv"
	"math"

	"github.com/spaolacci/murmur3"
)

const (
	falsePositiveProb = 0.01
	setBit            = true
)

type Bloom struct {
	m      uint
	n      uint
	k      uint
	h1     hash.Hash64
	h2     hash.Hash64
	bitset []bool
}

func NewBloom(n uint) *Bloom {
	if n == 0 {
		n = 1
	}
	m, k := estimateParams(n, falsePositiveProb)

	return &Bloom{
		m:      m,
		n:      n,
		k:      k,
		h1:     fnv.New64a(),
		h2:     murmur3.New64(),
		bitset: make([]bool, m),
	}
}

func (b *Bloom) Set(data []byte) error {
	for i := 0; i < int(b.k); i++ {
		pos, err := b.doubleHash(data, i)
		if err != nil {
			return err
		}
		b.bitset[pos] = setBit
	}
	return nil
}

func (b *Bloom) SetMany(data [][]byte) error {
	for _, item := range data {
		err := b.Set(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bloom) IsInSet(data []byte) bool {
	for i := 0; i < int(b.k); i++ {
		pos, err := b.doubleHash(data, i)
		if err != nil {
			return false
		}
		if !b.bitset[pos] {
			return false
		}
	}
	return true
}

func (b *Bloom) doubleHash(data []byte, i int) (uint64, error) {
	b.h1.Reset()
	b.h2.Reset()

	_, err := b.h1.Write(data)
	if err != nil {
		return 0, err
	}
	_, err = b.h2.Write(data)
	if err != nil {
		return 0, err
	}

	hSum := b.h1.Sum64() + (uint64(i) * b.h2.Sum64())
	g := hSum + uint64(math.Pow(float64(i), 2))
	return g % uint64(b.m), nil
}

func estimateParams(n uint, p float64) (uint, uint) {
	ln2 := math.Log(2)
	ln2sqrt := math.Pow(math.Log(2), 2)

	nlogp := float64(n) * math.Log(p)
	bitsetLen := uint(math.Ceil(-(nlogp) / ln2sqrt))
	numHashes := uint(math.Ceil(ln2 * float64(bitsetLen/n)))

	return bitsetLen, numHashes
}
