package log

import (
	"crypto/sha256"
	"encoding/hex"
)

func hashPair(a string, b string) string {
	h := sha256.Sum256([]byte(a + b))
	return hex.EncodeToString(h[:])
}

func ComputeMerkleRoot(leaves []string) string {

	if len(leaves) == 0 {
		return ""
	}

	nodes := leaves

	for len(nodes) > 1 {

		var next []string

		for i := 0; i < len(nodes); i += 2 {

			if i+1 < len(nodes) {
				next = append(next, hashPair(nodes[i], nodes[i+1]))
			} else {
				next = append(next, hashPair(nodes[i], nodes[i]))
			}

		}

		nodes = next
	}

	return nodes[0]
}

type MerkleProof struct {
	Leaf   string
	Path   []string
	Index  int
	Root   string
}

func GenerateProof(leaves []string, index int) MerkleProof {

	var path []string
	nodes := leaves
	i := index

	for len(nodes) > 1 {

		var next []string

		for j := 0; j < len(nodes); j += 2 {

			left := nodes[j]
			right := left

			if j+1 < len(nodes) {
				right = nodes[j+1]
			}

			next = append(next, hashPair(left, right))

			if j == i || j+1 == i {

				if j == i {
					path = append(path, right)
				} else {
					path = append(path, left)
				}

				i = len(next) - 1
			}
		}

		nodes = next
	}

	return MerkleProof{
		Leaf:  leaves[index],
		Path:  path,
		Index: index,
		Root:  nodes[0],
	}
}

func VerifyProof(proof MerkleProof) bool {

	hash := proof.Leaf
	index := proof.Index

	for _, sibling := range proof.Path {

		if index%2 == 0 {
			hash = hashPair(hash, sibling)
		} else {
			hash = hashPair(sibling, hash)
		}

		index = index / 2
	}

	return hash == proof.Root
}
