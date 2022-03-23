package main

import (
	"crypto/sha256"
	"github.com/alessandro-c/merkle"
	"hash"
	"log"
	"math/rand"
	"time"
)

func main() {
	algo := sha256.New()

	leaves := [][]byte{
		hashString(algo, "a"), hashString(algo, "b"),
		hashString(algo, "c"), hashString(algo, "d"),
		hashString(algo, "e"),
	}

	// you can change the order of leaves without affecting the end result, that is, the same merkle root
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(leaves), func(i, j int) {
		leaves[i], leaves[j] = leaves[j], leaves[i]
	})

	for i, l := range leaves {
		log.Printf("hex leaf #%d - %x\n", i, l)
	}

	// building up tree up to the merkle root
	tree := merkle.NewTree(algo, leaves)

	// merkle root
	log.Println("hex merkle root: ", tree.Root().Hex())

	// building proof for leaf c
	hashedLeafToProof := hashString(algo, "c")
	proof := tree.Proof(hashedLeafToProof)

	for i, h := range proof.ToHexStrings() {
		log.Printf("proof for leaf %x at index %d is : %s", hashedLeafToProof, i, h)
	}

	// verifying proof
	ok := merkle.Verify(algo, hashString(algo, "c"), tree.Root().Bytes(), proof.ToByteArrays())
	log.Println("proof is valid ?", ok)

	// printing whole tree
	log.Println("whole tree below : ")
	tree.Root().Graphify(log.Writer())
}

func hashString(algo hash.Hash, s string) []byte {
	algo.Reset()
	algo.Write([]byte(s))
	return algo.Sum(nil)
}
