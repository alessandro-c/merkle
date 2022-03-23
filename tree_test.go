package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"testing"
)

func hashStringSlice(algo hash.Hash, strings ...string) [][]byte {
	hashes := make([][]byte, len(strings))
	for i, s := range strings {
		algo.Reset()
		algo.Write([]byte(s))
		hashes[i] = algo.Sum(nil)
	}
	return hashes
}

func hexStringsToByteArrays(strings ...string) [][]byte {
	barr := make([][]byte, 0, len(strings))
	for _, s := range strings {
		b, _ := hex.DecodeString(s)
		barr = append(barr, b)
	}
	return barr
}

// prepping test data
var algo hash.Hash = sha256.New()

/*
Merkle tree with an odd number of leaves

└─ 3a64c13ffc8d22739538f49d901d909754e4ca185cf128ce7e64c8482f0cd8c6
   ├─ a26df13b366b0fc0e7a96ec9a1658d691d7640668de633333098d7952ce0c50b
   │  ├─ 800e03ddb2432933692401d1631850c0af91953fd9c8f3874488c0541dfcf413
   │  │  ├─ 18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4
   │  │  └─ 2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6
   │  └─ 28b5a66c8c61ee13ad5f708a561d758b24d10abe5a0e72133c85d59821539e05
   │     ├─ 3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d
   │     └─ 3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea
   └─ ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb
*/
var oddLeavesTree = NewTree(algo, hashStringSlice(algo, "a", "b", "c", "d", "e"))

var oddLeavesTreeProofs = map[string][]string{
	"18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4": []string{
		"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
		"28b5a66c8c61ee13ad5f708a561d758b24d10abe5a0e72133c85d59821539e05",
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
	},
	"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6": []string{
		"18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4",
		"28b5a66c8c61ee13ad5f708a561d758b24d10abe5a0e72133c85d59821539e05",
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
	},
	"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d": []string{
		"3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea",
		"800e03ddb2432933692401d1631850c0af91953fd9c8f3874488c0541dfcf413",
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
	},
	"3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea": []string{
		"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d",
		"800e03ddb2432933692401d1631850c0af91953fd9c8f3874488c0541dfcf413",
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
	},
	"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb": []string{
		"a26df13b366b0fc0e7a96ec9a1658d691d7640668de633333098d7952ce0c50b",
	},
}

/*
Merkle tree with an even number of leaves

└─ 4c6aae040ffada3d02598207b8485fcbe161c03f4cb3f660e4d341e7496ff3b2
   ├─ 800e03ddb2432933692401d1631850c0af91953fd9c8f3874488c0541dfcf413
   │  ├─ 18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4
   │  └─ 2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6
   └─ 18d79cb747ea174c59f3a3b41768672526d56fecc58360a99d283d0f9b0a3cc0
      ├─ 3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d
      └─ ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb
*/
var evenLeavesTree = NewTree(algo, hashStringSlice(algo, "a", "b", "c", "d"))

var evenLeavesTreeProofs = map[string][]string{
	"18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4": []string{
		"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
		"18d79cb747ea174c59f3a3b41768672526d56fecc58360a99d283d0f9b0a3cc0",
	},
	"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6": []string{
		"18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4",
		"18d79cb747ea174c59f3a3b41768672526d56fecc58360a99d283d0f9b0a3cc0",
	},
	"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d": []string{
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"800e03ddb2432933692401d1631850c0af91953fd9c8f3874488c0541dfcf413",
	},
	"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb": []string{
		"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d",
		"800e03ddb2432933692401d1631850c0af91953fd9c8f3874488c0541dfcf413",
	},
}

func TestNewTree(t *testing.T) {

	t.Run("With Even Leaves", func(t *testing.T) {
		t.Run("Should Return Expected Merkle Root", func(t *testing.T) {
			exp := "4c6aae040ffada3d02598207b8485fcbe161c03f4cb3f660e4d341e7496ff3b2"
			if act := evenLeavesTree.Root().String(); act != exp {
				t.Errorf("expected merkle root should have been %s, got %s", exp, act)
			}
		})
	})

	t.Run("With Odd Leaves", func(t *testing.T) {
		t.Run("Should Return Expected Merkle Root", func(t *testing.T) {
			exp := "3a64c13ffc8d22739538f49d901d909754e4ca185cf128ce7e64c8482f0cd8c6"
			if act := oddLeavesTree.Root().String(); act != exp {
				t.Errorf("expected merkle root should have been %s, got %s", exp, act)
			}
		})
	})

}

func TestTree_Proof(t *testing.T) {
	t.Run("With Non Existent Leaf", func(t *testing.T) {
		t.Run("Should Return Empty Proof", func(t *testing.T) {
			proof := evenLeavesTree.Proof([]byte("foo"))
			if len(proof) > 0 {
				t.Errorf("expected empty proof")
			}
		})
	})
	t.Run("With Even Leaves", func(t *testing.T) {
		for leaf, expProof := range evenLeavesTreeProofs {
			t.Run("Should Return Expected Proof For Leaf "+leaf, func(t *testing.T) {
				leafb, _ := hex.DecodeString(leaf)
				actProof := evenLeavesTree.Proof(leafb)
				if len(expProof) != len(actProof) {
					t.Errorf("expected length of proof to be %d, got %d", len(expProof), len(actProof))
					t.SkipNow()
				}
				actProofHex := actProof.ToHexStrings()
				for i := 0; i < len(actProofHex); i++ {
					if actProofHex[i] != expProof[i] {
						t.Errorf("expected node at index %d to be %s, got %s", i, expProof[i], actProofHex[i])
					}
				}
			})
		}
	})
	t.Run("With Odd Leaves", func(t *testing.T) {
		for leaf, expProof := range oddLeavesTreeProofs {
			t.Run("Should Return Expected Proof For Leaf "+leaf, func(t *testing.T) {
				leafb, _ := hex.DecodeString(leaf)
				actProof := oddLeavesTree.Proof(leafb)
				if len(expProof) != len(actProof) {
					t.Errorf("expected length of proof to be %d, got %d", len(expProof), len(actProof))
					t.SkipNow()
				}
				actProofHex := actProof.ToHexStrings()
				for i := 0; i < len(actProofHex); i++ {
					if actProofHex[i] != expProof[i] {
						t.Errorf("expected node at index %d to be %s, got %s", i, expProof[i], actProofHex[i])
					}
				}
			})
		}
	})
}

func TestVerify(t *testing.T) {
	t.Run("For Even Leaves Tree", func(t *testing.T) {
		for leaf, proof := range evenLeavesTreeProofs {
			t.Run("Should Be Verified For "+leaf, func(t *testing.T) {
				leafb, _ := hex.DecodeString(leaf)
				ok := Verify(algo, leafb, evenLeavesTree.root.val, hexStringsToByteArrays(proof...))
				if !ok {
					t.Errorf("proof should have been valid")
				}
			})
		}
	})
	t.Run("For Odd Leaves Tree", func(t *testing.T) {
		for leaf, proof := range oddLeavesTreeProofs {
			t.Run("Should Be Verified For "+leaf, func(t *testing.T) {
				leafb, _ := hex.DecodeString(leaf)
				ok := Verify(algo, leafb, oddLeavesTree.root.val, hexStringsToByteArrays(proof...))
				if !ok {
					t.Errorf("proof should have been valid")
				}
			})
		}
	})
}
