package processor

import (
	"regexp"
	"strconv"
	"strings"
)

func Process(words []string) ([]string, error) {

	vowels := "aeiouhAEIOUH"

	for i := 0; i < len(words); i++ {

		//HEX
		if strings.Contains(words[i], "(hex)") && i > 0 {
			num, err := strconv.ParseInt(words[i-1], 16, 64)
			if err != nil {
				return nil, err
			}
			words[i-1] = strconv.Itoa(int(num))
			words[i] = strings.ReplaceAll(words[i], "(hex)", "")
		}

		//BIN
		if strings.Contains(words[i], "(bin)") && i > 0 {
			num, err := strconv.ParseInt(words[i-1], 2, 64)
			if err != nil {
				return nil, err
			}
			words[i-1] = strconv.Itoa(int(num))
			words[i] = strings.ReplaceAll(words[i], "(bin)", "")
		}

		//SINGLE CASE
		if strings.Contains(words[i], "(up)") && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1])
			words[i] = strings.ReplaceAll(words[i], "(up)", "")
		}

		if strings.Contains(words[i], "(low)") && i > 0 {
			words[i-1] = strings.ToLower(words[i-1])
			words[i] = strings.ReplaceAll(words[i], "(low)", "")
		}

		if strings.Contains(words[i], "(cap)") && i > 0 {
			w := words[i-1]
			if len(w) > 0 {
				words[i-1] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
			}
			words[i] = strings.ReplaceAll(words[i], "(cap)", "")
		}

		//MULTI CASE (UP)
		if strings.Contains(words[i], "(up,") && i+1 < len(words) {

			part := words[i] + words[i+1]
			start := strings.Index(part, ",")
			end := strings.Index(part, ")")

			if start != -1 && end != -1 {
				numStr := strings.TrimSpace(part[start+1 : end])
				n, err := strconv.Atoi(numStr)

				if err == nil {
					startIdx := i - n
					if startIdx < 0 {
						startIdx = 0
					}

					for j := startIdx; j < i; j++ {
						words[j] = strings.ToUpper(words[j])
					}
				}
			}

			words[i] = ""
			words[i+1] = ""
		}

		//MULTI CASE (LOW)
		if strings.Contains(words[i], "(low,") && i+1 < len(words) {

			part := words[i] + words[i+1]

			start := strings.Index(part, ",")
			end := strings.Index(part, ")")

			if start != -1 && end != -1 {
				numStr := strings.TrimSpace(part[start+1 : end])
				n, err := strconv.Atoi(numStr)

				if err == nil {
					startIdx := i - n
					if startIdx < 0 {
						startIdx = 0
					}

					for j := startIdx; j < i; j++ {
						words[j] = strings.ToLower(words[j])
					}
				}
			}

			words[i] = ""
			words[i+1] = ""
		}

		//MULTI CASE (CAP)
		if strings.Contains(words[i], "(cap,") && i+1 < len(words) {

			part := words[i] + words[i+1]

			start := strings.Index(part, ",")
			end := strings.Index(part, ")")

			if start != -1 && end != -1 {
				numStr := strings.TrimSpace(part[start+1 : end])
				n, err := strconv.Atoi(numStr)

				if err == nil {
					startIdx := i - n
					if startIdx < 0 {
						startIdx = 0
					}

					for j := startIdx; j < i; j++ {
						if len(words[j]) > 0 {
							w := words[j]
							words[j] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
						}
					}
				}
			}

			words[i] = ""
			words[i+1] = ""
		}
	}

	//A → AN
	for i := 0; i < len(words)-1; i++ {
		if strings.ToLower(words[i]) == "a" {
			next := words[i+1]
			if len(next) > 0 && strings.ContainsRune(vowels, rune(next[0])) {
				if words[i] == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}

	return words, nil
}

func Clean(words []string) []string {
	var result []string
	for _, w := range words {
		if strings.TrimSpace(w) != "" {
			result = append(result, w)
		}
	}
	return result
}

func FixQuotes(text string) string {

	// remove spaces inside quotes
	re := regexp.MustCompile(`'\s*(.*?)\s*'`)
	text = re.ReplaceAllString(text, "'$1'")

	return text
}

func FixPunctuation(text string) string {

	// remove space BEFORE punctuation
	re := regexp.MustCompile(`\s+([.,!?;:])`)
	text = re.ReplaceAllString(text, "$1")

	// ensure ONE space AFTER punctuation (except end)
	re = regexp.MustCompile(`([.,!?;:])([^\s])`)
	text = re.ReplaceAllString(text, "$1 $2")

	// fix multiple dots (.... . . -> .......)
	re = regexp.MustCompile(`(\.\s*)+`)
	text = re.ReplaceAllStringFunc(text, func(s string) string {
		count := strings.Count(s, ".")
		return strings.Repeat(".", count)
	})

	// remove multiple spaces
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}
