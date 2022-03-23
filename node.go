package merkle

import (
	"bytes"
	"fmt"
	"github.com/xlab/treeprint"
	"io"
)

// Node is a merkle tree node
type Node struct {
	// val is the node hashed value
	val    []byte
	left   *Node
	right  *Node
	parent *Node
}

// Bytes return the raw hash
func (n Node) Bytes() []byte {
	return n.val
}

// Hex returns the Node val represented as an hexadecimal string
func (n Node) Hex() string {
	return fmt.Sprintf("%x", n.val)
}

// String implements most common interfaces
func (n Node) String() string {
	// using fmt.Sprintf to allow certain IDE debuggers
	// displaying the value during debug mode, realistically
	// this is not needed but worth it for debugging purposes.
	// See: https://youtrack.jetbrains.com/issue/GO-7821
	return fmt.Sprintf("%x", n.val)
}

// IsLeaf tells whether the Node is a leaf type of node
func (n *Node) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

// IsLeft tells whether the Node is a left child of its parent
func (n *Node) IsLeft() bool {
	return n.parent != nil && n.parent.left == n
}

// IsRight tells whether the Node is a right child of its parent
func (n *Node) IsRight() bool {
	return n.parent != nil && n.parent.right == n
}

// Sibling returns its opposite sibling.
// Given 2 nodes i, j if Node is i returns j else returns i.
// Returns nil if root.
func (n *Node) Sibling() *Node {
	if n.parent == nil {
		return nil
	}
	if n.parent.left == n {
		return n.parent.right
	}
	return n.parent.left
}

// Graphify builds up a hierarchical graphic representation
// from the Node to the very bottom of its children.
// Will write to the provided io.Writer for greater usability.
//
// For example, to print in your terminal you may do :
//
//  n.Graphify(os.Stdout)
//
// where n is the Node instance you want to print from
func (n *Node) Graphify(w io.Writer) {

	branches := map[string]treeprint.Tree{
		n.Hex(): treeprint.NewWithRoot(n.Hex()),
	}

	// this has its limitations as it assumes there won't be
	// any duplicate hash in the tree.
	n.WalkPreOrder(func(n *Node, depth int) {
		if n.IsLeaf() {
			branches[n.parent.Hex()].AddNode(n.Hex())
		} else if _, ok := branches[n.Hex()]; !ok {
			branches[n.Hex()] = branches[n.parent.Hex()].AddBranch(n.Hex())
		}
	})

	w.Write(branches[n.Hex()].Bytes())
}

// WalkPreOrder traverses from the tree *Node down
// to the very bottom using the "Pre Order" strategy.
func (n *Node) WalkPreOrder(fn func(n *Node, depth int)) {
	var por func(n *Node, depth int, fn func(n *Node, depth int))
	por = func(n *Node, depth int, fn func(n *Node, depth int)) {
		if n != nil {
			fn(n, depth)
			depth++
			por(n.left, depth, fn)
			por(n.right, depth, fn)
		}
	}
	por(n, 0, fn)
}

// Nodes is slice type of *Node
type Nodes []*Node

// Len implements the sort.Interface
func (n Nodes) Len() int {
	return len(n)
}

// Less implements the sort.Interface
func (n Nodes) Less(i, j int) bool {
	return bytes.Compare(n[i].val, n[j].val) == -1
}

// Swap implements the sort.Interface
func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// IteratePair iterates through all Nodes pairing with fn(i,j).
// If there is an odd number Nodes the last element Node len(n) - 1 will be returned.
func (n Nodes) IteratePair(fn func(i, j *Node)) (odd *Node) {
	if len(n)%2 != 0 {
		odd = n[len(n)-1]
	}
	for i := 0; i < len(n)-1; i = i + 2 {
		fn(n[i], n[i+1])
	}
	return
}

// IterateSortedPair iterate same as IteratePair but with sorted ascending i,j
func (n Nodes) IterateSortedPair(fn func(i, j *Node)) (odd *Node) {
	odd = n.IteratePair(func(i, j *Node) {
		if bytes.Compare(i.val, j.val) == 1 {
			// i > j
			fn(j, i)
			return
		}
		fn(i, j)
	})
	return
}

// ToHexStrings converts each Node in Nodes into an hex strings.
func (ns Nodes) ToHexStrings() []string {
	hexs := make([]string, 0, len(ns))
	for _, n := range ns {
		hexs = append(hexs, n.Hex())
	}
	return hexs
}

// ToByteArrays converts each Node in Nodes into a slice of byte array
func (ns Nodes) ToByteArrays() [][]byte {
	barr := make([][]byte, 0, len(ns))
	for _, n := range ns {
		barr = append(barr, n.val)
	}
	return barr
}

// newNode makes and return a new *Node
// with the provided hash set as val
func newNode(h []byte) *Node {
	return &Node{val: h}
}

// newParentNode makes and return a new *Node
// with the provided hash set as val.
// The l (left) and r (right) will be associated as children.
func newParentNode(h []byte, l, r *Node) *Node {
	n := newNode(h)
	n.left = l
	n.right = r
	return n
}

// byteArrSliceToNodes turns the byte array slice into Nodes
func byteArrSliceToNodes(bas ...[]byte) Nodes {
	nodes := make(Nodes, len(bas))
	for i := 0; i < len(bas); i++ {
		nodes[i] = newNode(bas[i])
	}
	return nodes
}
