package bbrpc

import (
	"fmt"
	"strings"
	"testing"
)

func TestUtilDataEncoding(t *testing.T) {
	data := "hello 大棒客为武汉加油" //suffix: 68656c6c6f20e5a4a7e6a392e5aea2e4b8bae6ada6e6b189e58aa0e6b2b9
	encData := UtilDataEncoding(data)
	tShouldTrue(t, strings.HasSuffix(encData, "68656c6c6f20e5a4a7e6a392e5aea2e4b8bae6ada6e6b189e58aa0e6b2b9"), "should has suffix ...")

	for _, tt := range []string{
		"", "abc", "123", "中文", "含标点符号！!:@", "^_^", "武汉加油！", "新型コロナウイルス",
		"make USA great again", "Tumpang lalu", "해질녘,새벽녘", "emoji:😂😀", "(´▽｀)", "Italia rifornimento",
		"卢本伟🐂🍺",
	} {
		en := UtilDataEncoding(tt)
		tShouldTrue(t, len(en) >= 42, "len should greater than 42")
		de, er := UtilDataDecoding(en)
		tShouldNil(t, er)
		tShouldTrue(t, tt == de.Data, fmt.Sprintf("expected: %s, got: %s\n", tt, de.Data))
	}
}
