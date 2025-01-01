package utils

// IsWhiteSpace checks if the given byte represents a white space character
func IsWhiteSpace(ch byte) bool {
	switch ch {
	case 0x00, 0x09, 0x0A, 0x0C, 0x0D, 0x20:
		return true
	default:
		return false
	}
}

// IsFloatDigit checks if the given byte can be a part of a float number string.
func IsFloatDigit(c byte) bool {
	return ('0' <= c && c <= '9') || c == '.'
}

// IsDecimalDigit checks if the given byte is a part of a decimal number string.
func IsDecimalDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsOctalDigit checks if the given byte can be part of an octal digit string.
func IsOctalDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

// IsPrintable checks if the given byte is a printable character.
func IsPrintable(c byte) bool {
	return 0x21 <= c && c <= 0x7E
}

// IsDelimiter checks if the given byte represents a delimiter character
// according to the PDF specification.
func IsDelimiter(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	default:
		return false
	}
}
