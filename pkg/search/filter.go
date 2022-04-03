package search

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"
	"github.com/blugelabs/bluge/analysis/token"
	"github.com/blugelabs/bluge/analysis/tokenizer"
	"github.com/longbridgeapp/opencc"
)

type ConvertFilter struct {
	cc    *opencc.OpenCC
	fw2hw unicode.SpecialCase
}

func NewConvertFilter() (f *ConvertFilter, err error) {
	f = &ConvertFilter{}
	f.cc, err = opencc.New("t2s")
	f.fw2hw = unicode.SpecialCase{
		unicode.CaseRange{
			Lo: 0xff01,
			Hi: 0xff5e,
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x0021 - 0xff01, // LowerCase
				0,               // TitleCase
			},
		},
	}
	return
}

func (c *ConvertFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	for _, t := range input {
		convert, err := c.cc.Convert(string(t.Term))
		if err != nil {
			panic(err)
		}
		convert = strings.ToLowerSpecial(c.fw2hw, convert)
		t.Term = []byte(convert)
	}
	return input
}

var filenameAnalyzer *analysis.Analyzer

func init() {
	filter, err := NewConvertFilter()
	if err != nil {
		panic(err)
	}

	filenameAnalyzer = &analysis.Analyzer{
		CharFilters: []analysis.CharFilter{
			// UnicodeTokenizer 不会分 .
			char.NewRegexpCharFilter(regexp.MustCompile(`[.,_]+`), []byte(" ")),
		},
		Tokenizer: tokenizer.NewUnicodeTokenizer(),
		TokenFilters: []analysis.TokenFilter{
			token.NewLowerCaseFilter(),
			filter,
		},
	}
}
