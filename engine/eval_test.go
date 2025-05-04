package engine

import (
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "simple addition",
			input:    "10 + 20",
			expected: int(30),
			wantErr:  false,
		},
		{
			name:     "function definition and call",
			input:    "fn add(a, b) { return a + b } add(10, 20)",
			expected: int(30),
			wantErr:  false,
		},
		{
			name:     "variable assignment",
			input:    "let x = 10\nlet y = 20\nx + y",
			expected: int(30),
			wantErr:  false,
		},
		{
			name:     "syntax error",
			input:    "10 + * 20",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Eval(tt.input, nil)
			
			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// If we expect an error, we don't need to check the result
			if tt.wantErr {
				return
			}
			
			// Check result value
			if result == nil {
				t.Errorf("Eval() result is nil, expected %v", tt.expected)
				return
			}
			
			// Check the type and value based on expected type
			switch expected := tt.expected.(type) {
			case int:
				// Try to convert the result to int
				var val int
				switch v := result.Value.(type) {
				case int:
					val = v
				case int64:
					val = int(v)
				case float64:
					val = int(v)
				default:
					t.Errorf("Eval() result type = %T, expected int", result.Value)
					return
				}
				
				if val != expected {
					t.Errorf("Eval() result = %v, expected %v", val, expected)
				}
			case string:
				if val, ok := result.Value.(string); !ok || val != expected {
					t.Errorf("Eval() result = %v, expected %v", result.Value, expected)
				}
			case bool:
				if val, ok := result.Value.(bool); !ok || val != expected {
					t.Errorf("Eval() result = %v, expected %v", result.Value, expected)
				}
			default:
				t.Errorf("Unsupported expected type: %T", expected)
			}
		})
	}
}