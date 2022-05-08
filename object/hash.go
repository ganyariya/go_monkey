package object

import (
	"fmt"
	"strings"
)

type HashKey struct {
	/*
		String or Integer or Boolean
		タイプごとに比較する
	*/
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	/*
		HashKey は Type と ハッシュ化された Key が入る（"Hello" -> 489281121）
		そのため HashPair にハッシュする前のオリジナルオブジェクト key: value を保存する
	*/
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
func (h *Hash) AsBool() bool { return len(h.Pairs) > 0 }
