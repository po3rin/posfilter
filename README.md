# posfilter

<img src="https://img.shields.io/badge/go-v1.13-blue.svg"/>

this package lets you to tokenize & filter words by part of speech.
Depends on Sudachi tokenizer.

You have to download dictionary file.
https://github.com/WorksApplications/SudachiDict

## Quick start

```go
package main

import (
	// ...
	"github.com/po3rin/posfilter"
)

func main() {
    var filter posfilter.PosFilter
	words, err := filter.Do("gosudachiは日本語形態素解析器であるSudachiのGo移植版です。")
	if err != nil {
		log.Fatal(err)
	}
    fmt.Println(words)
    // [gosudachi 日本語 形態素 解析 Sudachi Go 移植]
}
```

default target pos is here.

```go
// in posfilter.go
var defaultTargetPos = map[string]struct{}{
	"名詞,普通名詞,一般":      struct{}{},
	"名詞,普通名詞,サ変可能":    struct{}{},
	"名詞,普通名詞,形状詞可能":   struct{}{},
	"名詞,普通名詞,サ変形状詞可能": struct{}{},
	"名詞,普通名詞,副詞可能":    struct{}{},
	"名詞,固有名詞,一般":      struct{}{},
	"名詞,固有名詞,人名":      struct{}{},
	"名詞,固有名詞,地名":      struct{}{},
	"名詞,固有名詞,組織名":     struct{}{},
}
```

you can use custom settings using builder pattern.

```go
// ...
words, _ := NewPosFilter().
    SetMode(ModeA).
    SetTargetPos([]string{"名詞,固有名詞,地名","名詞,固有名詞,組織名"}).
    SetSettingFilePath("custom_setting.json").
    Do("スカイツリーには素晴らしいお店がある")

// ...
```
