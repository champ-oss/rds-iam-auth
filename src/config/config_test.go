package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_ParseBool_true(t *testing.T) {
	_ = os.Setenv("FOO", "true")
	assert.True(t, parseBool("FOO", false))
	_ = os.Setenv("FOO", "True")
	assert.True(t, parseBool("FOO", false))
}

func Test_ParseBool_false(t *testing.T) {
	_ = os.Setenv("FOO", "false")
	assert.False(t, parseBool("FOO", true))
	_ = os.Setenv("FOO", "False")
	assert.False(t, parseBool("FOO", true))
}

func Test_ParseBool_invalid(t *testing.T) {
	_ = os.Setenv("FOO", "nonsense")
	assert.False(t, parseBool("FOO", false))
}

func Test_ParseBool_unset(t *testing.T) {
	os.Clearenv()
	assert.False(t, parseBool("FOO", false))
	os.Clearenv()
	assert.True(t, parseBool("FOO", true))
}

func Test_ParseString_success(t *testing.T) {
	_ = os.Setenv("FOO", "stuff")
	assert.Equal(t, "stuff", parseString("FOO", ""))
}

func Test_ParseString_empty(t *testing.T) {
	_ = os.Setenv("FOO", "")
	assert.Equal(t, "something", parseString("FOO", "something"))
}

func Test_ParseString_unset(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, "something", parseString("FOO", "something"))
}

func Test_parseCommaSeparated_success(t *testing.T) {
	_ = os.Setenv("FOO", "foo,bar")
	assert.Equal(t, []string{"foo", "bar"}, parseCommaSeparated("FOO", []string{"something"}))
}

func Test_parseCommaSeparated_empty(t *testing.T) {
	_ = os.Setenv("FOO", "")
	assert.Equal(t, []string{"foo", "bar"}, parseCommaSeparated("FOO", []string{"foo", "bar"}))
}

func Test_parseCommaSeparated_unset(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, []string{"foo", "bar"}, parseCommaSeparated("FOO", []string{"foo", "bar"}))
}

func Test_LoadConfig(t *testing.T) {
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("AWS_REGION", "us-west-1")
	_ = os.Setenv("QUEUE_URL", "http://localhost")

	config := LoadConfig()
	assert.Equal(t, true, config.Debug)
	assert.Equal(t, "us-west-1", config.AwsRegion)
	assert.Equal(t, "http://localhost", config.QueueUrl)
}

func Test_getAWSConfig(t *testing.T) {
	assert.NotNil(t, getAWSConfig("us-east-1"))
}

func Test_setLogging(t *testing.T) {
	assert.NotPanics(t, func() {
		setLogging(true)
	})
}
