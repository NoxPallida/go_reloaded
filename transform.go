package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	reBeforePunct = regexp.MustCompile(`\s+([.,!?:;]+)`)
	reAfterPunct  = regexp.MustCompile(`([.,!?:;]+)([^\s.,!?:;'])`)
)

func handleHexBin(words []string) []string {
	for i := 1; i < len(words); i++ {
		var base int
		switch words[i] {
		case "(hex)":
			base = 16
		case "(bin)":
			base = 2
		default:
			continue
		}

		val, err := strconv.ParseInt(words[i-1], base, 64)
		if err != nil {
			continue
		}

		words[i-1] = strconv.FormatInt(val, 10)
		words = append(words[:i], words[i+1:]...)
		i--
	}
	return words
}

func handleUpLowCap(words []string) []string {
	for i := 0; i < len(words); i++ {
		var convert func(string) string
		var tag string

		switch words[i] {
		case "(up)", "(up,":
			convert, tag = strings.ToUpper, "(up"
		case "(low)", "(low,":
			convert, tag = strings.ToLower, "(low"
		case "(cap)", "(cap,":
			tag = "(cap"
			convert = func(s string) string {
				r := []rune(strings.ToLower(s))
				if len(r) > 0 {
					r[0] = unicode.ToUpper(r[0])
				}
				return string(r)
			}
		default:
			continue
		}

		count, size := 1, 1

		if words[i] == tag+"," {
			if i+1 >= len(words) {
				continue
			}
			val, err := strconv.Atoi(strings.TrimSuffix(words[i+1], ")"))
			if err != nil {
				continue
			}
			count, size = val, 2
		}

		for j := 1; j <= count && i-j >= 0; j++ {
			words[i-j] = convert(words[i-j])
		}

		words = append(words[:i], words[i+size:]...)
		i--
	}
	return words
}

func ProcessText(input string) string {
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		words := handleModifiers(strings.Fields(line))

		text := strings.Join(words, " ")
		text = fixQuotes(text)
		text = fixPunctuation(text)

		lines[i] = strings.Join(handleArticles(strings.Fields(text)), " ")
	}
	return strings.Join(lines, "\n")
}

func handleModifiers(words []string) []string {
	words = handleHexBin(words)
	words = handleUpLowCap(words)
	return words
}

func handleArticles(words []string) []string {
	for i := 0; i < len(words)-1; i++ {
		if words[i] != "a" && words[i] != "A" {
			continue
		}

		first := unicode.ToLower([]rune(words[i+1])[0])
		if strings.ContainsRune("aeiouh", first) {
			words[i] += "n"
		}
	}
	return words
}

func fixPunctuation(text string) string {
	text = reBeforePunct.ReplaceAllString(text, "$1")
	return reAfterPunct.ReplaceAllString(text, "$1 $2")
}

func fixQuotes(text string) string {
	words := strings.Fields(text)
	result := make([]string, 0, len(words))
	open, pending := false, false

	for _, w := range words {
		switch {
		case w != "'":
			if pending {
				w = "'" + w
				pending = false
			}
			result = append(result, w)
		case !open:
			open, pending = true, true
		case len(result) > 0:
			result[len(result)-1] += "'"
			open = false
		}
	}
	return strings.Join(result, " ")
}
