package abstraction

type HashFunction interface {
	Hash(s string) uint32
}
