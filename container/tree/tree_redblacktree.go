package tree

import (
	"fmt"
	"utils/container/var"
	"utils/utils/json"

	"utils/base/util"
	"utils/convert/conv"

	"utils/utils/rwmutex"
)

type color bool

const (
	black, red color = true, false
)

type RedBlackTree struct {
	mu         rwmutex.RWMutex
	root       *RedBlackTreeNode
	size       int
	comparator func(v1, v2 interface{}) int
}

type RedBlackTreeNode struct {
	Key    interface{}
	Value  interface{}
	color  color
	left   *RedBlackTreeNode
	right  *RedBlackTreeNode
	parent *RedBlackTreeNode
}

func NewRedBlackTree(comparator func(v1, v2 interface{}) int, safe ...bool) *RedBlackTree {
	return &RedBlackTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

func NewRedBlackTreeFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *RedBlackTree {
	tree := NewRedBlackTree(comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

func (tree *RedBlackTree) SetComparator(comparator func(a, b interface{}) int) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.comparator = comparator
	if tree.size > 0 {
		data := make(map[interface{}]interface{}, tree.size)
		tree.doIteratorAsc(tree.leftNode(), func(key, value interface{}) bool {
			data[key] = value
			return true
		})
		tree.root = nil
		tree.size = 0
		for k, v := range data {
			tree.doSet(k, v)
		}
	}
}

func (tree *RedBlackTree) Clone() *RedBlackTree {
	newTree := NewRedBlackTree(tree.comparator, !tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

func (tree *RedBlackTree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

func (tree *RedBlackTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

func (tree *RedBlackTree) doSet(key interface{}, value interface{}) {
	insertedNode := (*RedBlackTreeNode)(nil)
	if tree.root == nil {
		tree.getComparator()(key, key)
		tree.root = &RedBlackTreeNode{Key: key, Value: value, color: red}
		insertedNode = tree.root
	} else {
		node := tree.root
		loop := true
		for loop {
			compare := tree.getComparator()(key, node.Key)
			switch {
			case compare == 0:
				//node.Key   = key
				node.Value = value
				return
			case compare < 0:
				if node.left == nil {
					node.left = &RedBlackTreeNode{Key: key, Value: value, color: red}
					insertedNode = node.left
					loop = false
				} else {
					node = node.left
				}
			case compare > 0:
				if node.right == nil {
					node.right = &RedBlackTreeNode{Key: key, Value: value, color: red}
					insertedNode = node.right
					loop = false
				} else {
					node = node.right
				}
			}
		}
		insertedNode.parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++
}

func (tree *RedBlackTree) Get(key interface{}) (value interface{}) {
	value, _ = tree.Search(key)
	return
}

func (tree *RedBlackTree) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.Value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		tree.doSet(key, value)
	}
	return value
}

func (tree *RedBlackTree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (tree *RedBlackTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (tree *RedBlackTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (tree *RedBlackTree) GetVar(key interface{}) *vvar.Var {
	return vvar.New(tree.Get(key))
}

func (tree *RedBlackTree) GetVarOrSet(key interface{}, value interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSet(key, value))
}

func (tree *RedBlackTree) GetVarOrSetFunc(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSetFunc(key, f))
}

func (tree *RedBlackTree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSetFuncLock(key, f))
}

func (tree *RedBlackTree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (tree *RedBlackTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (tree *RedBlackTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (tree *RedBlackTree) Contains(key interface{}) bool {
	_, ok := tree.Search(key)
	return ok
}

func (tree *RedBlackTree) doRemove(key interface{}) (value interface{}) {
	child := (*RedBlackTreeNode)(nil)
	node, found := tree.doSearch(key)
	if !found {
		return
	}
	value = node.Value
	if node.left != nil && node.right != nil {
		p := node.left.maximumNode()
		node.Key = p.Key
		node.Value = p.Value
		node = p
	}
	if node.left == nil || node.right == nil {
		if node.right == nil {
			child = node.left
		} else {
			child = node.right
		}
		if node.color == black {
			node.color = tree.nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.parent == nil && child != nil {
			child.color = black
		}
	}
	tree.size--
	return
}

func (tree *RedBlackTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

func (tree *RedBlackTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

func (tree *RedBlackTree) IsEmpty() bool {
	return tree.Size() == 0
}

func (tree *RedBlackTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

func (tree *RedBlackTree) Keys() []interface{} {
	var (
		keys  = make([]interface{}, tree.Size())
		index = 0
	)
	tree.IteratorAsc(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

func (tree *RedBlackTree) Values() []interface{} {
	var (
		values = make([]interface{}, tree.Size())
		index  = 0
	)
	tree.IteratorAsc(func(key, value interface{}) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

func (tree *RedBlackTree) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func (tree *RedBlackTree) MapStrAny() map[string]interface{} {
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

func (tree *RedBlackTree) Left() *RedBlackTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.leftNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

func (tree *RedBlackTree) Right() *RedBlackTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.rightNode()
	if tree.mu.IsSafe() {
		return &RedBlackTreeNode{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

func (tree *RedBlackTree) leftNode() *RedBlackTreeNode {
	p := (*RedBlackTreeNode)(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.left
	}
	return p
}

func (tree *RedBlackTree) rightNode() *RedBlackTreeNode {
	p := (*RedBlackTreeNode)(nil)
	n := tree.root
	for n != nil {
		p = n
		n = n.right
	}
	return p
}

func (tree *RedBlackTree) Floor(key interface{}) (floor *RedBlackTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		compare := tree.getComparator()(key, n.Key)
		switch {
		case compare == 0:
			return n, true
		case compare < 0:
			n = n.left
		case compare > 0:
			floor, found = n, true
			n = n.right
		}
	}
	if found {
		return
	}
	return nil, false
}

func (tree *RedBlackTree) Ceiling(key interface{}) (ceiling *RedBlackTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		compare := tree.getComparator()(key, n.Key)
		switch {
		case compare == 0:
			return n, true
		case compare > 0:
			n = n.right
		case compare < 0:
			ceiling, found = n, true
			n = n.left
		}
	}
	if found {
		return
	}
	return nil, false
}

func (tree *RedBlackTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

func (tree *RedBlackTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

func (tree *RedBlackTree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.leftNode(), f)
}

func (tree *RedBlackTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if match {
		if found {
			tree.doIteratorAsc(node, f)
		}
	} else {
		tree.doIteratorAsc(node, f)
	}
}

func (tree *RedBlackTree) doIteratorAsc(node *RedBlackTreeNode, f func(key, value interface{}) bool) {
loop:
	if node == nil {
		return
	}
	if !f(node.Key, node.Value) {
		return
	}
	if node.right != nil {
		node = node.right
		for node.left != nil {
			node = node.left
		}
		goto loop
	}
	if node.parent != nil {
		old := node
		for node.parent != nil {
			node = node.parent
			if tree.getComparator()(old.Key, node.Key) <= 0 {
				goto loop
			}
		}
	}
}

func (tree *RedBlackTree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.rightNode(), f)
}

func (tree *RedBlackTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if match {
		if found {
			tree.doIteratorDesc(node, f)
		}
	} else {
		tree.doIteratorDesc(node, f)
	}
}

func (tree *RedBlackTree) doIteratorDesc(node *RedBlackTreeNode, f func(key, value interface{}) bool) {
loop:
	if node == nil {
		return
	}
	if !f(node.Key, node.Value) {
		return
	}
	if node.left != nil {
		node = node.left
		for node.right != nil {
			node = node.right
		}
		goto loop
	}
	if node.parent != nil {
		old := node
		for node.parent != nil {
			node = node.parent
			if tree.getComparator()(old.Key, node.Key) >= 0 {
				goto loop
			}
		}
	}
}

func (tree *RedBlackTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

func (tree *RedBlackTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for k, v := range data {
		tree.doSet(k, v)
	}
}

func (tree *RedBlackTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
	if tree.size != 0 {
		tree.output(tree.root, "", true, &str)
	}
	return str
}

func (tree *RedBlackTree) Print() {
	fmt.Println(tree.String())
}

func (tree *RedBlackTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, found := tree.doSearch(key)
	if found {
		return node.Value, true
	}
	return nil, false
}

func (tree *RedBlackTree) Flip(comparator ...func(v1, v2 interface{}) int) {
	t := (*RedBlackTree)(nil)
	if len(comparator) > 0 {
		t = NewRedBlackTree(comparator[0], !tree.mu.IsSafe())
	} else {
		t = NewRedBlackTree(tree.comparator, !tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value interface{}) bool {
		t.doSet(value, key)
		return true
	})
	tree.mu.Lock()
	tree.root = t.root
	tree.size = t.size
	tree.mu.Unlock()
}

func (tree *RedBlackTree) output(node *RedBlackTreeNode, prefix string, isTail bool, str *string) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		tree.output(node.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += fmt.Sprintf("%v\n", node.Key)
	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		tree.output(node.left, newPrefix, true, str)
	}
}

func (tree *RedBlackTree) doSearch(key interface{}) (node *RedBlackTreeNode, found bool) {
	node = tree.root
	for node != nil {
		compare := tree.getComparator()(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return node, false
}

func (node *RedBlackTreeNode) grandparent() *RedBlackTreeNode {
	if node != nil && node.parent != nil {
		return node.parent.parent
	}
	return nil
}

func (node *RedBlackTreeNode) uncle() *RedBlackTreeNode {
	if node == nil || node.parent == nil || node.parent.parent == nil {
		return nil
	}
	return node.parent.sibling()
}

func (node *RedBlackTreeNode) sibling() *RedBlackTreeNode {
	if node == nil || node.parent == nil {
		return nil
	}
	if node == node.parent.left {
		return node.parent.right
	}
	return node.parent.left
}

func (tree *RedBlackTree) rotateLeft(node *RedBlackTreeNode) {
	right := node.right
	tree.replaceNode(node, right)
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.left = node
	node.parent = right
}

func (tree *RedBlackTree) rotateRight(node *RedBlackTreeNode) {
	left := node.left
	tree.replaceNode(node, left)
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.right = node
	node.parent = left
}

func (tree *RedBlackTree) replaceNode(old *RedBlackTreeNode, new *RedBlackTreeNode) {
	if old.parent == nil {
		tree.root = new
	} else {
		if old == old.parent.left {
			old.parent.left = new
		} else {
			old.parent.right = new
		}
	}
	if new != nil {
		new.parent = old.parent
	}
}

func (tree *RedBlackTree) insertCase1(node *RedBlackTreeNode) {
	if node.parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *RedBlackTree) insertCase2(node *RedBlackTreeNode) {
	if tree.nodeColor(node.parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *RedBlackTree) insertCase3(node *RedBlackTreeNode) {
	uncle := node.uncle()
	if tree.nodeColor(uncle) == red {
		node.parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *RedBlackTree) insertCase4(node *RedBlackTreeNode) {
	grandparent := node.grandparent()
	if node == node.parent.right && node.parent == grandparent.left {
		tree.rotateLeft(node.parent)
		node = node.left
	} else if node == node.parent.left && node.parent == grandparent.right {
		tree.rotateRight(node.parent)
		node = node.right
	}
	tree.insertCase5(node)
}

func (tree *RedBlackTree) insertCase5(node *RedBlackTreeNode) {
	node.parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.parent.left && node.parent == grandparent.left {
		tree.rotateRight(grandparent)
	} else if node == node.parent.right && node.parent == grandparent.right {
		tree.rotateLeft(grandparent)
	}
}

func (node *RedBlackTreeNode) maximumNode() *RedBlackTreeNode {
	if node == nil {
		return nil
	}
	for node.right != nil {
		return node.right
	}
	return node
}

func (tree *RedBlackTree) deleteCase1(node *RedBlackTreeNode) {
	if node.parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *RedBlackTree) deleteCase2(node *RedBlackTreeNode) {
	sibling := node.sibling()
	if tree.nodeColor(sibling) == red {
		node.parent.color = red
		sibling.color = black
		if node == node.parent.left {
			tree.rotateLeft(node.parent)
		} else {
			tree.rotateRight(node.parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *RedBlackTree) deleteCase3(node *RedBlackTreeNode) {
	sibling := node.sibling()
	if tree.nodeColor(node.parent) == black &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == black &&
		tree.nodeColor(sibling.right) == black {
		sibling.color = red
		tree.deleteCase1(node.parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *RedBlackTree) deleteCase4(node *RedBlackTreeNode) {
	sibling := node.sibling()
	if tree.nodeColor(node.parent) == red &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == black &&
		tree.nodeColor(sibling.right) == black {
		sibling.color = red
		node.parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *RedBlackTree) deleteCase5(node *RedBlackTreeNode) {
	sibling := node.sibling()
	if node == node.parent.left &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.left) == red &&
		tree.nodeColor(sibling.right) == black {
		sibling.color = red
		sibling.left.color = black
		tree.rotateRight(sibling)
	} else if node == node.parent.right &&
		tree.nodeColor(sibling) == black &&
		tree.nodeColor(sibling.right) == red &&
		tree.nodeColor(sibling.left) == black {
		sibling.color = red
		sibling.right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *RedBlackTree) deleteCase6(node *RedBlackTreeNode) {
	sibling := node.sibling()
	sibling.color = tree.nodeColor(node.parent)
	node.parent.color = black
	if node == node.parent.left && tree.nodeColor(sibling.right) == red {
		sibling.right.color = black
		tree.rotateLeft(node.parent)
	} else if tree.nodeColor(sibling.left) == red {
		sibling.left.color = black
		tree.rotateRight(node.parent)
	}
}

func (tree *RedBlackTree) nodeColor(node *RedBlackTreeNode) color {
	if node == nil {
		return black
	}
	return node.color
}

func (tree *RedBlackTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(conv.Map(tree.Map()))
}

func (tree *RedBlackTree) UnmarshalJSON(b []byte) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = util.ComparatorString
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		tree.doSet(k, v)
	}
	return nil
}

func (tree *RedBlackTree) UnmarshalValue(value interface{}) (err error) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.comparator == nil {
		tree.comparator = util.ComparatorString
	}
	for k, v := range conv.Map(value) {
		tree.doSet(k, v)
	}
	return
}

func (tree *RedBlackTree) getComparator() func(a, b interface{}) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}
