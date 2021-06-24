package filters

type BloomFilter interface {
	Set(data []byte) error
	SetMany(data [][]byte) error
	IsInSet(data []byte) bool
}
