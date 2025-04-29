package entities

import "errors"

type Option[T any] struct {
	value    T
	hasValue bool
}

func EmptyOption[T any]() Option[T] {
	return Option[T]{
		value:    *new(T),
		hasValue: false,
	}
}

func OptionFrom[T any](value T) Option[T] {
	return Option[T]{
		value:    value,
		hasValue: true,
	}
}

func OptionFromNilable[T any](value *T) Option[T] {
	if value == nil {
		return EmptyOption[T]()
	}

	return OptionFrom[T](*value)
}

func (o Option[T]) HasValue() bool {
	return o.hasValue
}

func (o Option[T]) Value() (T, error) {
	if !o.hasValue {
		return o.value, errors.New("no value for option")
	}

	return o.value, nil
}

type ArticlePatch struct {
	Starred Option[bool]
	Read    Option[bool]
}

type FeedSourcePatch struct {
	Name     Option[string]
	Disabled Option[bool]
}
