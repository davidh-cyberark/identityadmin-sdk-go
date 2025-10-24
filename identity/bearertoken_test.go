package identity

import (
	"testing"
	"time"
)

func TestNewBearerToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		tokenType      string
		expireSeconds  int
		expectedToken  string
		expectedType   string
		validateExpiry bool
	}{
		{
			name:           "valid bearer token with default expiration",
			token:          "abc123xyz",
			tokenType:      "Bearer",
			expireSeconds:  300,
			expectedToken:  "abc123xyz",
			expectedType:   "Bearer",
			validateExpiry: true,
		},
		{
			name:           "valid custom token type",
			token:          "custom-token-456",
			tokenType:      "CustomAuth",
			expireSeconds:  600,
			expectedToken:  "custom-token-456",
			expectedType:   "CustomAuth",
			validateExpiry: true,
		},
		{
			name:           "empty token string",
			token:          "",
			tokenType:      "Bearer",
			expireSeconds:  300,
			expectedToken:  "",
			expectedType:   "Bearer",
			validateExpiry: true,
		},
		{
			name:           "empty token type",
			token:          "test-token",
			tokenType:      "",
			expireSeconds:  300,
			expectedToken:  "test-token",
			expectedType:   "",
			validateExpiry: true,
		},
		{
			name:           "zero expiration seconds",
			token:          "zero-expire-token",
			tokenType:      "Bearer",
			expireSeconds:  0,
			expectedToken:  "zero-expire-token",
			expectedType:   "Bearer",
			validateExpiry: false, // Don't validate expiry for edge case
		},
		{
			name:           "negative expiration seconds",
			token:          "negative-expire-token",
			tokenType:      "Bearer",
			expireSeconds:  -60,
			expectedToken:  "negative-expire-token",
			expectedType:   "Bearer",
			validateExpiry: false, // Don't validate expiry for edge case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTime := time.Now()
			result := NewBearerToken(tt.token, tt.tokenType, tt.expireSeconds)
			afterTime := time.Now()

			// Validate basic structure
			if result == nil {
				t.Fatalf("NewBearerToken() returned nil")
			}

			if result.Token != tt.expectedToken {
				t.Errorf("NewBearerToken() Token = %q, expected %q", result.Token, tt.expectedToken)
			}

			if result.TokenType != tt.expectedType {
				t.Errorf("NewBearerToken() TokenType = %q, expected %q", result.TokenType, tt.expectedType)
			}

			// Validate expiration time
			if result.Expires == nil {
				t.Errorf("NewBearerToken() Expires is nil")
			} else if tt.validateExpiry {
				// Check that the expiration is set to the specified seconds from now
				expectedExpiry := beforeTime.Add(time.Duration(tt.expireSeconds) * time.Second)
				maxExpectedExpiry := afterTime.Add(time.Duration(tt.expireSeconds) * time.Second)

				if result.Expires.Before(expectedExpiry) || result.Expires.After(maxExpectedExpiry) {
					t.Errorf("NewBearerToken() expiration time %v is not within expected range [%v, %v]",
						result.Expires, expectedExpiry, maxExpectedExpiry)
				}
			}
		})
	}
}

func TestBearerToken_SetExpiration(t *testing.T) {
	tests := []struct {
		name    string
		seconds int
		wantErr bool
	}{
		{
			name:    "positive seconds",
			seconds: 300,
			wantErr: false,
		},
		{
			name:    "zero seconds",
			seconds: 0,
			wantErr: false,
		},
		{
			name:    "negative seconds",
			seconds: -60,
			wantErr: false,
		},
		{
			name:    "large positive value",
			seconds: 86400, // 24 hours
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &BearerToken{
				Token:     "test-token",
				TokenType: "Bearer",
			}

			beforeTime := time.Now()
			token.SetExpiration(tt.seconds)
			afterTime := time.Now()

			// Validate that expiration was set
			if token.Expires == nil {
				t.Errorf("SetExpiration() did not set Expires field")
				return
			}

			// Calculate expected expiration time range
			expectedMin := beforeTime.Add(time.Duration(tt.seconds) * time.Second)
			expectedMax := afterTime.Add(time.Duration(tt.seconds) * time.Second)

			// Validate expiration time is within expected range
			if token.Expires.Before(expectedMin) || token.Expires.After(expectedMax) {
				t.Errorf("SetExpiration(%d) set expiration to %v, expected between %v and %v",
					tt.seconds, token.Expires, expectedMin, expectedMax)
			}
		})
	}
}

func TestBearerToken_SetExpiration_NilToken(t *testing.T) {
	// Test that calling SetExpiration on a nil token doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetExpiration() panicked when called on nil token: %v", r)
		}
	}()

	var token *BearerToken
	// This should cause a panic due to nil pointer dereference
	// We expect this test to fail, demonstrating the need for nil checks
	t.Run("nil token should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("SetExpiration() should panic when called on nil token")
			}
		}()
		token.SetExpiration(300)
	})
}

func TestNewBearerToken_Integration_WithSetExpiration(t *testing.T) {
	t.Run("create token and modify expiration", func(t *testing.T) {
		// Create token with constructor
		token := NewBearerToken("integration-token", "Bearer", 300)

		// Verify initial state
		if token.Token != "integration-token" {
			t.Errorf("Token = %q, expected %q", token.Token, "integration-token")
		}
		if token.TokenType != "Bearer" {
			t.Errorf("TokenType = %q, expected %q", token.TokenType, "Bearer")
		}

		// Store original expiration
		originalExpiration := *token.Expires

		// Wait a small amount to ensure time difference
		time.Sleep(1 * time.Millisecond)

		// Modify expiration
		newExpirationSeconds := 600
		token.SetExpiration(newExpirationSeconds)

		// Verify expiration was changed
		if token.Expires.Equal(originalExpiration) {
			t.Errorf("SetExpiration() did not change the expiration time")
		}

		// Verify token is still valid
		if !token.IsValid() {
			t.Errorf("Token should still be valid after SetExpiration()")
		}

		// Test string representation still works
		expected := "Bearer integration-token"
		if token.String() != expected {
			t.Errorf("String() = %q, expected %q", token.String(), expected)
		}
	})
}

func TestNewBearerToken_EdgeCases(t *testing.T) {
	t.Run("extremely large expiration", func(t *testing.T) {
		token := NewBearerToken("test", "Bearer", 31536000) // 1 year
		if !token.IsValid() {
			t.Errorf("Token with 1 year expiration should be valid")
		}
	})

	t.Run("unicode token and type", func(t *testing.T) {
		token := NewBearerToken("🔑token", "🛡️Auth", 300)
		expected := "🛡️Auth 🔑token"
		if token.String() != expected {
			t.Errorf("String() = %q, expected %q", token.String(), expected)
		}
	})

	t.Run("very long token string", func(t *testing.T) {
		longToken := make([]byte, 1000)
		for i := range longToken {
			longToken[i] = 'a'
		}
		token := NewBearerToken(string(longToken), "Bearer", 300)
		if len(token.Token) != 1000 {
			t.Errorf("Long token not preserved correctly")
		}
	})
}

func TestBearerToken_SetExpiration_EdgeCases(t *testing.T) {
	t.Run("multiple calls to SetExpiration", func(t *testing.T) {
		token := &BearerToken{Token: "test", TokenType: "Bearer"}

		// First call
		token.SetExpiration(100)
		firstExpiration := *token.Expires

		// Wait to ensure time difference
		time.Sleep(1 * time.Millisecond)

		// Second call
		token.SetExpiration(200)
		secondExpiration := *token.Expires

		if secondExpiration.Before(firstExpiration) {
			t.Errorf("Second SetExpiration call should result in later expiration")
		}
	})

	t.Run("SetExpiration with max int", func(t *testing.T) {
		token := &BearerToken{Token: "test", TokenType: "Bearer"}

		// This might overflow or cause issues depending on the system
		// We just want to ensure it doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SetExpiration with large value should not panic: %v", r)
			}
		}()

		token.SetExpiration(2147483647) // Max int32
		if token.Expires == nil {
			t.Errorf("SetExpiration should set expiration even with large values")
		}
	})
}

func TestBearerToken_String(t *testing.T) {
	tests := []struct {
		name     string
		token    *BearerToken
		expected string
	}{
		{
			name: "valid token with Bearer type",
			token: &BearerToken{
				Token:     "abc123",
				TokenType: "Bearer",
			},
			expected: "Bearer abc123",
		},
		{
			name: "valid token with custom type",
			token: &BearerToken{
				Token:     "xyz789",
				TokenType: "Custom",
			},
			expected: "Custom xyz789",
		},
		{
			name: "empty token",
			token: &BearerToken{
				Token:     "",
				TokenType: "Bearer",
			},
			expected: "Bearer ",
		},
		{
			name: "empty token type",
			token: &BearerToken{
				Token:     "abc123",
				TokenType: "",
			},
			expected: " abc123",
		},
		{
			name: "both empty",
			token: &BearerToken{
				Token:     "",
				TokenType: "",
			},
			expected: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.String()
			if result != tt.expected {
				t.Errorf("BearerToken.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestBearerToken_IsValid(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(5 * time.Minute)
	pastTime := now.Add(-5 * time.Minute)

	tests := []struct {
		name     string
		token    *BearerToken
		expected bool
	}{
		{
			name:     "nil token",
			token:    nil,
			expected: false,
		},
		{
			name: "empty token string",
			token: &BearerToken{
				Token:     "",
				TokenType: "Bearer",
				Expires:   &futureTime,
			},
			expected: false,
		},
		{
			name: "nil expires",
			token: &BearerToken{
				Token:     "abc123",
				TokenType: "Bearer",
				Expires:   nil,
			},
			expected: false,
		},
		{
			name: "token not expired (future expiration)",
			token: &BearerToken{
				Token:     "abc123",
				TokenType: "Bearer",
				Expires:   &futureTime,
			},
			expected: true, // Fixed: valid token should return true
		},
		{
			name: "token expired (past expiration)",
			token: &BearerToken{
				Token:     "abc123",
				TokenType: "Bearer",
				Expires:   &pastTime,
			},
			expected: false, // Fixed: expired token should return false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.IsValid()
			if result != tt.expected {
				t.Errorf("BearerToken.IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBearerTokenConstants(t *testing.T) {
	t.Run("BearerTokenDefaultExpirationSeconds", func(t *testing.T) {
		expected := 300
		if BearerTokenDefaultExpirationSeconds != expected {
			t.Errorf("BearerTokenDefaultExpirationSeconds = %d, expected %d",
				BearerTokenDefaultExpirationSeconds, expected)
		}
	})
}

func TestBearerToken_IsValid_EdgeCases(t *testing.T) {

	t.Run("token expires exactly now", func(t *testing.T) {
		exactNow := time.Now()
		token := &BearerToken{
			Token:     "abc123",
			TokenType: "Bearer",
			Expires:   &exactNow,
		}

		// This test might be flaky due to timing, but tests the edge case
		result := token.IsValid()
		if result {
			t.Errorf("BearerToken.IsValid() for token expiring exactly now should return false, got %v", result)
		}
	})
}

// TestBearerToken_Integration tests the complete workflow
func TestBearerToken_Integration(t *testing.T) {
	t.Run("create token and validate lifecycle", func(t *testing.T) {
		futureTime := time.Now().Add(5 * time.Minute)

		token := &BearerToken{
			Token:     "integration-test-token",
			TokenType: "Bearer",
			Expires:   &futureTime,
		}

		// Test string representation
		expectedString := "Bearer integration-test-token"
		if token.String() != expectedString {
			t.Errorf("Token string representation = %q, expected %q",
				token.String(), expectedString)
		}

		// Test validity (with corrected implementation logic)
		if !token.IsValid() {
			t.Errorf("Fresh token should be valid with corrected implementation logic, got invalid")
		}
	})
}

// Benchmark tests for performance-critical methods
func BenchmarkBearerToken_String(b *testing.B) {
	token := &BearerToken{
		Token:     "benchmark-token-123",
		TokenType: "Bearer",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = token.String()
	}
}

func BenchmarkBearerToken_IsValid(b *testing.B) {
	futureTime := time.Now().Add(5 * time.Minute)
	token := &BearerToken{
		Token:     "benchmark-token-123",
		TokenType: "Bearer",
		Expires:   &futureTime,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = token.IsValid()
	}
}

func BenchmarkNewBearerToken(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBearerToken("benchmark-token", "Bearer", 300)
	}
}

func BenchmarkBearerToken_SetExpiration(b *testing.B) {
	token := &BearerToken{
		Token:     "benchmark-token-123",
		TokenType: "Bearer",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token.SetExpiration(300)
	}
}
