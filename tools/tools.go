package tools

import (
	"bytes"
)

func RemoveIntFromSlice(slice []int, cell int) []int {
	indexToRemove := 0
	for i, c := range slice {
		if c == cell {
			indexToRemove = i
		}
	}
	return append(slice[:indexToRemove], slice[indexToRemove+1:]...)
}

func EncodeChar(ch int) uint8 {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	return alphabet[ch]
}

func DecodeChar(c uint8) uint8 {
	if c >= 'a' && c <= 'z' {
		return c - 'a'
	}
	if c >= 'A' && c <= 'Z' {
		return (c - 'A') + 26
	}
	if c >= '0' && c <= '9' {
		return (c - '0') + 26 * 2
	}
	if c == '-' {
		return 62
	}
	if c == '_' {
		return 63
	}
	return 0
}

func ExtractPackets(b* []byte) [][]byte {
	var packets [][]byte
	current := bytes.NewBuffer([]byte{})

	by := bytes.NewBuffer(*b)

	for {
		bi, err := by.ReadByte()
		if err != nil {
			return packets
		}
		current.WriteByte(bi)
		if bi == 0 {
			packets = append(packets, current.Bytes())
			current = bytes.NewBuffer([]byte{})
		}
	}
}
