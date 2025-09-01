package logs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
)

// TestLevel_String 测试 Level 的字符串表示
func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARNING"},
		{LevelError, "ERROR"},
		{LevelPanic, "PANIC"},
		{LevelFatal, "FATAL"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			if got := LevelNameMapping[tt.level]; got != tt.expected {
				t.Errorf("LevelNameMapping[%d] = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}

// TestLevel_UnmarshalText 测试 Level 的文本反序列化
func TestLevel_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
		wantErr  bool
	}{
		{"debug lowercase", "debug", LevelDebug, false},
		{"info lowercase", "info", LevelInfo, false},
		{"warning lowercase", "warning", LevelWarn, false},
		{"error lowercase", "error", LevelError, false},
		{"panic lowercase", "panic", LevelPanic, false},
		{"fatal lowercase", "fatal", LevelFatal, false},
		{"DEBUG uppercase", "DEBUG", LevelDebug, false},
		{"INFO uppercase", "INFO", LevelInfo, false},
		{"WARNING uppercase", "WARNING", LevelWarn, false},
		{"ERROR uppercase", "ERROR", LevelError, false},
		{"PANIC uppercase", "PANIC", LevelPanic, false},
		{"FATAL uppercase", "FATAL", LevelFatal, false},
		{"Mixed case", "DeBuG", LevelDebug, false},
		{"invalid level", "invalid", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l Level
			err := l.UnmarshalText([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalText(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("UnmarshalText(%q) unexpected error: %v", tt.input, err)
				}
				if l != tt.expected {
					t.Errorf("UnmarshalText(%q) = %v, want %v", tt.input, l, tt.expected)
				}
			}
		})
	}
}

// TestLevel_UnmarshalText_NilReceiver 测试空指针情况
func TestLevel_UnmarshalText_NilReceiver(t *testing.T) {
	var l *Level
	err := l.UnmarshalText([]byte("debug"))

	if err == nil {
		t.Error("Expected error when unmarshaling to nil Level pointer")
	}

	if !errors.Is(err, errUnmarshalNilLevel) {
		t.Errorf("Expected errUnmarshalNilLevel, got %v", err)
	}
}

// TestLevel_unmarshalText 测试私有方法 unmarshalText
func TestLevel_unmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
		success  bool
	}{
		{"valid debug", "debug", LevelDebug, true},
		{"valid info", "info", LevelInfo, true},
		{"valid warning", "warning", LevelWarn, true},
		{"valid error", "error", LevelError, true},
		{"valid panic", "panic", LevelPanic, true},
		{"valid fatal", "fatal", LevelFatal, true},
		{"invalid level", "unknown", 0, false},
		{"empty string", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l Level
			success := l.unmarshalText([]byte(tt.input))

			if success != tt.success {
				t.Errorf("unmarshalText(%q) success = %v, want %v", tt.input, success, tt.success)
			}

			if tt.success && l != tt.expected {
				t.Errorf("unmarshalText(%q) = %v, want %v", tt.input, l, tt.expected)
			}
		})
	}
}

// TestInitOptions 测试选项初始化
func TestInitOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		opts := initOptions()

		if opts.output != os.Stderr {
			t.Errorf("Default output = %v, want %v", opts.output, os.Stderr)
		}

		if opts.formatter == nil {
			t.Error("Default formatter should not be nil")
		}

		if _, ok := opts.formatter.(*TextFormatter); !ok {
			t.Errorf("Default formatter type = %T, want *TextFormatter", opts.formatter)
		}

		if opts.level != 0 {
			t.Errorf("Default level = %v, want 0", opts.level)
		}

		if opts.stdLevel != 0 {
			t.Errorf("Default stdLevel = %v, want 0", opts.stdLevel)
		}

		if opts.disableCaller != false {
			t.Errorf("Default disableCaller = %v, want false", opts.disableCaller)
		}
	})

	t.Run("with custom options", func(t *testing.T) {
		var buf bytes.Buffer
		jsonFormatter := &JSONFormatter{}

		opts := initOptions(
			WithOutput(&buf),
			WithLevel(LevelInfo),
			WithStdLevel(LevelError),
			WithFormatter(jsonFormatter),
			WithDisableCaller(true),
		)

		if opts.output != &buf {
			t.Errorf("Output = %v, want %v", opts.output, &buf)
		}

		if opts.level != LevelInfo {
			t.Errorf("Level = %v, want %v", opts.level, LevelInfo)
		}

		if opts.stdLevel != LevelError {
			t.Errorf("StdLevel = %v, want %v", opts.stdLevel, LevelError)
		}

		if opts.formatter != jsonFormatter {
			t.Errorf("Formatter = %v, want %v", opts.formatter, jsonFormatter)
		}

		if opts.disableCaller != true {
			t.Errorf("DisableCaller = %v, want true", opts.disableCaller)
		}
	})
}

// TestWithOutput 测试输出选项
func TestWithOutput(t *testing.T) {
	var buf bytes.Buffer
	opt := WithOutput(&buf)

	opts := &options{}
	opt(opts)

	if opts.output != &buf {
		t.Errorf("WithOutput() set output = %v, want %v", opts.output, &buf)
	}
}

// TestWithLevel 测试日志级别选项
func TestWithLevel(t *testing.T) {
	opt := WithLevel(LevelError)

	opts := &options{}
	opt(opts)

	if opts.level != LevelError {
		t.Errorf("WithLevel() set level = %v, want %v", opts.level, LevelError)
	}
}

// TestWithStdLevel 测试标准日志级别选项
func TestWithStdLevel(t *testing.T) {
	opt := WithStdLevel(LevelFatal)
	opts := &options{}
	opt(opts)

	if opts.stdLevel != LevelFatal {
		t.Errorf("WithStdLevel() set stdLevel = %v, want %v", opts.stdLevel, LevelFatal)
	}
}

// TestWithFormatter 测试格式化器选项
func TestWithFormatter(t *testing.T) {
	jsonFormatter := &JSONFormatter{}
	opt := WithFormatter(jsonFormatter)

	opts := &options{}
	opt(opts)

	if opts.formatter != jsonFormatter {
		t.Errorf("WithFormatter() set formatter = %v, want %v", opts.formatter, jsonFormatter)
	}
}

// TestWithDisableCaller 测试禁用调用者选项
func TestWithDisableCaller(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected bool
	}{
		{"disable caller true", true, true},
		{"disable caller false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDisableCaller(tt.input)

			opts := &options{}
			opt(opts)

			if opts.disableCaller != tt.expected {
				t.Errorf("WithDisableCaller(%v) set disableCaller = %v, want %v",
					tt.input, opts.disableCaller, tt.expected)
			}
		})
	}
}

// TestMultipleOptions 测试多个选项组合
func TestMultipleOptions(t *testing.T) {
	var buf bytes.Buffer
	jsonFormatter := &JSONFormatter{}

	opts := initOptions(
		WithOutput(&buf),
		WithLevel(LevelWarn),
		WithStdLevel(LevelPanic),
		WithFormatter(jsonFormatter),
		WithDisableCaller(true),
	)

	// 验证所有选项都被正确设置
	if opts.output != &buf {
		t.Errorf("Output not set correctly")
	}
	if opts.level != LevelWarn {
		t.Errorf("Level not set correctly")
	}
	if opts.stdLevel != LevelPanic {
		t.Errorf("StdLevel not set correctly")
	}
	if opts.formatter != jsonFormatter {
		t.Errorf("Formatter not set correctly")
	}
	if opts.disableCaller != true {
		t.Errorf("DisableCaller not set correctly")
	}
}

// TestLevelConstants 测试级别常量的值
func TestLevelConstants(t *testing.T) {
	expectedValues := map[Level]int8{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
		LevelPanic: 4,
		LevelFatal: 5,
	}

	for level, expectedValue := range expectedValues {
		if int8(level) != expectedValue {
			t.Errorf("Level %v = %d, want %d", level, int8(level), expectedValue)
		}
	}
}

// TestErrorMessages 测试错误信息
func TestErrorMessages(t *testing.T) {
	t.Run("errUnmarshalNilLevel", func(t *testing.T) {
		expectedMsg := "can't unmarshal a nil *Level"
		if errUnmarshalNilLevel.Error() != expectedMsg {
			t.Errorf("errUnmarshalNilLevel.Error() = %q, want %q",
				errUnmarshalNilLevel.Error(), expectedMsg)
		}
	})

	t.Run("unrecognized level error", func(t *testing.T) {
		var l Level
		err := l.UnmarshalText([]byte("invalid"))

		expectedMsg := `unrecognized level:"invalid"`
		if err.Error() != expectedMsg {
			t.Errorf("UnmarshalText error = %q, want %q", err.Error(), expectedMsg)
		}
	})
}

// BenchmarkLevel_UnmarshalText 性能基准测试
func BenchmarkLevel_UnmarshalText(b *testing.B) {
	testCases := []string{"debug", "info", "warning", "error", "panic", "fatal"}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			text := []byte(tc)
			for i := 0; i < b.N; i++ {
				var l Level
				l.UnmarshalText(text)
			}
		})
	}
}

// BenchmarkInitOptions 初始化选项性能测试
func BenchmarkInitOptions(b *testing.B) {
	var buf bytes.Buffer
	jsonFormatter := &JSONFormatter{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		initOptions(
			WithOutput(&buf),
			WithLevel(LevelInfo),
			WithStdLevel(LevelError),
			WithFormatter(jsonFormatter),
			WithDisableCaller(true),
		)
	}
}

// TestLevelNameMappingCompleteness 确保所有级别都有对应的名称映射
func TestLevelNameMappingCompleteness(t *testing.T) {
	levels := []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal}

	for _, level := range levels {
		if name, exists := LevelNameMapping[level]; !exists || name == "" {
			t.Errorf("Level %d missing or empty in LevelNameMapping", level)
		}
	}
}
