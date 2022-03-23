// Package merkle
//
// A merkle tree is a kind of binary tree where each node is labelled with a hash. Starting from the very bottom,
// leaves will be paired and hashed together to make their parent inner-node, recursively up to the root a.k.a merkle root.
//
// Merkle trees, are commonly used in distributed systems to efficiently compare large data set ensuring validity of such data.
//
// Some examples that leverage merkle trees are the Git version control, AWS's QLDB, Apache's Cassandra and last but not least, blockchains.
//
// There are different flavours of implementations, this package doesn't attempt to build an abstraction for all the
// possible ones out there, it rather implements a fairly efficient specific one that can be used to experiment with the data structure and concepts.
//
// With that said, if you're looking to use this package to validate merkle proofs for existing blockchains you should look
// elsewhere as their implementation may be different. For example, Bitcoin's merkle, duplicates eventual odd nodes to re-balance
// the tree and this implementation doesn't, thus producing a different merkle root and proof.
package merkle

import (
	"bytes"
	"hash"
	"sort"
)

// Tree is a whole merkle tree
type Tree struct {
	// the merkle root Node
	root *Node
	// stored for convenience to avoid traversing
	leaves Nodes
}

// NewTree builds up a new merkle tree with the provided
// hashing algorithm and set of leaves that have been
// hashed with the same algorithm.
func NewTree(h hash.Hash, hl [][]byte) *Tree {
	// turning leaves into nodes
	leaves := byteArrSliceToNodes(hl...)
	// sorting leaves lexicographically this will come
	// in handy to efficiently build proofs and find leaves
	sort.Sort(leaves)
	// building up tree up to root
	root := buildTree(h, leaves)
	return &Tree{root, leaves}
}

// Root returns the root *Node a.k.a merkle root
func (t Tree) Root() *Node {
	return t.root
}

func buildTree(h hash.Hash, n Nodes) *Node {

	// allocating with just enough capacity.
	// +1 to give space for eventual odd to re-balance
	ps := make(Nodes, 0, len(n)/2+1)

	// pairing sorted nodes and making parents hashing pairs.
	// if an odd number of nodes was provided the last
	// item will be removed and will be re-used later to re-balance
	odd := n.IterateSortedPair(func(i, j *Node) {
		// hashing paired nodes
		h.Reset()
		h.Write(i.val)
		h.Write(j.val)
		// making parent node from hashed pair
		p := newParentNode(h.Sum(nil), i, j)
		// attaching parent node
		i.parent = p
		j.parent = p
		// appending parent for next batch of recursive iteration
		ps = append(ps, p)
	})

	// if there is an odd push it back to re-balance
	if odd != nil {
		ps = append(ps, odd)
	}

	// recursively building up tree
	// until we have only one node (aka merkle root)
	if len(ps) > 1 {
		return buildTree(h, ps)
	}

	// merkle root reached
	return ps[0]
}

// Proof builds and returns the merkle proof for the provided hashed leaf.
func (t Tree) Proof(hl []byte) Nodes {

	// at first, let's find out whether the leaf actually
	// exists. Given that the leaves were originally sorted
	// we can use binary search to efficiently find the leaf.
	ihl := sort.Search(len(t.leaves), func(i int) bool {
		cmp := bytes.Compare(t.leaves[i].val, hl)
		return cmp == 1 || cmp == 0 // t.leaves[i].val >= hl
	})

	// checking whether the leaf was actually found, if not
	// we will just simply return an empty slice of Nodes
	if ihl >= len(t.leaves) || bytes.Compare(t.leaves[ihl].val, hl) != 0 {
		return Nodes{}
	}

	// allocating just enough capacity leaving
	// enough space for an eventual odd as well
	proof := make(Nodes, 0, len(t.leaves)/2)
	var buildProof func(n *Node)
	buildProof = func(n *Node) {
		if n != t.root {
			proof = append(proof, n.Sibling())
			buildProof(n.parent)
		}
	}
	buildProof(t.leaves[ihl])

	return proof
}

// Verify verifies whether the provided proof for leaf is valid.
func Verify(algo hash.Hash, leaf, root []byte, proof [][]byte) bool {
	for _, h := range proof {
		// leaf is a left child node
		i, j := leaf, h
		if cmp := bytes.Compare(leaf, h); cmp == 1 {
			// leaf is a right child node
			i, j = h, leaf
		}
		algo.Reset()
		algo.Write(i)
		algo.Write(j)
		leaf = algo.Sum(nil)
	}
	return bytes.Compare(leaf, root) == 0
}
