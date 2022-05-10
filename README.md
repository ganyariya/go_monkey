# 🐒 go_monkey

[`Go 言語でつくるインタプリタ`](https://www.amazon.co.jp/Go%E8%A8%80%E8%AA%9E%E3%81%A7%E3%81%A4%E3%81%8F%E3%82%8B%E3%82%A4%E3%83%B3%E3%82%BF%E3%83%97%E3%83%AA%E3%82%BF-Thorsten-Ball/dp/4873118220) の ganyariya の実装です。


[Zenn Scrap](https://zenn.dev/ganariya/scraps/ce02513c03b094)

# 🐒 実行方法

```shell
git clone https://github.com/ganyariya/go_monkey.git
cd go_monkey
go run ./main.go
```

```txt
❯ go run ./main.go
Hello ganariya! This is the Monkey Programming Language!
>> let adder = fn (x) { fn (y) { x + y }};    
fn(x) {
fn(y)(x + y)
}
>> let two = adder(2)
fn(y) {
(x + y)
}
>> two(10)
12
```

# 🐒 History

- 1 章: [字句解析](https://github.com/ganyariya/go_monkey/tree/d9d62d8f28704a1d2c5655757c3d898e1fc95069)
- 2 章: [構文解析](https://github.com/ganyariya/go_monkey/tree/7b5e3786ae233e183379c3f46b9d7f35c5383dae)
- 3 章: [評価](https://github.com/ganyariya/go_monkey/tree/161db914a9c5092de4a26367578e3a5bcb5edefa)
- 4 章: [インタプリタの拡張](https://github.com/ganyariya/go_monkey/tree/964938bc166be1145b265a34a4c38d6531beb9f0)
- 付録:

