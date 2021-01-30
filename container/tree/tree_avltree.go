package tree

import (
	"fmt"
	vvar "utils/container/var"
	"utils/util/json"

	"utils/os/conv"

	"utils/util/rwmutex"
)

type AVLTree struct {
	mu         rwmutex.RWMutex
	root       *AVLTreeNode
	comparator func(v1, v2 interface{}) int
	size       int
}

type AVLTreeNode struct {
	Key      interface{}
	Value    interface{}
	parent   *AVLTreeNode
	children [2]*AVLTreeNode
	b        int8
}

func NewAVLTree(comparator func(v1, v2 interface{}) int, safe ...bool) *AVLTree {
	return &AVLTree{
		mu:         rwmutex.Create(safe...),
		comparator: comparator,
	}
}

func NewAVLTreeFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *AVLTree {
	tree := NewAVLTree(comparator, safe...)
	for k, v := range data {
		tree.put(k, v, nil, &tree.root)
	}
	return tree
}

func (tree *AVLTree) Clone() *AVLTree {
	newTree := NewAVLTree(tree.comparator, !tree.mu.IsSafe())
	newTree.Sets(tree.Map())
	return newTree
}

func (tree *AVLTree) Set(key interface{}, value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.put(key, value, nil, &tree.root)
}

func (tree *AVLTree) Sets(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for key, value := range data {
		tree.put(key, value, nil, &tree.root)
	}
}

func (tree *AVLTree) Search(key interface{}) (value interface{}, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if node, found := tree.doSearch(key); found {
		return node.Value, true
	}
	return nil, false
}

func (tree *AVLTree) doSearch(key interface{}) (node *AVLTreeNode, found bool) {
	node = tree.root
	for node != nil {
		cmp := tree.getComparator()(key, node.Key)
		switch {
		case cmp == 0:
			return node, true
		case cmp < 0:
			node = node.children[0]
		case cmp > 0:
			node = node.children[1]
		}
	}
	return nil, false
}

func (tree *AVLTree) Get(key interface{}) (value interface{}) {
	value, _ = tree.Search(key)
	return
}

func (tree *AVLTree) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if node, found := tree.doSearch(key); found {
		return node.Value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		tree.put(key, value, nil, &tree.root)
	}
	return value
}

func (tree *AVLTree) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (tree *AVLTree) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (tree *AVLTree) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := tree.Search(key); !ok {
		return tree.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (tree *AVLTree) GetVar(key interface{}) *vvar.Var {
	return vvar.New(tree.Get(key))
}

func (tree *AVLTree) GetVarOrSet(key interface{}, value interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSet(key, value))
}

func (tree *AVLTree) GetVarOrSetFunc(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSetFunc(key, f))
}

func (tree *AVLTree) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(tree.GetOrSetFuncLock(key, f))
}

func (tree *AVLTree) SetIfNotExist(key interface{}, value interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (tree *AVLTree) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (tree *AVLTree) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !tree.Contains(key) {
		tree.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (tree *AVLTree) Contains(key interface{}) bool {
	_, ok := tree.Search(key)
	return ok
}

func (tree *AVLTree) Remove(key interface{}) (value interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	value, _ = tree.remove(key, &tree.root)
	return
}

func (tree *AVLTree) Removes(keys []interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for _, key := range keys {
		tree.remove(key, &tree.root)
	}
}

func (tree *AVLTree) IsEmpty() bool {
	return tree.Size() == 0
}

func (tree *AVLTree) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

func (tree *AVLTree) Keys() []interface{} {
	keys := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		keys[index] = key
		index++
		return true
	})
	return keys
}

func (tree *AVLTree) Values() []interface{} {
	values := make([]interface{}, tree.Size())
	index := 0
	tree.IteratorAsc(func(key, value interface{}) bool {
		values[index] = value
		index++
		return true
	})
	return values
}

func (tree *AVLTree) Left() *AVLTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.bottom(0)
	if tree.mu.IsSafe() {
		return &AVLTreeNode{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

func (tree *AVLTree) Right() *AVLTreeNode {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	node := tree.bottom(1)
	if tree.mu.IsSafe() {
		return &AVLTreeNode{
			Key:   node.Key,
			Value: node.Value,
		}
	}
	return node
}

func (tree *AVLTree) Floor(key interface{}) (floor *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		c := tree.getComparator()(key, n.Key)
		switch {
		case c == 0:
			return n, true
		case c < 0:
			n = n.children[0]
		case c > 0:
			floor, found = n, true
			n = n.children[1]
		}
	}
	if found {
		return
	}
	return nil, false
}

func (tree *AVLTree) Ceiling(key interface{}) (ceiling *AVLTreeNode, found bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	n := tree.root
	for n != nil {
		c := tree.getComparator()(key, n.Key)
		switch {
		case c == 0:
			return n, true
		case c > 0:
			n = n.children[1]
		case c < 0:
			ceiling, found = n, true
			n = n.children[0]
		}
	}
	if found {
		return
	}
	return nil, false
}

func (tree *AVLTree) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
}

func (tree *AVLTree) Replace(data map[interface{}]interface{}) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.root = nil
	tree.size = 0
	for key, value := range data {
		tree.put(key, value, nil, &tree.root)
	}
}

func (tree *AVLTree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	str := ""
	if tree.size != 0 {
		output(tree.root, "", true, &str)
	}
	return str
}

func (tree *AVLTree) Print() {
	fmt.Println(tree.String())
}

func (tree *AVLTree) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func (tree *AVLTree) MapStrAny() map[string]interface{} {
	m := make(map[string]interface{}, tree.Size())
	tree.IteratorAsc(func(key, value interface{}) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

func (tree *AVLTree) Flip(comparator ...func(v1, v2 interface{}) int) {
	t := (*AVLTree)(nil)
	if len(comparator) > 0 {
		t = NewAVLTree(comparator[0], !tree.mu.IsSafe())
	} else {
		t = NewAVLTree(tree.comparator, !tree.mu.IsSafe())
	}
	tree.IteratorAsc(func(key, value interface{}) bool {
		t.put(value, key, nil, &t.root)
		return true
	})
	tree.mu.Lock()
	tree.root = t.root
	tree.size = t.size
	tree.mu.Unlock()
}

func (tree *AVLTree) Iterator(f func(key, value interface{}) bool) {
	tree.IteratorAsc(f)
}

func (tree *AVLTree) IteratorFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
	tree.IteratorAscFrom(key, match, f)
}

func (tree *AVLTree) IteratorAsc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorAsc(tree.bottom(0), f)
}

func (tree *AVLTree) IteratorAscFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

func (tree *AVLTree) doIteratorAsc(node *AVLTreeNode, f func(key, value interface{}) bool) {
	for node != nil {
		if !f(node.Key, node.Value) {
			return
		}
		node = node.Next()
	}
}

func (tree *AVLTree) IteratorDesc(f func(key, value interface{}) bool) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	tree.doIteratorDesc(tree.bottom(1), f)
}

func (tree *AVLTree) IteratorDescFrom(key interface{}, match bool, f func(key, value interface{}) bool) {
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

func (tree *AVLTree) doIteratorDesc(node *AVLTreeNode, f func(key, value interface{}) bool) {
	for node != nil {
		if !f(node.Key, node.Value) {
			return
		}
		node = node.Prev()
	}
}

func (tree *AVLTree) put(key interface{}, value interface{}, p *AVLTreeNode, qp **AVLTreeNode) bool {
	q := *qp
	if q == nil {
		tree.size++
		*qp = &AVLTreeNode{Key: key, Value: value, parent: p}
		return true
	}

	c := tree.getComparator()(key, q.Key)
	if c == 0 {
		q.Key = key
		q.Value = value
		return false
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	if tree.put(key, value, q, &q.children[a]) {
		return putFix(int8(c), qp)
	}
	return false
}

func (tree *AVLTree) remove(key interface{}, qp **AVLTreeNode) (value interface{}, fix bool) {
	q := *qp
	if q == nil {
		return nil, false
	}

	c := tree.getComparator()(key, q.Key)
	if c == 0 {
		tree.size--
		value = q.Value
		fix = true
		if q.children[1] == nil {
			if q.children[0] != nil {
				q.children[0].parent = q.parent
			}
			*qp = q.children[0]
			return
		}
		if removeMin(&q.children[1], &q.Key, &q.Value) {
			return value, removeFix(-1, qp)
		}
		return
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	value, fix = tree.remove(key, &q.children[a])
	if fix {
		return value, removeFix(int8(-c), qp)
	}
	return value, false
}

func removeMin(qp **AVLTreeNode, minKey *interface{}, minVal *interface{}) bool {
	q := *qp
	if q.children[0] == nil {
		*minKey = q.Key
		*minVal = q.Value
		if q.children[1] != nil {
			q.children[1].parent = q.parent
		}
		*qp = q.children[1]
		return true
	}
	fix := removeMin(&q.children[0], minKey, minVal)
	if fix {
		return removeFix(1, qp)
	}
	return false
}

func putFix(c int8, t **AVLTreeNode) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.children[(c+1)/2].b == c {
		s = singleRotate(c, s)
	} else {
		s = doubleRotate(c, s)
	}
	*t = s
	return false
}

func removeFix(c int8, t **AVLTreeNode) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return false
	}

	if s.b == -c {
		s.b = 0
		return true
	}

	a := (c + 1) / 2
	if s.children[a].b == 0 {
		s = rotate(c, s)
		s.b = -c
		*t = s
		return false
	}

	if s.children[a].b == c {
		s = singleRotate(c, s)
	} else {
		s = doubleRotate(c, s)
	}
	*t = s
	return true
}

func singleRotate(c int8, s *AVLTreeNode) *AVLTreeNode {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doubleRotate(c int8, s *AVLTreeNode) *AVLTreeNode {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = rotate(-c, s.children[a])
	p := rotate(c, s)

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate(c int8, s *AVLTreeNode) *AVLTreeNode {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = r.children[a^1]
	if s.children[a] != nil {
		s.children[a].parent = s
	}
	r.children[a^1] = s
	r.parent = s.parent
	s.parent = r
	return r
}

func (tree *AVLTree) bottom(d int) *AVLTreeNode {
	n := tree.root
	if n == nil {
		return nil
	}

	for c := n.children[d]; c != nil; c = n.children[d] {
		n = c
	}
	return n
}

func (node *AVLTreeNode) Prev() *AVLTreeNode {
	return node.walk1(0)
}

func (node *AVLTreeNode) Next() *AVLTreeNode {
	return node.walk1(1)
}

func (node *AVLTreeNode) walk1(a int) *AVLTreeNode {
	if node == nil {
		return nil
	}
	n := node
	if n.children[a] != nil {
		n = n.children[a]
		for n.children[a^1] != nil {
			n = n.children[a^1]
		}
		return n
	}

	p := n.parent
	for p != nil && p.children[a] == n {
		n = p
		p = p.parent
	}
	return p
}

func output(node *AVLTreeNode, prefix string, isTail bool, str *string) {
	if node.children[1] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.children[1], newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += fmt.Sprintf("%v\n", node.Key)
	if node.children[0] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.children[0], newPrefix, true, str)
	}
}

func (tree *AVLTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(tree.Map())
}

func (tree *AVLTree) getComparator() func(a, b interface{}) int {
	if tree.comparator == nil {
		panic("comparator is missing for tree")
	}
	return tree.comparator
}
