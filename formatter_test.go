package main

import "testing"

func TestFormat(t *testing.T) {
	tests := []struct {
		name, input, want string
	}{
		{"hex", "1E (hex) files were added", "30 files were added"},
		{"bin", "It has been 10 (bin) years", "It has been 2 years"},
		{"up", "Ready, set, go (up) !", "Ready, set, GO!"},
		{"low", "I should stop SHOUTING (low)", "I should stop shouting"},
		{"cap", "Welcome to the Brooklyn bridge (cap)", "Welcome to the Brooklyn Bridge"},
		{"count", "This is so exciting (up, 2)", "This is SO EXCITING"},
		{"low count", "ONE TWO THREE (low,2)", "ONE two three"},
		{"cap count", "it was the age of foolishness (cap, 6)", "It Was The Age Of Foolishness"},
		{"punctuation", "I was sitting over there ,and then BAMM !!", "I was sitting over there, and then BAMM!!"},
		{"ellipsis", "I was thinking ... You were right", "I was thinking... You were right"},
		{"clusters", "Wait ! ? Really :yes ;no .", "Wait!? Really: yes; no."},
		{"quote", "I am exactly how they describe me: ' awesome '", "I am exactly how they describe me: 'awesome'"},
		{"quote phrase", "As Elton John said: ' I am the most well-known homosexual in the world '", "As Elton John said: 'I am the most well-known homosexual in the world'"},
		{"attached quotes", "'hello' and ' world '!", "'hello' and 'world'!"},
		{"articles", "a apple, a elephant, a idea, a orange, a umbrella, a house, a cat", "an apple, an elephant, an idea, an orange, an umbrella, an house, a cat"},
		{"capital article", "There it was. A amazing rock!", "There it was. An amazing rock!"},
		{"numbers sample", "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.", "Simply add 66 and 2 and you will see the result is 68."},
		{"article sample", "There is no greater agony than bearing a untold story inside you.", "There is no greater agony than bearing an untold story inside you."},
		{"punctuation sample", "Punctuation tests are ... kinda boring ,what do you think ?", "Punctuation tests are... kinda boring, what do you think?"},
		{"chained", "hELLO (up) (low) (cap)", "Hello"},
		{"punctuation before command", "hello, world! (up, 2)", "HELLO, WORLD!"},
		{"count exceeds words", "one two (up, 20)", "ONE TWO"},
		{"article after command", "a APPLE (low)", "an apple"},
		{"unicode", "élan über (cap, 2)", "Élan Über"},
		{"contractions", "don't shout (up), I'm here.", "don't SHOUT, I'm here."},
		{"large integer", "FFFFFFFFFFFFFFFFFFFFFFFF (hex)", "79228162514264337593543950335"},
		{"whitespace", "  hello\n\tworld  !  ", "hello world!"},
		{"empty", " \n\t", ""},
		{"ordinary parentheses", "hello (world)", "hello (world)"},
		{"article sentence boundary", "a. apple", "a. apple"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("Format(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatErrors(t *testing.T) {
	for _, input := range []string{"(up)", "xyz (hex)", "102 (bin)", "hello (up, 0)", "10 (bin, 2)", "word (low, 999999999999999999999999999)"} {
		if _, err := Format(input); err == nil {
			t.Errorf("Format(%q) should fail", input)
		}
	}
}
