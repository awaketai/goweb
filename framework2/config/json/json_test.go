package json

import (
	"os"
	"reflect"
	"testing"
)

func createTempFile(content string) (string, error) {
	file, err := os.CreateTemp("", "test.json")
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return "", err
	}
	return file.Name(), nil
}

func TestJSONConfig_Parse(t *testing.T) {
	content := `{
		"name": "test",
		"age": 25,
		"active": true,
		"score": 95.5,
		"tags": "go;json;test"
	}`
	tempFile, err := createTempFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile)

	j := &JSONConfig{}
	configer, err := j.Parse(tempFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val, err := configer.String("name"); err != nil || val != "test" {
		t.Errorf("expected name = test, got %v, error: %v", val, err)
	}

	if val, err := configer.Int("age"); err != nil || val != 25 {
		t.Errorf("expected age = 25, got %v, error: %v", val, err)
	}

	if val, err := configer.Bool("active"); err != nil || val != true {
		t.Errorf("expected active = true, got %v, error: %v", val, err)
	}

	if val, err := configer.Float("score"); err != nil || val != 95.5 {
		t.Errorf("expected score = 95.5, got %v, error: %v", val, err)
	}

	if val, err := configer.Strings("tags"); err != nil || !reflect.DeepEqual(val, []string{"go", "json", "test"}) {
		t.Errorf("expected tags = [go json test], got %v, error: %v", val, err)
	}
}

func TestJSONConfig_Defaults(t *testing.T) {
	content := `{
		"name": "test",
		"age": 25
	}`

	j := &JSONConfig{}
	configer, err := j.ParseData([]byte(content))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val := configer.DefaultString("name", "default"); val != "test" {
		t.Errorf("expected default string to be 'test', got %v", val)
	}

	if val := configer.DefaultString("unknown", "default"); val != "default" {
		t.Errorf("expected default string to be 'default', got %v", val)
	}

	if val := configer.DefaultInt("age", 10); val != 25 {
		t.Errorf("expected default int to be 25, got %v", val)
	}

	if val := configer.DefaultInt("unknown", 10); val != 10 {
		t.Errorf("expected default int to be 10, got %v", val)
	}

	if val := configer.DefaultBool("active", true); val != true {
		t.Errorf("expected default bool to be true, got %v", val)
	}
}

func TestJSONConfig_ParseInvalidData(t *testing.T) {
	invalidContent := `invalid-json`
	j := &JSONConfig{}
	_, err := j.ParseData([]byte(invalidContent))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestJSONConfig_Set(t *testing.T) {
	content := `{
		"name": "test"
	}`

	j := &JSONConfig{}
	configer, err := j.ParseData([]byte(content))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = configer.Set("name", "newName")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val, err := configer.String("name"); err != nil || val != "newName" {
		t.Errorf("expected name = newName, got %v, error: %v", val, err)
	}
}
