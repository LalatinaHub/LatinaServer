package model

// KeyValue represents a generic key-value configuration entry.
type KeyValue struct {
	ID    int64
	Key   string
	Value any
}
