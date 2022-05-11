package object

/*
`AST ノード`をオブジェクトシステムの Object に変換する。
Object の Type を区別するためのタイプを定義する。
*/
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	QUOTE_OBJ        = "QUOTE"
)

type ObjectType string

type Object interface {
	Type() ObjectType // オブジェクトのタイプを表す
	Inspect() string  // オブジェクトがラップしている値を表す
	AsBool() bool     // オブジェクトの真偽値
}

/*
ハッシュ可能オブジェクトの共通インターフェース
*/
type Hashable interface {
	HashKey() HashKey
}
