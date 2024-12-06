package linkedList

import (
	"fmt"
	"reflect"
)

type ILinkedList[T any] interface {
	AddFront(item *T)
	AddBack(item *T)
	FindByStructField(fieldName string, fieldValue any) (*T, bool)
	RemoveByStructField(fieldName string, fieldValue any) error
	RemoveFront() *T
	GetSize() int
	PrintList()
	Clear()
}

type node[T any] struct {
	element  *T
	nextNode *node[T]
}

type LinkedList[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
}

func New[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func NewWithInitValue[T any](item *T) *LinkedList[T] {
	newNode := &node[T]{element: item, nextNode: nil}
	return &LinkedList[T]{
		head:   newNode,
		tail:   newNode,
		length: 1,
	}
}

func (list *LinkedList[T]) GetSize() int {
	return list.length
}

func (list *LinkedList[T]) AddFront(item *T) {
	newNode := &node[T]{element: item, nextNode: nil}
	if list.head == nil {
		list.head = newNode
		list.tail = newNode
		list.length++
		return
	}

	newNode.nextNode = list.head
	list.head = newNode
	list.length++
}

func (list *LinkedList[T]) AddBack(item *T) {
	newNode := &node[T]{element: item, nextNode: nil}
	if list.head == nil {
		list.head = newNode
		list.tail = newNode
		list.length++
		return
	}

	list.tail.nextNode = newNode
	list.tail = newNode
	list.length++
}

// FindByStructField finds element in linked list by fieldName only if this field is accessible.
func (list *LinkedList[T]) FindByStructField(fieldName string, fieldValue interface{}) (*T, bool) {
	current := list.head

	var field reflect.Value
	for current != nil {
		field = reflect.ValueOf(*current.element).FieldByName(fieldName)
		if !field.IsValid() || !field.CanInterface() {
			return nil, false
		}

		if reflect.DeepEqual(field.Interface(), fieldValue) {
			return current.element, true
		}
		current = current.nextNode
	}

	return nil, false
}

// RemoveByStructField removes element in linked list by fieldName only if this field is accessible.
func (list *LinkedList[T]) RemoveByStructField(fieldName string, fieldValue any) error {
	if list.length == 0 {
		return ErrorElementNotRemoved
	}
	field := reflect.ValueOf(*list.head.element).FieldByName(fieldName)
	if !field.IsValid() || !field.CanInterface() {
		return ErrorElementNotRemoved
	}

	if reflect.DeepEqual(field.Interface(), fieldValue) {
		list.head = list.head.nextNode
		list.length--
		if list.head == nil {
			list.tail = nil
		}
		return nil
	}

	var prev *node[T]
	current := list.head
	for current != nil {
		field = reflect.ValueOf(*current.element).FieldByName(fieldName)
		if !field.IsValid() || !field.CanInterface() {
			return ErrorElementNotRemoved
		}

		if reflect.DeepEqual(field.Interface(), fieldValue) {
			prev.nextNode = current.nextNode
			list.length--

			if current.nextNode == nil {
				list.tail = prev
			}
			return nil
		}

		prev = current
		current = current.nextNode
	}

	return ErrorElementNotRemoved
}

func (list *LinkedList[T]) RemoveFront() *T {
	element := list.head.element
	list.head = list.head.nextNode
	list.length--

	if list.head == nil {
		list.tail = nil
	}

	return element
}

func (list *LinkedList[T]) PrintList() {
	current := list.head
	for current != nil {
		fmt.Printf("%v ", current.element)
		current = current.nextNode
	}
	fmt.Printf("\n")
}

func (list *LinkedList[T]) Clear() {
	list.head = nil
	list.tail = nil
	list.length = 0
}
