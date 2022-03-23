package merkle

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestNode_Sibling(t *testing.T) {
	left := &Node{val: []byte("left")}
	right := &Node{val: []byte("right")}
	root := &Node{
		val:   []byte("root"),
		right: right,
		left:  left,
	}
	left.parent = root
	right.parent = root

	t.Run("Should Return Right", func(t *testing.T) {
		if left.Sibling() != right {
			t.Error("expected opposite to be right")
		}
	})

	t.Run("Should Return Left", func(t *testing.T) {
		if left.Sibling() != right {
			t.Error("expected opposite to be left")
		}
	})

	t.Run("Should Return nil", func(t *testing.T) {
		if root.Sibling() != nil {
			t.Error("expected opposite to be nil")
		}
	})
}

func TestNode_IsLeaf(t *testing.T) {

	leaf := &Node{val: []byte("leaf")}
	inner := &Node{val: []byte("inner"), left: leaf, right: leaf}
	innerNoLeft := &Node{val: []byte("inner"), right: leaf}
	innerNoRight := &Node{val: []byte("inner"), left: leaf}

	t.Run("Should Return True", func(t *testing.T) {
		if !leaf.IsLeaf() {
			t.Errorf("expected to return true")
		}
	})

	t.Run("Should Return False", func(t *testing.T) {
		if inner.IsLeaf() || innerNoLeft.IsLeaf() || innerNoRight.IsLeaf() {
			t.Errorf("expected to return false")
		}
	})
}

func TestNode_IsLeft(t *testing.T) {
	left := &Node{val: []byte("left")}
	right := &Node{val: []byte("right")}
	root := &Node{
		val:   []byte("root"),
		right: right,
		left:  left,
	}
	left.parent = root
	right.parent = root
	if !left.IsLeft() {
		t.Errorf("expected to be a left children")
	}
}

func TestNode_IsRight(t *testing.T) {
	left := &Node{val: []byte("left")}
	right := &Node{val: []byte("right")}
	root := &Node{
		val:   []byte("root"),
		right: right,
		left:  left,
	}
	left.parent = root
	right.parent = root
	if !right.IsRight() {
		t.Errorf("expected to be a right children")
	}
}

func TestNode_WalkPreOrder(t *testing.T) {
	leftLeftChild := &Node{val: []byte("left - child - left")}
	leftRightChild := &Node{val: []byte("left - child - right")}

	rightRightChild := &Node{val: []byte("right - child - right")}
	rightLeftChild := &Node{val: []byte("right - child - left")}

	rootLeftChild := &Node{val: []byte("root - child - left")}
	rootLeftChild.left = leftLeftChild
	rootLeftChild.right = leftRightChild

	rootRightChild := &Node{val: []byte("root - child - right")}
	rootRightChild.left = rightLeftChild
	rootRightChild.right = rightRightChild

	root := &Node{
		val:   []byte("root"),
		right: rootRightChild,
		left:  rootLeftChild,
	}

	// expected depth and order of walk Node
	expOrderWalkWalk := [][]interface{}{
		[]interface{}{0, root},
		[]interface{}{1, root.left},
		[]interface{}{2, root.left.left},
		[]interface{}{2, root.left.right},
		[]interface{}{1, root.right},
		[]interface{}{2, root.right.left},
		[]interface{}{2, root.right.right},
	}

	iteration := 0

	root.WalkPreOrder(func(n *Node, depth int) {
		expDepth := expOrderWalkWalk[iteration][0].(int)
		expNode := expOrderWalkWalk[iteration][1].(*Node)
		if depth != expDepth {
			t.Errorf("expected depth at %d to be %d, got %d", iteration, expDepth, depth)
		}
		if n != expNode {
			t.Errorf("expected node at %d to be %s, got %s", iteration, expNode.val, n.val)
		}
		iteration++
	})
}

func TestNodes_Len(t *testing.T) {
	nodes := Nodes{
		&Node{val: []byte("1")},
		&Node{val: []byte("2")},
		&Node{val: []byte("3")},
		&Node{val: []byte("4")},
		&Node{val: []byte("5")},
	}

	exp := len(nodes)

	if act := nodes.Len(); exp != act {
		t.Errorf("expected Len to be %d, got %d", exp, act)
	}
}

func TestNodes_Less(t *testing.T) {
	nodes := Nodes{
		&Node{val: []byte("a")},
		&Node{val: []byte("b")},
	}
	if !nodes.Less(0, 1) {
		t.Errorf("i should be lower than j")
	}
}

func TestNodes_Swap(t *testing.T) {
	nodes := Nodes{
		&Node{val: []byte("a")},
		&Node{val: []byte("b")},
	}

	expi := nodes[1]
	expj := nodes[0]

	nodes.Swap(0, 1)

	if nodes[0] != expi {
		t.Errorf("i should have been %s", expi.val)
	}

	if nodes[1] != expj {
		t.Errorf("j should have been %s", expi.val)
	}
}

func TestNodes_IteratePair(t *testing.T) {

	nodes := Nodes{
		&Node{val: []byte("1")},
		&Node{val: []byte("2")},
		&Node{val: []byte("3")},
		&Node{val: []byte("4")},
		&Node{val: []byte("5")},
	}

	t.Run("Should Return The Odd Node", func(t *testing.T) {
		odd := nodes.IteratePair(func(i, j *Node) {})
		if odd == nil {
			t.Errorf("expected odd Node to be returned")
		} else if bytes.Compare(nodes[len(nodes)-1].val, odd.val) != 0 {
			t.Errorf("wrong odd node was returned")
		}
	})

	t.Run("Should Return A nil Odd Node", func(t *testing.T) {
		odd := nodes[:len(nodes)-1].IteratePair(func(i, j *Node) {})
		if odd != nil {
			t.Errorf("unexpected odd Node returned")
		}
	})

	t.Run("Should Iterate And Pair Correctly", func(t *testing.T) {
		exp4iters := map[int][]*Node{
			0: []*Node{nodes[0], nodes[1]},
			1: []*Node{nodes[2], nodes[3]},
		}
		iteration := 0
		nodes.IteratePair(func(i, j *Node) {
			acti := i.val
			actj := j.val
			expi := exp4iters[iteration][0].val
			expj := exp4iters[iteration][1].val
			if bytes.Compare(acti, expi) != 0 {
				t.Errorf("expected i[%d] to be %s, got %s", iteration, expi, acti)
			}
			if bytes.Compare(actj, expj) != 0 {
				t.Errorf("expected j[%d] to be %s, got %s", iteration, expj, actj)
			}
			iteration++
		})
	})
}

func TestNodes_IterateSortedPair(t *testing.T) {

	nodes := Nodes{
		&Node{val: []byte("e")},
		&Node{val: []byte("d")},
		&Node{val: []byte("c")},
		&Node{val: []byte("b")},
		&Node{val: []byte("a")},
	}

	t.Run("Should Return The Odd Node", func(t *testing.T) {
		odd := nodes.IterateSortedPair(func(i, j *Node) {})
		if odd == nil {
			t.Errorf("expected odd Node to be returned")
		} else if bytes.Compare(nodes[len(nodes)-1].val, odd.val) != 0 {
			t.Errorf("wrong odd node was returned")
		}
	})

	t.Run("Should Return A nil Odd Node", func(t *testing.T) {
		odd := nodes[:len(nodes)-1].IterateSortedPair(func(i, j *Node) {})
		if odd != nil {
			t.Errorf("unexpected odd Node returned")
		}
	})

	t.Run("Should Iterate And Sort Pair Correctly", func(t *testing.T) {
		exp4iters := map[int][]*Node{
			0: []*Node{nodes[1], nodes[0]},
			1: []*Node{nodes[3], nodes[2]},
		}
		iteration := 0
		nodes.IterateSortedPair(func(i, j *Node) {
			acti := i.val
			actj := j.val
			expi := exp4iters[iteration][0].val
			expj := exp4iters[iteration][1].val
			if bytes.Compare(acti, expi) != 0 {
				t.Errorf("expected i[%d] to be %s, got %s", iteration, expi, acti)
			}
			if bytes.Compare(actj, expj) != 0 {
				t.Errorf("expected j[%d] to be %s, got %s", iteration, expj, actj)
			}
			iteration++
		})
	})
}

func TestNodes_ToHexStrings(t *testing.T) {
	nodes := byteArrSliceToNodes(hashStringSlice(sha256.New(), "a", "b", "c")...)
	expHex := []string{
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d",
		"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
	}
	actHex := nodes.ToHexStrings()
	for i, exp := range expHex {
		if actHex[i] != exp {
			t.Errorf("expected hex at %d to be %s, got %s", i, exp, actHex)
		}
	}
}

func TestNodes_ToByteArrays(t *testing.T) {
	nodes := Nodes{
		&Node{val: []byte("a")},
		&Node{val: []byte("b")},
	}
	for i, val := range nodes.ToByteArrays() {
		exp := string(nodes[i].val)
		act := string(val)
		if exp != act {
			t.Errorf("expected val at index %d to be %s, got %s", i, exp, act)
		}
	}
}
