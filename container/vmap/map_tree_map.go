package vmap

import (
	"github.com/luyingjie/utils/container/vtree"
)

type TreeMap = vtree.RedBlackTree

func NewTreeMap(comparator func(v1, v2 interface{}) int, safe ...bool) *TreeMap {
	return vtree.NewRedBlackTree(comparator, safe...)
}

func NewTreeMapFrom(comparator func(v1, v2 interface{}) int, data map[interface{}]interface{}, safe ...bool) *TreeMap {
	return vtree.NewRedBlackTreeFrom(comparator, data, safe...)
}
