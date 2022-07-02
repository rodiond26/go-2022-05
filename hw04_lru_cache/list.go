package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// Node
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// Open doubly linked list
type list struct {
	length int
	head   *ListItem
	tail   *ListItem
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	switch {

	case l.length == 0:
		l.head = newNode
		l.tail = newNode

	case l.length == 1:
		l.head = newNode
		l.head.Next = l.tail
		l.tail.Prev = l.head

	case l.length > 1:
		newNode.Next = l.head
		l.head.Prev = newNode
		l.head = newNode
	}

	l.length++
	return newNode
}

func (l *list) PushBack(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	switch {

	case l.length == 0:
		l.head = newNode
		l.tail = newNode

	case l.length == 1:
		l.tail = newNode
		l.tail.Prev = l.head
		l.head.Next = l.tail

	case l.length > 1:
		newNode.Prev = l.tail
		l.tail.Next = newNode
		l.tail = newNode
	}

	l.length++
	return newNode
}

func (l *list) Remove(rem *ListItem) {
	switch {

	case l.length == 1:
		l.head = nil
		l.tail = nil

	case l.length == 2 && rem == l.head:
		l.head = l.tail
		l.head.Next = nil
		l.head.Prev = nil

	case l.length == 2 && rem == l.tail:
		l.tail = l.head
		l.tail.Next = nil
		l.tail.Prev = nil

	case l.length > 2 && rem == l.head:
		l.head = l.head.Next
		l.head.Prev = nil

	case l.length > 2 && rem == l.tail:
		l.tail = l.tail.Prev
		l.tail.Next = nil

	default:
		rem.Next.Prev = rem.Prev
		rem.Prev.Next = rem.Next
	}

	l.length--
}

func (l *list) MoveToFront(moved *ListItem) {
	switch {

	case l.length == 1 || l.head == moved:

	case l.length == 2 && l.tail == moved:
		l.head, l.tail = l.tail, l.head
		l.head.Next = l.tail
		l.head.Prev = nil
		l.tail.Prev = l.head
		l.tail.Next = nil

	case l.length > 2 && l.tail == moved:
		l.tail = l.tail.Prev
		l.tail.Next = nil
		l.head.Prev = moved
		moved.Next = l.head
		l.head = moved
		l.head.Prev = nil

	default:
		moved.Prev.Next = moved.Next
		moved.Next.Prev = moved.Prev
		l.head.Prev = moved
		moved.Next = l.head
		moved.Prev = nil
		l.head = moved
	}
}
