package object

/*
組み込み関数の型
任意のオブジェクトを受け取り単一のオブジェクトを返す
インタプリタとホスト言語(Go)の受け渡しをする
*/
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) AsBool() bool     { return true }
