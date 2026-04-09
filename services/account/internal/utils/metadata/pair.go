package metadata

type Pair struct {
	key   string
	value string
}

func NewPair(key, value string) Pair {
	return Pair{key: key, value: value}
}
