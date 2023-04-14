package bloom

type DefaultHash struct {
}

func (*DefaultHash) Hash(s string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		hash = (hash * 16777619) ^ uint32(s[i])
	}
	return hash
}
