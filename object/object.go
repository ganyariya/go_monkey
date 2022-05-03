package object

/*
`AST ノード`をオブジェクトシステムの Object に変換する。
Object の Type を区別するためのタイプを定義する。
*/
const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type ObjectType string

type Object interface {
	Type() ObjectType // オブジェクトのタイプを表す
	Inspect() string  // オブジェクトがラップしている値を表す
}
