// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// For more information, please refer to <https://unlicense.org>
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"zero", 0, true},
		{"single digit", 5, true},
		{"negative number", -121, false},
		{"palindrome 121", 121, true},
		{"palindrome 1221", 1221, true},
		{"palindrome 12321", 12321, true},
		{"not palindrome 123", 123, false},
		{"not palindrome 1234", 1234, false},
		{"large palindrome", 123454321, true},
		{"large not palindrome", 123456789, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("isPalindrome(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPalindromeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", true},
		{"single char", "a", true},
		{"palindrome aba", "aba", true},
		{"palindrome abba", "abba", true},
		{"palindrome 12321", "12321", true},
		{"palindrome with spaces", "a b a", true}, // compares characters positionally
		{"not palindrome abc", "abc", false},
		{"not palindrome abcd", "abcd", false},
		{"palindrome with decimal", "50.05", true},
		{"not palindrome 32.14", "32.14", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPalindromeString(tt.input)
			if result != tt.expected {
				t.Errorf("isPalindromeString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"hello", "hello", "olleh"},
		{"palindrome", "aba", "aba"},
		{"numbers", "12345", "54321"},
		{"unicode", "🚗✨", "✨🚗"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reverse(tt.input)
			if result != tt.expected {
				t.Errorf("reverse(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGeneratePalindromesForDigits(t *testing.T) {
	tests := []struct {
		name        string
		digits      int
		expectedLen int
		firstFew    []int
	}{
		{"1 digit", 1, 9, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"2 digits", 2, 9, []int{11, 22, 33, 44, 55, 66, 77, 88, 99}},
		{"3 digits", 3, 90, []int{101, 111, 121, 131, 141, 151, 161, 171, 181, 191}},
		{"4 digits", 4, 90, []int{1001, 1111, 1221, 1331, 1441, 1551, 1661, 1771, 1881, 1991}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generatePalindromesForDigits(tt.digits)
			if len(result) != tt.expectedLen {
				t.Errorf("generatePalindromesForDigits(%d) length = %d, want %d", tt.digits, len(result), tt.expectedLen)
			}
			if len(result) >= len(tt.firstFew) {
				for i, expected := range tt.firstFew {
					if result[i] != expected {
						t.Errorf("generatePalindromesForDigits(%d)[%d] = %d, want %d", tt.digits, i, result[i], expected)
					}
				}
			}
		})
	}
}

func TestFormatPounds(t *testing.T) {
	tests := []struct {
		name     string
		pence    int
		expected string
	}{
		{"zero", 0, "0.00"},
		{"single pence", 1, "0.01"},
		{"ten pence", 10, "0.10"},
		{"one pound", 100, "1.00"},
		{"one pound one pence", 101, "1.01"},
		{"large amount", 12345, "123.45"},
		{"palindrome pence", 3223, "32.23"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPounds(tt.pence)
			if result != tt.expected {
				t.Errorf("formatPounds(%d) = %q, want %q", tt.pence, result, tt.expected)
			}
		})
	}
}

func TestIsEffectivelyInteger(t *testing.T) {
	tests := []struct {
		name     string
		f        float64
		epsilon  float64
		expected bool
	}{
		{"exact integer", 5.0, 0.01, true},
		{"close to integer", 5.001, 0.01, true},
		{"not close enough", 5.02, 0.01, false},
		{"negative close", -3.001, 0.01, true},
		{"decimal", 3.14, 0.01, false},
		{"large epsilon", 5.1, 0.2, true},
		{"zero", 0.0, 0.01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEffectivelyInteger(tt.f, tt.epsilon)
			if result != tt.expected {
				t.Errorf("isEffectivelyInteger(%f, %f) = %v, want %v", tt.f, tt.epsilon, result, tt.expected)
			}
		})
	}
}

func TestGetPalindromicPencesInRange(t *testing.T) {
	tests := []struct {
		name     string
		minPence int
		maxPence int
		expected []int
		checkLen bool // if true, just check length matches (for large ranges)
	}{
		{"small range with palindromes", 10, 50, []int{11, 22, 33, 44}, false},
		{"range with 1-3 digits", 1, 999, nil, true}, // just check we get results
		{"empty range", 50, 10, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPalindromicPencesInRange(tt.minPence, tt.maxPence)
			if tt.checkLen {
				if len(result) == 0 {
					t.Errorf("getPalindromicPencesInRange(%d, %d) returned empty slice, expected some results", tt.minPence, tt.maxPence)
				}
			} else {
				// For empty ranges, just check that result is empty
				if tt.minPence > tt.maxPence {
					if len(result) != 0 {
						t.Errorf("getPalindromicPencesInRange(%d, %d) returned %d results, expected 0 for empty range", tt.minPence, tt.maxPence, len(result))
					}
				} else if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("getPalindromicPencesInRange(%d, %d) = %v, want %v", tt.minPence, tt.maxPence, result, tt.expected)
				}
			}
		})
	}
}

func TestFindPalindromicFuelCosts(t *testing.T) {
	tests := []struct {
		name           string
		pricePerVolume float64
		maxVolume      int
		expectedCount  int // check count since full results would be long
		checkSpecific  bool
		specificResult *Result
	}{
		{"standard price", 128.9, 100, 4, true, &Result{
			Volume:             25.0,
			CostPounds:         "32.23",
			VolumeIsPalindrome: false, // 25 is not a palindrome
			Type:               "whole",
		}},
		{"zero price", 0, 10, 0, false, nil},
		{"very high price", 1000, 1, 0, false, nil},                       // no results in range
		{"fractional price", 0.5, 10, 0, false, nil},                      // results would be < 1 litre
		{"very low max litres", 128.9, 1, 0, false, nil},                  // no results fit in 1 litre
		{"high price with small range", 500.0, 5, 2, false, nil},          // some results still possible
		{"very small price", 0.01, 10000, 0, false, nil},                  // results would be huge litres
		{"price causing fractional litres", 200.0, 10, 1, false, nil},     // should skip some fractional results
		{"price producing sub-litre results", 1000.0, 50, 18, false, nil}, // very high price, small volumes
		{"price for decimal palindrome test", 131.0, 50, 2, false, nil},   // should trigger decimal palindrome check
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindPalindromicFuelCosts(tt.pricePerVolume, tt.maxVolume, 0.01)
			if len(results) != tt.expectedCount {
				t.Errorf("FindPalindromicFuelCosts(%f, %d) returned %d results, want %d",
					tt.pricePerVolume, tt.maxVolume, len(results), tt.expectedCount)
			}

			if tt.checkSpecific && tt.specificResult != nil {
				found := false
				for _, result := range results {
					if result.Volume == tt.specificResult.Volume &&
						result.CostPounds == tt.specificResult.CostPounds &&
						result.VolumeIsPalindrome == tt.specificResult.VolumeIsPalindrome &&
						result.Type == tt.specificResult.Type {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("FindPalindromicFuelCosts(%f, %d) did not contain expected result: %+v",
						tt.pricePerVolume, tt.maxVolume, tt.specificResult)
				}
			}
		})
	}
}

func TestFindNearestPalindromicCost(t *testing.T) {
	tests := []struct {
		name          string
		pricePerVolume float64
		targetVolume  float64
		searchRadius  int
		expectResult  bool
	}{
		{"find near 25", 128.9, 25.0, 10, true},
		{"find near 30", 128.9, 30.0, 10, true},
		{"no result in radius", 128.9, 1000.0, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindNearestPalindromicCost(tt.pricePerVolume, tt.targetVolume, tt.searchRadius, 0.01)
			if tt.expectResult && result == nil {
				t.Errorf("FindNearestPalindromicCost(%f, %f, %d) expected result but got nil",
					tt.pricePerVolume, tt.targetVolume, tt.searchRadius)
			} else if !tt.expectResult && result != nil {
				t.Errorf("FindNearestPalindromicCost(%f, %f, %d) expected nil but got result",
					tt.pricePerVolume, tt.targetVolume, tt.searchRadius)
			}
		})
	}
}

func TestFindPalindromicCostForTarget(t *testing.T) {
	tests := []struct {
		name              string
		pricePerVolume    float64
		targetPounds      float64
		searchRadiusPence int
		expectedCount     int
	}{
		{"find near £32.23", 128.9, 32.23, 100, 1},
		{"find near £50.05", 128.9, 50.05, 100, 1},
		{"no results in radius", 128.9, 1000.00, 10, 0},
		{"multiple results", 128.9, 50.00, 500, 2}, // Should find both 32.23 and 50.05
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindPalindromicCostForTarget(tt.pricePerVolume, tt.targetPounds, tt.searchRadiusPence, 0.01)
			if len(results) != tt.expectedCount {
				t.Errorf("FindPalindromicCostForTarget(%f, %f, %d) returned %d results, want %d",
					tt.pricePerVolume, tt.targetPounds, tt.searchRadiusPence, len(results), tt.expectedCount)
			}
		})
	}
}

func TestBatchFindPalindromicCosts(t *testing.T) {
	tests := []struct {
		name         string
		prices       []float64
		maxVolume    int
		expectedKeys int
	}{
		{"single price", []float64{128.9}, 100, 1},
		{"multiple prices", []float64{128.9, 135.7}, 50, 2},
		{"empty prices", []float64{}, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := BatchFindPalindromicCosts(tt.prices, tt.maxVolume, 0.01)
			if len(results) != tt.expectedKeys {
				t.Errorf("BatchFindPalindromicCosts(%v, %d) returned %d results, want %d keys",
					tt.prices, tt.maxVolume, len(results), tt.expectedKeys)
			}

			// Verify each price has results
			for _, price := range tt.prices {
				if _, exists := results[price]; !exists {
					t.Errorf("BatchFindPalindromicCosts missing results for price %f", price)
				}
			}
		})
	}
}

func TestPrintResult(t *testing.T) {
	tests := []struct {
		name   string
		result Result
	}{
		{"whole number litres", Result{
			Volume:             25.0,
			CostPounds:         "32.23",
			VolumeIsPalindrome: false,
			Type:               "whole",
		}},
		{"palindromic whole litres", Result{
			Volume:             121.0,
			CostPounds:         "50.05",
			VolumeIsPalindrome: true,
			Type:               "whole",
		}},
		{"palindromic decimal litres", Result{
			Volume:             38.83,
			CostPounds:         "50.05",
			VolumeIsPalindrome: true,
			Type:               "palindromic_decimal",
		}},
		{"decimal litres", Result{
			Volume:             15.75,
			CostPounds:         "20.31",
			VolumeIsPalindrome: false,
			Type:               "whole",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that it doesn't panic
			printResult(tt.result)
		})
	}
}

func TestPrintResults(t *testing.T) {
	results := []Result{
		{Volume: 25.0, CostPounds: "32.23", VolumeIsPalindrome: false, Type: "whole"},
		{Volume: 38.83, CostPounds: "50.05", VolumeIsPalindrome: true, Type: "palindromic_decimal"},
	}

	// Test that it doesn't panic
	printResults(results, 128.9)
}

func TestExportToCSV(t *testing.T) {
	results := []Result{
		{Volume: 25.0, CostPounds: "32.23", VolumeIsPalindrome: false, Type: "whole"},
		{Volume: 38.83, CostPounds: "50.05", VolumeIsPalindrome: true, Type: "palindromic_decimal"},
	}

	// Test export to temporary file
	tmpfile, err := os.CreateTemp("", "test_export_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	err = exportToCSV(tmpfile.Name(), results, 128.9)
	if err != nil {
		t.Errorf("exportToCSV failed: %v", err)
	}

	// Read back the file and verify content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("Failed to read exported file: %v", err)
	}

	contentStr := string(content)
	// Check for expected content
	expectedLines := []string{
		"Price per Litre (p),Litres,Cost (£),Litres is Palindrome,Type",
		"128.9,25,32.23,No,whole",
		"128.9,38.83,50.05,Yes,palindromic_decimal",
	}

	for _, line := range expectedLines {
		if !strings.Contains(contentStr, line) {
			t.Errorf("Expected line %q not found in exported CSV", line)
		}
	}
}

func TestHandleAPI(t *testing.T) {
	// Test GET request
	req, err := http.NewRequest("GET", "/api/calculate?price=128.9&max=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response CalculateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response.Results) == 0 {
		t.Errorf("Expected some results, got 0")
	}

	if response.Error != "" {
		t.Errorf("Unexpected error in response: %s", response.Error)
	}
}

func TestHandleAPI_POST(t *testing.T) {
	// Test POST request
	reqBody := CalculateRequest{PricePerVolume: 128.9, MaxVolume: 50}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/api/calculate", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response CalculateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response.Results) == 0 {
		t.Errorf("Expected some results, got 0")
	}
}

func TestHandleAPI_InvalidInput(t *testing.T) {
	// Test invalid GET request
	req, err := http.NewRequest("GET", "/api/calculate?price=invalid&max=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPI)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response CalculateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Error == "" {
		t.Errorf("Expected error in response for invalid input")
	}
}

func TestExportBatchToCSV(t *testing.T) {
	batchResults := map[float64][]Result{
		128.9: {
			{Volume: 25.0, CostPounds: "32.23", VolumeIsPalindrome: false, Type: "whole"},
		},
		135.7: {
			{Volume: 20.0, CostPounds: "27.14", VolumeIsPalindrome: false, Type: "whole"},
		},
	}
	prices := []float64{128.9, 135.7}

	// Test export to temporary file
	tmpfile, err := os.CreateTemp("", "test_batch_export_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	err = exportBatchToCSV(tmpfile.Name(), batchResults, prices)
	if err != nil {
		t.Errorf("exportBatchToCSV failed: %v", err)
	}

	// Read back the file and verify content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("Failed to read exported file: %v", err)
	}

	contentStr := string(content)
	// Check for header
	if !strings.Contains(contentStr, "Price per Litre (p),Litres,Cost (£),Litres is Palindrome,Type") {
		t.Errorf("CSV header not found in exported batch CSV")
	}

	// Check for data from both prices
	if !strings.Contains(contentStr, "128.9,25,32.23,No,whole") {
		t.Errorf("First price data not found in exported batch CSV")
	}
	if !strings.Contains(contentStr, "135.7,20,27.14,No,whole") {
		t.Errorf("Second price data not found in exported batch CSV")
	}
}

func TestExportToCSVEdgeCases(t *testing.T) {
	// Test empty results
	emptyResults := []Result{}
	tmpfile, err := os.CreateTemp("", "test_empty_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	err = exportToCSV(tmpfile.Name(), emptyResults, 128.9)
	if err != nil {
		t.Errorf("exportToCSV with empty results failed: %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("Failed to read exported file: %v", err)
	}

	contentStr := string(content)
	// Should still have header
	if !strings.Contains(contentStr, "Price per Litre (p)") {
		t.Errorf("CSV header not found in empty export")
	}

	// Test with various result types
	diverseResults := []Result{
		{Volume: 25.0, CostPounds: "32.23", VolumeIsPalindrome: true, Type: "whole"},                 // integer litres, palindromic
		{Volume: 38.83, CostPounds: "50.05", VolumeIsPalindrome: false, Type: "palindromic_decimal"}, // decimal, not palindromic litres
		{Volume: 1.5, CostPounds: "1.93", VolumeIsPalindrome: false, Type: "whole"},                  // decimal litres
	}

	tmpfile2, err := os.CreateTemp("", "test_diverse_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile2.Name())
	defer tmpfile2.Close()

	err = exportToCSV(tmpfile2.Name(), diverseResults, 128.9)
	if err != nil {
		t.Errorf("exportToCSV with diverse results failed: %v", err)
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"valid float", "32.23", 32.23},
		{"integer string", "50", 50.0},
		{"negative", "-10.5", -10.5},
		{"zero", "0.00", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("parseFloat(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmarks

func BenchmarkFindPalindromicFuelCosts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPalindromicFuelCosts(128.9, 1000, 0.01)
	}
}

func BenchmarkGeneratePalindromesForDigits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generatePalindromesForDigits(5)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isPalindrome(12321)
	}
}

func BenchmarkIsPalindromeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isPalindromeString("12321")
	}
}
