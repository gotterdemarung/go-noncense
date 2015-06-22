package noncense

import (
	"hash/crc32"
)

// Simple string with hash code
type HString struct {
	Value string
	HashCode uint32
}

// Constructor
func NewHString(value string) HString {

	h := crc32.NewIEEE()
	h.Write([]byte(value))
	return HString{Value: value, HashCode: h.Sum32()}
}

// Returns trimmed to provided size hashcode
// Used to calculate position inside hashmap
func (s *HString) trim(mod uint32) uint32 {
	return s.HashCode % mod;
}