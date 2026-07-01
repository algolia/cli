package interactive

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequiredString(t *testing.T) {
	v := requiredString("name")
	require.Error(t, v(""))
	require.NoError(t, v("x"))
}

func TestIntValidator(t *testing.T) {
	req := intValidator("count", 32, false)
	require.Error(t, req(""))     // required: empty rejected
	require.Error(t, req("abc"))  // not a number
	require.Error(t, req("3.5"))  // not a whole number
	require.NoError(t, req("42")) // ok

	opt := intValidator("max", 32, true)
	require.NoError(t, opt(""))   // optional: empty allowed
	require.Error(t, opt("abc"))  // still must parse when present
	require.NoError(t, opt("-7")) // ok

	// bit width is enforced.
	require.Error(t, intValidator("n", 32, true)("9999999999999")) // overflows int32
}

func TestFloatValidator(t *testing.T) {
	req := floatValidator("ratio", 64, false)
	require.Error(t, req(""))
	require.Error(t, req("abc"))
	require.NoError(t, req("1.5"))

	opt := floatValidator("ratio", 64, true)
	require.NoError(t, opt(""))
	require.Error(t, opt("x"))
	require.NoError(t, opt("2"))

	// 32-bit overflow is rejected at the field's precision.
	require.Error(t, floatValidator("ratio", 32, true)("1e40"))
	require.NoError(t, floatValidator("ratio", 64, true)("1e40"))
}

func TestCountValidator(t *testing.T) {
	v := countValidator()
	require.NoError(t, v("")) // none
	require.NoError(t, v("0"))
	require.NoError(t, v("3"))
	require.Error(t, v("-1")) // negative rejected
	require.Error(t, v("abc"))
}

func TestBoolValidator(t *testing.T) {
	v := boolValidator()
	require.NoError(t, v("")) // optional skip
	require.NoError(t, v("true"))
	require.NoError(t, v("false"))
	require.NoError(t, v("1"))
	require.Error(t, v("maybe"))
}

func TestJSONValidator(t *testing.T) {
	v := jsonValidator()
	require.NoError(t, v("")) // skip
	require.NoError(t, v(`{"a":1}`))
	require.NoError(t, v(`[1,2,3]`))
	require.Error(t, v(`{not json}`))
}

func TestValidatorMessagesAreUserFacing(t *testing.T) {
	// The message survey shows on retry should be terse, not a wrapped Go error.
	assert.EqualError(t, intValidator("count", 32, false)("abc"), "must be a whole number")
	assert.EqualError(t, boolValidator()("x"), "must be true or false")
	assert.EqualError(t, jsonValidator()("{"), "must be valid JSON")
}
