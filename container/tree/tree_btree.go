package tree

import (
	"bytes"
	"fmt"
	"strings"
	vvar "utils/container/var"
	"utils/util/json"

	"utils/conv"

	"utils/util/rwmutex"
)

type Tree struct {
	mu         rwmutex.RWMutex
	root       *TreeNode
	comparator func(v1, v2 interface{}) int
	size       int
	m          int
}

type TreeNode struct {
	Parent   *TreeNode
	Entries  []*TreeEntry
	Children []*TreeNode
}

type TreeEntry struct {
	Key   interface{}
	Value interface{}
}

func NewTree(m int, comparator func(v1, v2 interface{}) int, safe ...bool) *Tree {
	if m < 3 {
		panic("Invalid order, should be at least 3")
	}
	return &Tree{
		comparator: comparator,
		mu:         rwmutex.Create(safe...),
		m:          m,
	}
}

func NewTreeFrom(m int, comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *Tree {
	tree := NewTree(m, comparator, safe...)
	for k, v := range data {
		tree.doSet(k, v)
	}
	return tree
}

func (tree *Tree) Clone() *Tree {
	newTree := NewTree(tree.m, tree.comparator, !tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

func (tree *Tree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.doSet(key, value)
}

func (tree *Tree) doSet(key interface{}, value interface{}) {
	entry := &TreeEntry{Key: key, Value: value}
	if tree.root == nil {
		tree.root = &TreeNode{Entries: []*TreeEntry{entry}, Children: []*TreeNode{}}
		tree.size++
		return
	}

	if tree.insert(tree.root, entry) {
		tree.size++
	}
}

func (tree *Tree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for k, v := range data {
		tree.doSet(k, v)
	}
}

func (tree *Tree) Get(key interface{}) (value interface{}) {
	value, _ = tree.Search(key)
	return
}

func (tree *Tree) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if entry := tree.doSearch(key); entry != nil {
		return entry.Value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		tree.doSet(key, value)
	}
	return value
}

func (tree *Tree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (tree *Tree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (tree *Tree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (tree *Tree) GetVar(key interface{}) *vvar.Var {
	return vvar.New(tree.Get(key))
}

func (tree *Tree) GetVarOrSet(key interface{}, value interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSet(key, value))
}

func (tree *Tree) GetVarOrSetFunc(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSetFunc(key, f))
}

func (tree *Tree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSetFuncLock(key, f))
}

func (tree *Tree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (tree *Tree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (tree *Tree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (tree *Tree) Contains(key interface{}) bool {
	_, ok := tree.Search(key)
	return ok
}

func (tree *Tree) doRemove(key interface{}) (value interface{}) {
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		value = node.Entries[index].Value
		tree.delete(node, index)
		tree.size--
	}
	return
}

func (tree *Tree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	return tree.doRemove(key)
}

func (tree *Tree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.doRemove(key)
	}
}

func (tree *Tree) IsEmpty() bool {
	return tree.Size() == 0
}

func (tree *Tree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

func (tree *Tree) Keys() []interface{} {
	keys := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

func (tree *Tree) Values() []interface{} {
	values := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

func (tree *Tree) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func (tree *Tree) MapStrAny() map[string]interface{} {
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

func (tree *Tree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

func (tree *Tree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for k, v := range data {
		tree.doSet(k, v)
	}
}

func (tree *Tree) Height() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.root.height()
}

func (tree *Tree) Left() *TreeEntry {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.left(tree.root)
	return node.Entries[0]
}

func (tree *Tree) Right() *TreeEntry {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.right(tree.root)
	return node.Entries[len(node.Entries)-1]
}

func (tree *Tree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	var buffer bytes.Buffer
	if tree.size != 0 {
		tree.output(&buffer, tree.root, 0, true)
	}
	return buffer.String()
}

func (tree *Tree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		return node.Entries[index].Value, true
	}
	return nil, false
}

func (tree *Tree) doSearch(key interface{}) *TreeEntry {
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		return node.Entries[index]
	}
	return nil
}

func (tree *Tree) Print() {
	fmt.Println(tree.String())
}

func (tree *Tree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

func (tree *Tree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

func (tree *Tree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.left(tree.root)
	if node == nil {
		return
	}
	tree.doIteratorAsc(node, node.Entries[0], 0, f)
}

func (tree *Tree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, index, found := tree.searchRecursively(tree.root, key)
	if match {
		if found {
			tree.doIteratorAsc(node, node.Entries[index], index, f)
		}
	} else {
		tree.doIteratorAsc(node, node.Entries[index], index, f)
	}
}

func (tree *Tree) doIteratorAsc(node *TreeNode, entry *TreeEntry, index int, f func(key, value interface{}) bool) {
	first := true
loop:
	if entry == nil {
		return
	}
	if !f(entry.Key, entry.Value) {
		return
	}

	if !first {
		index, _ = tree.search(node, entry.Key)
	} else {
		first = false
	}

	if index+1 < len(node.Children) {
		node = node.Children[index+1]

		for len(node.Children) > 0 {
			node = node.Children[0]
		}

		entry = node.Entries[0]
		goto loop
	}

	if index+1 < len(node.Entries) {
		entry = node.Entries[index+1]
		goto loop
	}

	for node.Parent != nil {
		node = node.Parent
		index, _ = tree.search(node, entry.Key)

		if index < len(node.Entries) {
			entry = node.Entries[index]
			goto loop
		}
	}
}

func (tree *Tree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.right(tree.root)
	if node == nil {
		return
	}
	index := len(node.Entries) - 1
	entry := node.Entries[index]
	tree.doIteratorDesc(node, entry, index, f)
}

func (tree *Tree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node, index, found := tree.searchRecursively(tree.root, key)
	if match {
		if found {
			tree.doIteratorDesc(node, node.Entries[index], index, f)
		}
	} else {
		tree.doIteratorDesc(node, node.Entries[index], index, f)
	}
}

func (tree *Tree) doIteratorDesc(node *TreeNode, entry *TreeEntry, index int, f func(key, value interface{}) bool) {
	first := true
loop:
	if entry == nil {
		return
	}
	if !f(entry.Key, entry.Value) {
		return
	}

	if !first {
		index, _ = tree.search(node, entry.Key)
	} else {
		first = false
	}

	if index < len(node.Children) {
		node = node.Children[index]

		for len(node.Children) > 0 {
			node = node.Children[len(node.Children)-1]
		}

		entry = node.Entries[len(node.Entries)-1]
		goto loop
	}

	if index-1 >= 0 {
		entry = node.Entries[index-1]
		goto loop
	}

	for node.Parent != nil {
		node = node.Parent

		index, _ = tree.search(node, entry.Key)

		if index-1 >= 0 {
			entry = node.Entries[index-1]
			goto loop
		}
	}
}

func (tree *Tree) output(buffer *bytes.Buffer, node *TreeNode, level int, isTail bool) {
	for e := 0; e < len(node.Entries)+1; e++ {
		if e < len(node.Children) {
			tree.output(buffer, node.Children[e], level+1, true)
		}
		if e < len(node.Entries) {
			if _, err := buffer.WriteString(strings.Repeat("    ", level)); err != nil {
			}
			if _, err := buffer.WriteString(fmt.Sprintf("%v", node.Entries[e].Key) + "\n"); err != nil {
			}
		}
	}
}

func (node *TreeNode) height() int {
	h := 0
	n := node
	for ; n != nil; n = n.Children[0] {
		h++
		if len(n.Children) == 0 {
			break
		}
	}
	return h
}

func (tree *Tree) isLeaf(node *TreeNode) bool {
	return len(node.Children) == 0
}

//func (tree *Tree) isFull(node *TreeNode) bool {
//	return len(node.Entries) == tree.maxEntries()
//}

func (tree *Tree) shouldSplit(node *TreeNode) bool {
	return len(node.Entries) > tree.maxEntries()
}

func (tree *Tree) maxChildren() int {
	return tree.m
}

func (tree *Tree) minChildren() int {
	return (tree.m + 1) / 2 // ceil(m/2)
}

func (tree *Tree) maxEntries() int {
	return tree.maxChildren() - 1
}

func (tree *Tree) minEntries() int {
	return tree.minChildren() - 1
}

func (tree *Tree) middle() int {
	return (tree.m - 1) / 2
}

func (tree *Tree) search(node *TreeNode, key interface{}) (index int, found bool) {
	low, mid, high := 0, 0, len(node.Entries)-1
	for low <= high {
		mid = (high + low) / 2
		compare := tree.getComparator()(key, node.Entries[mid].Key)
		switch {
		case compare > 0:
			low = mid + 1
		case compare < 0:
			high = mid - 1
		case compare == 0:
			return mid, true
		}
	}
	return low, false
}

func (tree *Tree) searchRecursively(startNode *TreeNode, key interface{}) (node *TreeNode, index int, found bool) {
	if tree.size == 0 {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = tree.search(node, key)
		if found {
			return node, index, true
		}
		if tree.isLeaf(node) {
			return node, index, false
		}
		node = node.Children[index]
	}
}

func (tree *Tree) insert(node *TreeNode, entry *TreeEntry) (inserted bool) {
	if tree.isLeaf(node) {
		return tree.insertIntoLeaf(node, entry)
	}
	return tree.insertIntoInternal(node, entry)
}

func (tree *Tree) insertIntoLeaf(node *TreeNode, entry *TreeEntry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	node.Entries = append(node.Entries, nil)
	copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])
	node.Entries[insertPosition] = entry
	tree.split(node)
	return true
}

func (tree *Tree) insertIntoInternal(node *TreeNode, entry *TreeEntry) (inserted bool) {
	insertPosition, found := tree.search(node, entry.Key)
	if found {
		node.Entries[insertPosition] = entry
		return false
	}
	return tree.insert(node.Children[insertPosition], entry)
}

func (tree *Tree) split(node *TreeNode) {
	if !tree.shouldSplit(node) {
		return
	}

	if node == tree.root {
		tree.splitRoot()
		return
	}

	tree.splitNonRoot(node)
}

func (tree *Tree) splitNonRoot(node *TreeNode) {
	middle := tree.middle()
	parent := node.Parent

	left := &TreeNode{Entries: append([]*TreeEntry(nil), node.Entries[:middle]...), Parent: parent}
	right := &TreeNode{Entries: append([]*TreeEntry(nil), node.Entries[middle+1:]...), Parent: parent}

	if !tree.isLeaf(node) {
		left.Children = append([]*TreeNode(nil), node.Children[:middle+1]...)
		right.Children = append([]*TreeNode(nil), node.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	insertPosition, _ := tree.search(parent, node.Entries[middle].Key)

	parent.Entries = append(parent.Entries, nil)
	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	parent.Entries[insertPosition] = node.Entries[middle]

	parent.Children[insertPosition] = left

	parent.Children = append(parent.Children, nil)
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	parent.Children[insertPosition+1] = right

	tree.split(parent)
}

func (tree *Tree) splitRoot() {
	middle := tree.middle()
	left := &TreeNode{Entries: append([]*TreeEntry(nil), tree.root.Entries[:middle]...)}
	right := &TreeNode{Entries: append([]*TreeEntry(nil), tree.root.Entries[middle+1:]...)}

	if !tree.isLeaf(tree.root) {
		left.Children = append([]*TreeNode(nil), tree.root.Children[:middle+1]...)
		right.Children = append([]*TreeNode(nil), tree.root.Children[middle+1:]...)
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	newRoot := &TreeNode{
		Entries:  []*TreeEntry{tree.root.Entries[middle]},
		Children: []*TreeNode{left, right},
	}

	left.Parent = newRoot
	right.Parent = newRoot
	tree.root = newRoot
}

func setParent(nodes []*TreeNode, parent *TreeNode) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (tree *Tree) left(node *TreeNode) *TreeNode {
	if tree.size == 0 {
		return nil
	}
	current := node
	for {
		if tree.isLeaf(current) {
			return current
		}
		current = current.Children[0]
	}
}

func (tree *Tree) right(node *TreeNode) *TreeNode {
	if tree.size == 0 {
		return nil
	}
	current := node
	for {
		if tree.isLeaf(current) {
			return current
		}
		current = current.Children[len(current.Children)-1]
	}
}

func (tree *Tree) leftSibling(node *TreeNode, key interface{}) (*TreeNode, int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index--
		if index >= 0 && index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

func (tree *Tree) rightSibling(node *TreeNode, key interface{}) (*TreeNode, int) {
	if node.Parent != nil {
		index, _ := tree.search(node.Parent, key)
		index++
		if index < len(node.Parent.Children) {
			return node.Parent.Children[index], index
		}
	}
	return nil, -1
}

func (tree *Tree) delete(node *TreeNode, index int) {
	if tree.isLeaf(node) {
		deletedKey := node.Entries[index].Key
		tree.deleteEntry(node, index)
		tree.rebalance(node, deletedKey)
		if len(tree.root.Entries) == 0 {
			tree.root = nil
		}
		return
	}

	leftLargestNode := tree.right(node.Children[index])
	leftLargestEntryIndex := len(leftLargestNode.Entries) - 1
	node.Entries[index] = leftLargestNode.Entries[leftLargestEntryIndex]
	deletedKey := leftLargestNode.Entries[leftLargestEntryIndex].Key
	tree.deleteEntry(leftLargestNode, leftLargestEntryIndex)
	tree.rebalance(leftLargestNode, deletedKey)
}

func (tree *Tree) rebalance(node *TreeNode, deletedKey interface{}) {
	if node == nil || len(node.Entries) >= tree.minEntries() {
		return
	}

	leftSibling, leftSiblingIndex := tree.leftSibling(node, deletedKey)
	if leftSibling != nil && len(leftSibling.Entries) > tree.minEntries() {
		node.Entries = append([]*TreeEntry{node.Parent.Entries[leftSiblingIndex]}, node.Entries...)
		node.Parent.Entries[leftSiblingIndex] = leftSibling.Entries[len(leftSibling.Entries)-1]
		tree.deleteEntry(leftSibling, len(leftSibling.Entries)-1)
		if !tree.isLeaf(leftSibling) {
			leftSiblingRightMostChild := leftSibling.Children[len(leftSibling.Children)-1]
			leftSiblingRightMostChild.Parent = node
			node.Children = append([]*TreeNode{leftSiblingRightMostChild}, node.Children...)
			tree.deleteChild(leftSibling, len(leftSibling.Children)-1)
		}
		return
	}

	rightSibling, rightSiblingIndex := tree.rightSibling(node, deletedKey)
	if rightSibling != nil && len(rightSibling.Entries) > tree.minEntries() {

		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Parent.Entries[rightSiblingIndex-1] = rightSibling.Entries[0]
		tree.deleteEntry(rightSibling, 0)
		if !tree.isLeaf(rightSibling) {
			rightSiblingLeftMostChild := rightSibling.Children[0]
			rightSiblingLeftMostChild.Parent = node
			node.Children = append(node.Children, rightSiblingLeftMostChild)
			tree.deleteChild(rightSibling, 0)
		}
		return
	}

	if rightSibling != nil {
		node.Entries = append(node.Entries, node.Parent.Entries[rightSiblingIndex-1])
		node.Entries = append(node.Entries, rightSibling.Entries...)
		deletedKey = node.Parent.Entries[rightSiblingIndex-1].Key
		tree.deleteEntry(node.Parent, rightSiblingIndex-1)
		tree.appendChildren(node.Parent.Children[rightSiblingIndex], node)
		tree.deleteChild(node.Parent, rightSiblingIndex)
	} else if leftSibling != nil {
		entries := append([]*TreeEntry(nil), leftSibling.Entries...)
		entries = append(entries, node.Parent.Entries[leftSiblingIndex])
		node.Entries = append(entries, node.Entries...)
		deletedKey = node.Parent.Entries[leftSiblingIndex].Key
		tree.deleteEntry(node.Parent, leftSiblingIndex)
		tree.prependChildren(node.Parent.Children[leftSiblingIndex], node)
		tree.deleteChild(node.Parent, leftSiblingIndex)
	}

	if node.Parent == tree.root && len(tree.root.Entries) == 0 {
		tree.root = node
		node.Parent = nil
		return
	}

	tree.rebalance(node.Parent, deletedKey)
}

func (tree *Tree) prependChildren(fromNode *TreeNode, toNode *TreeNode) {
	children := append([]*TreeNode(nil), fromNode.Children...)
	toNode.Children = append(children, toNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree) appendChildren(fromNode *TreeNode, toNode *TreeNode) {
	toNode.Children = append(toNode.Children, fromNode.Children...)
	setParent(fromNode.Children, toNode)
}

func (tree *Tree) deleteEntry(node *TreeNode, index int) {
	copy(node.Entries[index:], node.Entries[index+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

func (tree *Tree) deleteChild(node *TreeNode, index int) {
	if index >= len(node.Children) {
		return
	}
	copy(node.Children[index:], node.Children[index+1:])
	node.Children[len(node.Children)-1] = nil
	node.Children = node.Children[:len(node.Children)-1]
}

func (tree *Tree) MarshalJSON() ([]byte, error) {
	return json.Marshal(tree.Map())
}

func (tree *Tree) getComparator() func(a, b interface{}) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}
