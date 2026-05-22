package utils_test

import (
	"testing"

	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidationErrorMessage(t *testing.T) {
	t.Run("RequiredErrMsgCheck", func(t *testing.T) {
		assert.Equal(
			t,
			"field1 is required.",
			utils.ValidationErrorMessage("field1", "required"),
		)
	})

	t.Run("InvalidUrlErrMsgCheck", func(t *testing.T) {
		assert.Equal(
			t,
			"field2 must be a valid URL.",
			utils.ValidationErrorMessage("field2", "url"),
		)
	})

	t.Run("LengthExceedErrMsgCheck", func(t *testing.T) {
		assert.Equal(
			t,
			"field3 exceeds maximum length.",
			utils.ValidationErrorMessage("field3", "max"),
		)
	})

	t.Run("OtherErrMsgCheck", func(t *testing.T) {
		assert.Equal(
			t,
			"field4 is invalid.",
			utils.ValidationErrorMessage("field4", "others"),
		)
	})
}
