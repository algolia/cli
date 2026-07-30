package crawler

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestRenderJavaScript_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    bool
		wantErr bool
	}{
		{name: "bool true", payload: `true`, want: true},
		{name: "bool false", payload: `false`, want: false},
		{
			name:    "array of patterns",
			payload: `["https://example.com/docs/**","https://example.com/api/**"]`,
			want:    true,
		},
		{
			name:    "object with enabled false",
			payload: `{"enabled":false,"waitTime":{"min":1000,"max":5000}}`,
			want:    false,
		},
		{
			name:    "object with enabled true",
			payload: `{"enabled":true,"adblock":true,"patterns":["https://example.com/**"]}`,
			want:    true,
		},
		{
			name:    "object without enabled",
			payload: `{"waitTime":{"min":1000,"max":5000}}`,
			want:    true,
		},
		{name: "number", payload: `42`, wantErr: true},
		{name: "string", payload: `"yes"`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got RenderJavaScript
			err := json.Unmarshal([]byte(tt.payload), &got)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for payload %s, got nil", tt.payload)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for payload %s: %v", tt.payload, err)
			}
			if got.Enabled != tt.want {
				t.Errorf("expected Enabled %t, got %t", tt.want, got.Enabled)
			}
		})
	}
}

func TestRenderJavaScript_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{name: "bool", payload: `{"renderJavaScript":true}`},
		{
			name:    "array",
			payload: `{"renderJavaScript":["https://example.com/docs/**","https://example.com/api/**"]}`,
		},
		{
			name:    "object",
			payload: `{"renderJavaScript":{"enabled":true,"waitTime":{"min":1000,"max":5000},"adblock":true,"patterns":["https://example.com/**"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config Config
			if err := json.Unmarshal([]byte(tt.payload), &config); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			got, err := json.Marshal(&config)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			var want, have interface{}
			if err := json.Unmarshal([]byte(tt.payload), &want); err != nil {
				t.Fatalf("unmarshal expected payload: %v", err)
			}
			if err := json.Unmarshal(got, &have); err != nil {
				t.Fatalf("unmarshal marshalled payload: %v", err)
			}
			if !reflect.DeepEqual(want, have) {
				t.Errorf("expected %s, got %s", tt.payload, got)
			}
		})
	}
}

func TestCrawler_UnmarshalJSON_RenderJavaScriptPatterns(t *testing.T) {
	payload := `{
		"id": "1a2b3c4d-1234-5678-90ab-cdef12345678",
		"name": "my-crawler",
		"running": true,
		"createdAt": "2024-01-02T03:04:05.000Z",
		"config": {
			"appId": "APP_ID",
			"indexPrefix": "crawler_",
			"startUrls": ["https://example.com"],
			"renderJavaScript": ["https://example.com/docs/**"],
			"rateLimit": 8,
			"actions": [
				{
					"indexName": "docs",
					"pathsToMatch": ["https://example.com/docs/**"],
					"recordExtractor": {"__type": "function", "source": "() => {}"}
				}
			]
		}
	}`

	var crawler Crawler
	if err := json.Unmarshal([]byte(payload), &crawler); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if crawler.Config == nil {
		t.Fatal("expected a config, got nil")
	}
	if crawler.Config.RenderJavaScript == nil {
		t.Fatal("expected a renderJavaScript value, got nil")
	}
	if !crawler.Config.RenderJavaScript.Enabled {
		t.Error("expected Enabled true, got false")
	}
}

func TestConfig_MarshalJSON_WithoutRenderJavaScript(t *testing.T) {
	got, err := json.Marshal(&Config{AppID: "APP_ID"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(got), "renderJavaScript") {
		t.Errorf("expected no renderJavaScript key, got %s", got)
	}
}

func TestRenderJavaScript_MarshalJSON_FromGo(t *testing.T) {
	tests := []struct {
		name  string
		value *RenderJavaScript
		want  string
	}{
		{name: "enabled", value: &RenderJavaScript{Enabled: true}, want: `true`},
		{name: "disabled", value: &RenderJavaScript{}, want: `false`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
