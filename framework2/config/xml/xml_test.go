package xml

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigContainer_Sub(t *testing.T) {
	tests := []struct {
		name           string
		initialData    map[string]any
		key            string
		expectedResult map[string]any
		expectError    bool
	}{
		{
			name: "Existing nested map",
			initialData: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": "5432",
				},
			},
			key: "database",
			expectedResult: map[string]any{
				"host": "localhost",
				"port": "5432",
			},
			expectError: false,
		},
		{
			name:           "Empty key returns full data",
			initialData:    map[string]any{"key": "value"},
			key:            "",
			expectedResult: map[string]any{"key": "value"},
			expectError:    false,
		},
		{
			name:           "Non-existent key",
			initialData:    map[string]any{},
			key:            "missing",
			expectedResult: nil,
			expectError:    true,
		},
		{
			name: "Non-map value",
			initialData: map[string]any{
				"simple": "string",
			},
			key:            "simple",
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigContainer{
				data: tt.initialData,
			}

			result, err := c.Sub(tt.key)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				subContainer, ok := result.(*ConfigContainer)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedResult, subContainer.data)
			}
		})
	}
}

func TestConfigContainer_SetAndGet(t *testing.T) {
	c := &ConfigContainer{
		data: make(map[string]any),
	}

	// Test Set and String
	err := c.Set("name", "test")
	assert.NoError(t, err)

	val, err := c.String("name")
	assert.NoError(t, err)
	assert.Equal(t, "test", val)

	// Test Strings
	err = c.Set("tags", "go;test;config")
	assert.NoError(t, err)

	tags, err := c.Strings("tags")
	assert.NoError(t, err)
	assert.Equal(t, []string{"go", "test", "config"}, tags)
}

func TestConfigContainer_TypeConversions(t *testing.T) {
	c := &ConfigContainer{
		data: map[string]any{
			"age":      "30",
			"bignum":   "9223372036854775807",
			"isActive": "true",
			"price":    "99.99",
		},
	}

	// Test Int
	age, err := c.Int("age")
	assert.NoError(t, err)
	assert.Equal(t, 30, age)

	// Test Int64
	bignum, err := c.Int64("bignum")
	assert.NoError(t, err)
	assert.Equal(t, int64(9223372036854775807), bignum)

	// Test Bool
	isActive, err := c.Bool("isActive")
	assert.NoError(t, err)
	assert.True(t, isActive)

	// Test Float
	price, err := c.Float("price")
	assert.NoError(t, err)
	assert.Equal(t, 99.99, price)
}

func TestConfigContainer_DefaultValues(t *testing.T) {
	c := &ConfigContainer{
		data: map[string]any{
			"existing_string": "value",
			"existing_int":    "42",
		},
	}

	// Test DefaultString
	assert.Equal(t, "default", c.DefaultString("missing", "default"))
	assert.Equal(t, "value", c.DefaultString("existing_string", "default"))

	// Test DefaultStrings
	defaultTags := []string{"default"}
	assert.Equal(t, defaultTags, c.DefaultStrings("missing", defaultTags))
	assert.Equal(t, []string{"go", "test"}, c.DefaultStrings("missing", []string{"go", "test"}))

	// Test DefaultInt
	assert.Equal(t, 100, c.DefaultInt("missing", 100))
	assert.Equal(t, 42, c.DefaultInt("existing_int", 100))

	// Test DefaultBool
	c.data["existing_bool"] = "true"
	assert.Equal(t, true, c.DefaultBool("missing", true))
	assert.Equal(t, true, c.DefaultBool("existing_bool", false))

	// Test DefaultFloat
	c.data["existing_float"] = "3.14"
	assert.Equal(t, 2.0, c.DefaultFloat("missing", 2.0))
	assert.Equal(t, 3.14, c.DefaultFloat("existing_float", 2.0))
}

func TestConfigContainer_ConcurrentSet(t *testing.T) {
	c := &ConfigContainer{
		data: make(map[string]any),
	}

	// 并发写入测试
	for i := 0; i < 100; i++ {
		go func(n int) {
			err := c.Set(fmt.Sprintf("key_%d", n), fmt.Sprintf("value_%d", n))
			assert.NoError(t, err)
		}(i)
	}
}
