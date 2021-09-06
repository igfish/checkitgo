package checkit

type Tree struct{}
type color uint8

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

const (
	black color = iota
	red
)

type entry[K Ordered, V any] struct {
	key    K
	value  V
	color  color
	parent *entry[K, V]
	left   *entry[K, V]
	right  *entry[K, V]
}

func (e *entry[K, V]) black() {
	if e != nil {
		e.color = black
	}
}

func (e *entry[K, V]) red() {
	if e != nil {
		e.color = red
	}
}

func (e *entry[K, V]) isLeftChild() bool {
	return e.parent != nil && e.parent.left == e
}

func (e *entry[K, V]) isRightChild() bool {
	return e.parent != nil && e.parent.right == e
}

func (e *entry[K, V]) sibling() *entry[K, V] {
	switch {
	case e.parent == nil:
		return nil
	case e.isLeftChild():
		return e.parent.right
	default:
		return e.parent.left
	}
}

func (e *entry[K, V]) uncle() *entry[K, V] {
	switch {
	case e.parent == nil:
		return nil
	default:
		return e.parent.sibling()
	}
}

func (e *entry[K, V]) isBlack() bool {
	return e == nil || e.color == black
}

func newEntry[K Ordered, V any](key K, value V, color color, parent *entry[K, V]) *entry[K, V] {
	return &entry[K, V]{
		key:    key,
		value:  value,
		color:  color,
		parent: parent,
		left:   nil,
		right:  nil,
	}

}

type rbtree[K Ordered, V any] struct {
	size int
	root *entry[K, V]
	m    map[K]V
}

func (r *rbtree[K, V]) Len() int { return r.size }

func (r *rbtree[K, V]) Put(key K, value V) {
	if r.root == nil {
		r.root = newEntry(key, value, black, nil)
		r.size++
		return
	}
	root := r.root
	for root != nil {
		switch {
		case root.key == key:
			root.value = value
			return
		case root.key < key:
			if root.right != nil {
				root = root.right
				continue
			} else {
				root.right = newEntry(key, value, red, root)
				r.size++
				r.afterPut(root.right)
				return
			}
		default:
			if root.left != nil {
				root = root.left
				continue
			} else {
				root.left = newEntry(key, value, red, root)
				r.size++
				r.afterPut(root.left)
				return
			}
		}
	}
}

func (r *rbtree[K, V]) afterPut(e *entry[K, V]) {
	parent := e.parent
	if parent == nil {
		// The parent node is nil, added the first root node,
		// set its color black.
		e.black()
		return
	}
	if parent.isBlack() {
		// The parent node is black and no adjustment is required.
		return
	}
	// The parent node is red. Let's adjust.
	uncle := e.uncle()
	grad := parent.parent // The grandfather node must be black.
	grad.red()
	parent.black()
	switch {
	case uncle.isBlack():
		// Uncle node is nil (nil node is considered black).
		// Unbalance and rebalance.
		r.rebalance(grad, parent, e)
	default:
		// The uncle node is red, set its color black.
		uncle.black()
		// Grandfather nodes is colored red, which may cause imbalances.
		// Adjust the grandfather node as if it were newly added.
		r.afterPut(grad)
	}
}

func (r *rbtree[K, V]) rotateLeft(e *entry[K, V]) {
	if e == nil || e.right == nil {
		return
	}
	grand, parent, child := e, e.right, e.right.left
	grand.right = child
	parent.left = grand
	r.afterRotate(grand, parent, child)
}

func (r *rbtree[K, V]) rotateRight(e *entry[K, V]) {
	if e == nil || e.left == nil {
		return
	}
	grand, parent, child := e, e.left, e.left.right
	grand.left = child
	parent.right = grand
	r.afterRotate(grand, parent, child)
}

func (r *rbtree[K, V]) afterRotate(grand, parent, child *entry[K, V]) {
	switch {
	case grand.isLeftChild():
		grand.parent.left = parent
	case grand.isRightChild():
		grand.parent.right = parent
	default:
		r.root = parent
	}
	parent.parent = grand.parent
	grand.parent = parent
	if child != nil {
		child.parent = grand
	}
}

func (r *rbtree[K, V]) rebalance(grad, parent, child *entry[K, V]) {
	switch {
	case parent.isLeftChild(): // L
		if child.isRightChild() { // LR
			r.rotateLeft(parent)
		}
		r.rotateRight(grad)
	default: // R
		if child.isLeftChild() { // RL
			r.rotateRight(parent)
		}
		r.rotateLeft(grad)
	}
}

func (r *rbtree[K, V]) Get(key K) (v V, ok bool) {
	root := r.root
	for root != nil {
		if root.key == key {
			return root.value, true
		} else if root.key < key {
			root = root.right
		} else {
			root = root.left
		}
	}
	return
}

func (r *rbtree[K, V]) Range(from, to K) {}

func NewRBTreeMap[K Ordered, V any]() *rbtree[K, V] {
	return &rbtree[K, V]{
		size: 0,
		root: nil,
		m:    make(map[K]V),
	}
}

type Pair[K, V any] struct {
	K K
	V V
}

func newPair[K, V any](k K, v V) *Pair[K, V] {
	return &Pair[K, V]{k, v}
}

type rbtreeItor[K Ordered, V, T any] struct {
	r *entry[K, V]
	s []*entry[K, V]
}

func (itor *rbtreeItor[K, V, T]) Next() (*Pair[K, V], bool) {
	for itor.r != nil || len(itor.s) > 0 {
		if itor.r != nil {
			itor.s = append(itor.s, itor.r)
			itor.r = itor.r.left
		} else {
			itor.r = itor.s[len(itor.s)-1]
			itor.s = itor.s[:len(itor.s)-1]
			p := Pair[K, V]{itor.r.key, itor.r.value}
			itor.r = itor.r.right
			return &p, true
		}
	}
	return nil, false
}

func (r *rbtree[K, V]) Iter() Iterator[*Pair[K, V]] {
	itor := &rbtreeItor[K, V, *Pair[K, V]]{
		r: r.root,
		s: make([]*entry[K, V], 0),
	}
	return itor
}

var _ Iterator[*Pair[int, int]] = (*rbtreeItor[int, int, Pair[int, int]])(nil)
