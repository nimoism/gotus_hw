package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(value interface{}) *ListItem {
	newFront := &ListItem{Value: value, Prev: nil, Next: nil}
	l.pushItemToFront(newFront)
	return newFront
}

func (l *list) pushItemToFront(item *ListItem) {
	item.Prev, item.Next = nil, l.front
	if l.front != nil {
		l.front.Prev = item
	}
	if l.back == nil {
		l.back = item
	}
	l.front = item
	l.len++
}

func (l *list) PushBack(value interface{}) *ListItem {
	newBack := &ListItem{Value: value, Prev: l.back, Next: nil}
	if l.back != nil {
		l.back.Next = newBack
	}
	if l.front == nil {
		l.front = newBack
	}
	l.back = newBack
	l.len++
	return newBack
}

func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}
	if l.front == item {
		l.front = item.Next
	}
	if l.back == item {
		l.back = item.Prev
	}
	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}
	l.len--
}

func (l *list) MoveToFront(item *ListItem) {
	if item == nil {
		return
	}
	l.Remove(item)
	l.pushItemToFront(item)
}

func NewList() List {
	return &list{}
}
