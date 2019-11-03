package posfilter

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
