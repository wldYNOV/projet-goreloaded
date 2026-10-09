package main

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type tokenKind int

const (
	word tokenKind = iota
	punctuation
	quote
	directive
)

type token struct {
	text string
	kind tokenKind
}

var commandPattern = regexp.MustCompile(`^\((hex|bin|up|low|cap)(?:\s*,\s*([0-9]+))?\)$`)

func Format(input string) (string, error) {
	tokens, err := applyCommands(tokenize(input))
	if err != nil {
		return "", err
	}
	fixArticles(tokens)
	return render(tokens), nil
}

func isPunctuation(r rune) bool {
	return strings.ContainsRune(".,!?:;", r)
}

func internalApostrophe(rs []rune, i int) bool {
	return rs[i] == '\'' && i > 0 && i+1 < len(rs) &&
		unicode.IsLetter(rs[i-1]) && unicode.IsLetter(rs[i+1])
}

func tokenize(input string) []token {
	rs := []rune(input)
	var out []token
	for i := 0; i < len(rs); {
		if unicode.IsSpace(rs[i]) {
			i++
			continue
		}
		if rs[i] == '(' {
			j := i + 1
			for j < len(rs) && rs[j] != ')' {
				j++
			}
			if j < len(rs) {
				candidate := string(rs[i : j+1])
				if commandPattern.MatchString(candidate) {
					out = append(out, token{candidate, directive})
					i = j + 1
					continue
				}
			}
		}
		if isPunctuation(rs[i]) {
			out = append(out, token{string(rs[i]), punctuation})
			i++
			continue
		}
		if rs[i] == '\'' && !internalApostrophe(rs, i) {
			out = append(out, token{"'", quote})
			i++
			continue
		}
		start := i
		i++
		for i < len(rs) && !unicode.IsSpace(rs[i]) && !isPunctuation(rs[i]) &&
			(rs[i] != '\'' || internalApostrophe(rs, i)) && rs[i] != '(' {
			i++
		}
		out = append(out, token{string(rs[start:i]), word})
	}
	return out
}

func applyCommands(tokens []token) ([]token, error) {
	out := make([]token, 0, len(tokens))
	var words []int
	for _, t := range tokens {
		if t.kind != directive {
			out = append(out, t)
			if t.kind == word {
				words = append(words, len(out)-1)
			}
			continue
		}
		parts := commandPattern.FindStringSubmatch(t.text)
		command, count := parts[1], 1
		if parts[2] != "" {
			if command == "hex" || command == "bin" {
				return nil, fmt.Errorf("%s does not accept a word count", t.text)
			}
			var err error
			count, err = strconv.Atoi(parts[2])
			if err != nil || count < 1 {
				return nil, fmt.Errorf("invalid word count in %s", t.text)
			}
		}
		if len(words) == 0 {
			return nil, fmt.Errorf("%s requires a preceding word", t.text)
		}
		if count > len(words) {
			count = len(words)
		}
		for _, index := range words[len(words)-count:] {
			s := out[index].text
			switch command {
			case "up":
				out[index].text = strings.ToUpper(s)
			case "low":
				out[index].text = strings.ToLower(s)
			case "cap":
				rs := []rune(strings.ToLower(s))
				for i, r := range rs {
					if unicode.IsLetter(r) {
						rs[i] = unicode.ToUpper(r)
						break
					}
				}
				out[index].text = string(rs)
			case "hex", "bin":
				base := 16
				if command == "bin" {
					base = 2
				}
				n, ok := new(big.Int).SetString(s, base)
				if !ok {
					return nil, fmt.Errorf("invalid base-%d number %q before %s", base, s, t.text)
				}
				out[index].text = n.String()
			}
		}
	}
	return out, nil
}

func fixArticles(tokens []token) {
	for i, t := range tokens {
		if t.kind != word || (t.text != "a" && t.text != "A") {
			continue
		}
		j := i + 1
		for j < len(tokens) && tokens[j].kind == quote {
			j++
		}
		if j == len(tokens) || tokens[j].kind != word {
			continue
		}
		rs := []rune(tokens[j].text)
		if len(rs) > 0 && strings.ContainsRune("aeiouh", unicode.ToLower(rs[0])) {
			if t.text == "A" {
				tokens[i].text = "An"
			} else {
				tokens[i].text = "an"
			}
		}
	}
}

func render(tokens []token) string {
	var out strings.Builder
	insideQuote, previousOpening := false, false
	for i, t := range tokens {
		opening := t.kind == quote && !insideQuote
		attach := t.kind == punctuation || (t.kind == quote && !opening) || previousOpening
		if i > 0 && !attach {
			out.WriteByte(' ')
		}
		out.WriteString(t.text)
		if t.kind == quote {
			insideQuote = !insideQuote
		}
		previousOpening = opening
	}
	return out.String()
}
