package vmap

import (
	"github.com/luyingjie/utils/container/tree"
)

type TreeMap = tree.RedBlackTree

func NewTreeMap(comparator func(v1, v2 interface{}) int, safe ...bool) *TreeMap {
	return tree.NewRedBlackTree(comparator, safe...)
}

func NewTreeMapFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *TreeMap {
	return tree.NewRedBlackTreeFrom(comparator, data, safe...)
}
