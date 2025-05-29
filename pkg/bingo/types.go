package bingo

const (
	Regular  = 1
	AllHands = 2
)

func TypeExists(t int64) bool {
	switch t {
	case Regular, AllHands:
		return true
	}
	return false
}
