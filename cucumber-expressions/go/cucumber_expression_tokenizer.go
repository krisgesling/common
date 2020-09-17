package cucumberexpressions

func tokenize(expression string) ([]token, error) {

	tokens := make([]token, 0)

	runes := []rune(expression)

	var buffer []rune
	previousTokenType := startOfLine
	treatAsText := false
	escaped := 0
	bufferStartIndex := 0

	convertBufferToToken := func(tokenType tokenType) token {
		escapeTokens := 0
		if tokenType == text {
			escapeTokens = escaped
			escaped = 0
		}
		consumedIndex := bufferStartIndex + len(buffer) + escapeTokens
		t := token{string(buffer), tokenType, bufferStartIndex, consumedIndex}
		buffer = []rune{}
		bufferStartIndex = consumedIndex
		return t
	}

	tokens = append(tokens, token{"", startOfLine, 0, 0})

	for _, rune := range runes {
		if !treatAsText && isEscapeCharacter(rune) {
			escaped++
			treatAsText = true
			continue
		}

		currentTokenType, err := tokenTypeOf(rune, treatAsText)
		if err != nil {
			return nil, err
		}
		treatAsText = false

		if previousTokenType != startOfLine && (currentTokenType != previousTokenType || (currentTokenType != whiteSpace && currentTokenType != text)) {
			token := convertBufferToToken(previousTokenType)
			previousTokenType = currentTokenType
			buffer = append(buffer, rune)
			tokens = append(tokens, token)
		} else {
			previousTokenType = currentTokenType
			buffer = append(buffer, rune)
		}

	}

	if len(buffer) > 0 {
		token := convertBufferToToken(previousTokenType)
		tokens = append(tokens, token)
	}

	if treatAsText {
		return nil, NewCucumberExpressionError("can't escape EOL")
	}
	token := token{"", endOfLine, len(runes), len(runes)}
	tokens = append(tokens, token)
	return tokens, nil
}

func tokenTypeOf(r rune, treatAsText bool) (tokenType, error) {
	if !treatAsText {
		return typeOf(r)
	}
	if canEscape(r) {
		return text, nil
	}
	return startOfLine, NewCucumberExpressionError("can't escape")

}
