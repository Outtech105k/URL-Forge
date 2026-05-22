package utils

// JSON validatorが発生させたエラーからメッセージを作成
func ValidationErrorMessage(field, tag string) string {
	switch tag {
	case "required":
		return field + " is required."
	case "url":
		return field + " must be a valid URL."
	case "max":
		return field + " exceeds maximum length."
	default:
		return field + " is invalid."
	}
}
