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

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	list  map[*ListItem]struct{}
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return len(l.list)
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	current := &ListItem{Value: v}
	l.list[current] = struct{}{}
	if l.front != nil {
		current.Next = l.front
		current.Next.Prev = current
		l.front = current

		return current
	}

	l.front = current
	l.back = current

	return current
}

func (l *list) PushBack(v interface{}) *ListItem {
	current := &ListItem{Value: v}
	l.list[current] = struct{}{}

	if l.back != nil {
		current.Prev = l.back
		current.Prev.Next = current
		l.back = current

		return current
	}

	l.front = current
	l.back = current

	return current
}

func (l *list) Remove(i *ListItem) {
	if _, ok := l.list[i]; !ok {
		return
	}
	delete(l.list, i)

	count := l.Len()
	if count == 0 {
		l.front = nil
		l.back = nil

		return
	}

	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev

		return
	}

	if i.Prev != nil {
		i.Prev.Next = nil
		l.back = i.Prev
		if count == 1 {
			l.front = l.back
		}
	}
	if i.Next != nil {
		i.Next.Prev = nil
		l.front = i.Next
		if count == 1 {
			l.back = l.front
		}
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{
		list:  map[*ListItem]struct{}{},
		front: nil,
		back:  nil,
	}
}
