package util

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	ALIGN_LEFT   = 0
	ALIGN_CENTER = 1
	ALIGN_RIGHT  = 2
)

func Pad(s, pad string, width int, align int) string {
	switch align {
	case ALIGN_CENTER:
		return PadCenter(s, pad, width)
	case ALIGN_RIGHT:
		return PadLeft(s, pad, width)
	default:
		return PadRight(s, pad, width)
	}
}

func PadRight(s, pad string, width int) string {
	gap := widthValue(s, width)
	if gap > 0 {
		return s + strings.Repeat(string(pad), gap)
	}
	return s
}

func PadLeft(s, pad string, width int) string {
	gap := widthValue(s, width)
	if gap > 0 {
		return strings.Repeat(string(pad), gap) + s
	}
	return s
}

func PadCenter(s, pad string, width int) string {
	gap := widthValue(s, width)
	if gap > 0 {
		gapLeft := int(math.Ceil(float64(gap / 2)))
		gapRight := gap - gapLeft
		return strings.Repeat(string(pad), gapLeft) + s + strings.Repeat(string(pad), gapRight)
	}
	return s
}

func isHan(s string) (isHan bool) {
	wh := []rune(s)
	for _, r := range wh {
		if isHan != unicode.Is(unicode.Han, r) {
			break
		}
	}
	return
}

func countCN(s string) (count int) {
	wh := []rune(s)
	for _, r := range wh {
		if unicode.Is(unicode.Han, r) {
			count++
		}
	}
	return
}

func widthValue(s string, width int) (gap int) {
	l := utf8.RuneCountInString(s)
	ln := len(s)
	isHan := isHan(s)
	count := countCN(s)
	if ln != l {
		if isHan {
			gap = width - (ln - l)
		} else {
			gap = width - (ln - count)
		}
	} else {
		gap = width - l
	}
	return
}
