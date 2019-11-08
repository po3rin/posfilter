package posfilter_test

import (
	"fmt"
	"testing"

	"github.com/po3rin/posfilter"
)

func equalSlice(tb *testing.T, want []string, got []string) error {
	tb.Helper()
	if len(want) != len(got) {
		return fmt.Errorf("unexpected length: want: %v, got: %v", want, got)
	}
	for i, w := range want {
		if w != got[i] {
			return fmt.Errorf("unexpected length: want: %v, got: %v", want, got)
		}
	}
	return nil
}

func TestPosFilter(t *testing.T) {
	tests := []struct {
		text string
		want []string
	}{
		{
			text: "東京都へ行く",
			want: []string{"東京都"},
		},
		{
			text: "スカイツリーには素晴らしいお店がある",
			want: []string{"スカイツリー", "お店"},
		},
		{
			text: "Rustはシステム言語なので、低水準の操作を行います。",
			want: []string{"Rust", "システム", "言語", "水準", "操作"},
		},
		{
			text: "gosudachiは日本語形態素解析器であるSudachiのGo移植版です。",
			want: []string{"gosudachi", "日本語", "形態素", "解析", "Sudachi", "Go", "移植"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.text, func(t *testing.T) {
			t.Parallel()

			filter, err := posfilter.NewPosFilter()
			defer filter.Close()
			if err != nil {
				t.Fatalf("unexpected fatal error: %v", err)
			}

			got, err := filter.Do(tt.text)
			if err != nil {
				t.Errorf("unexpected err: %v", err)
			}

			if err := equalSlice(t, tt.want, got); err != nil {
				t.Error(err)
			}
		})
	}
}
