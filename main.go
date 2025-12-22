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
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Result represents a palindromic fuel cost finding
type Result struct {
	Litres             float64
	CostPounds         string
	LitresIsPalindrome bool
	Type               string
}

// isPalindrome checks if a number is palindromic
func isPalindrome(n int) bool {
	if n < 0 {
		return false
	}
	if n < 10 {
		return true
	}

	original := n
	reversed := 0

	for n > 0 {
		reversed = reversed*10 + n%10
		n /= 10
	}

	return original == reversed
}

// isPalindromeString checks if a string is palindromic
func isPalindromeString(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

// generatePalindromesForDigits generates all palindromic numbers with a given number of digits
func generatePalindromesForDigits(digits int) []int {
	var palindromes []int

	if digits == 1 {
		for i := 1; i <= 9; i++ {
			palindromes = append(palindromes, i)
		}
		return palindromes
	}

	if digits == 2 {
		for i := 1; i <= 9; i++ {
			palindromes = append(palindromes, i*11)
		}
		return palindromes
	}

	halfDigits := (digits + 1) / 2
	min := int(math.Pow(10, float64(halfDigits-1)))
	max := int(math.Pow(10, float64(halfDigits))) - 1

	for i := min; i <= max; i++ {
		str := strconv.Itoa(i)
		var palindrome string

		if digits%2 == 0 {
			// Even digits: mirror completely
			palindrome = str + reverse(str)
		} else {
			// Odd digits: mirror without center digit
			palindrome = str + reverse(str[:len(str)-1])
		}

		num, _ := strconv.Atoi(palindrome)
		palindromes = append(palindromes, num)
	}

	return palindromes
}

// reverse reverses a string
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// getPalindromicPencesInRange gets all palindromic pence values in a range
func getPalindromicPencesInRange(minPence, maxPence int) []int {
	minDigits := len(strconv.Itoa(minPence))
	maxDigits := len(strconv.Itoa(maxPence))

	var results []int

	for d := minDigits; d <= maxDigits; d++ {
		pals := generatePalindromesForDigits(d)
		for _, pal := range pals {
			if pal >= minPence && pal <= maxPence {
				results = append(results, pal)
			} else if pal > maxPence {
				break
			}
		}
	}

	return results
}

// formatPounds formats pence as pounds string
func formatPounds(pence int) string {
	pounds := float64(pence) / 100.0
	return fmt.Sprintf("%.2f", pounds)
}

// isEffectivelyInteger checks if a float is close enough to an integer
func isEffectivelyInteger(f float64, epsilon float64) bool {
	return math.Abs(f-math.Round(f)) < epsilon
}

// FindPalindromicFuelCosts finds all palindromic fuel costs for a given price
func FindPalindromicFuelCosts(pricePerLitre float64, maxLitres int) []Result {
	var results []Result

	minPence := int(math.Floor(pricePerLitre))
	maxPence := int(math.Ceil(float64(maxLitres) * pricePerLitre))

	// Get all palindromic pence values
	palindromicPences := getPalindromicPencesInRange(minPence, maxPence)

	// Pre-calculate reciprocal for faster division
	reciprocalPrice := 1.0 / pricePerLitre

	for _, pencePrice := range palindromicPences {
		// Check if this palindromic pence is also palindromic as pounds
		poundsStr := formatPounds(pencePrice)
		if !isPalindromeString(poundsStr) {
			continue
		}

		litres := float64(pencePrice) * reciprocalPrice

		// Skip if exceeds max litres or less than 1
		if litres > float64(maxLitres) {
			break
		}
		if litres < 1.0 {
			continue
		}

		// Check if litres is effectively a whole number
		if isEffectivelyInteger(litres, 0.01) {
			wholeLitres := int(math.Round(litres))
			results = append(results, Result{
				Litres:             float64(wholeLitres),
				CostPounds:         poundsStr,
				LitresIsPalindrome: isPalindrome(wholeLitres),
				Type:               "whole",
			})
		} else {
			// Check if litres value itself is palindromic
			litresRounded := math.Round(litres*100) / 100
			litresStr := fmt.Sprintf("%.2f", litresRounded)

			if isPalindromeString(litresStr) {
				results = append(results, Result{
					Litres:             litresRounded,
					CostPounds:         poundsStr,
					LitresIsPalindrome: true,
					Type:               "palindromic_decimal",
				})
			}
		}
	}

	return results
}

// FindNearestPalindromicCost finds the nearest palindromic cost to a target amount
func FindNearestPalindromicCost(pricePerLitre float64, targetLitres float64, searchRadius int) *Result {
	minLitres := int(math.Max(1, targetLitres-float64(searchRadius)))
	maxLitres := int(targetLitres + float64(searchRadius))

	results := FindPalindromicFuelCosts(pricePerLitre, maxLitres)

	var nearest *Result
	minDiff := math.MaxFloat64

	for i := range results {
		if results[i].Litres < float64(minLitres) {
			continue
		}

		diff := math.Abs(results[i].Litres - targetLitres)
		if diff < minDiff {
			minDiff = diff
			nearest = &results[i]
		}
	}

	return nearest
}

// FindPalindromicCostForTarget finds palindromic costs near a target price
func FindPalindromicCostForTarget(pricePerLitre float64, targetPounds float64, searchRadiusPence int) []Result {
	var results []Result

	targetPence := int(math.Round(targetPounds * 100))
	minPence := targetPence - searchRadiusPence
	maxPence := targetPence + searchRadiusPence

	if minPence < 1 {
		minPence = 1
	}

	// Get palindromic pences in range
	palindromicPences := getPalindromicPencesInRange(minPence, maxPence)
	reciprocalPrice := 1.0 / pricePerLitre

	for _, pencePrice := range palindromicPences {
		poundsStr := formatPounds(pencePrice)
		if !isPalindromeString(poundsStr) {
			continue
		}

		litres := float64(pencePrice) * reciprocalPrice

		if litres < 1.0 {
			continue
		}

		if isEffectivelyInteger(litres, 0.01) {
			wholeLitres := int(math.Round(litres))
			results = append(results, Result{
				Litres:             float64(wholeLitres),
				CostPounds:         poundsStr,
				LitresIsPalindrome: isPalindrome(wholeLitres),
				Type:               "whole",
			})
		} else {
			litresRounded := math.Round(litres*100) / 100
			litresStr := fmt.Sprintf("%.2f", litresRounded)

			if isPalindromeString(litresStr) {
				results = append(results, Result{
					Litres:             litresRounded,
					CostPounds:         poundsStr,
					LitresIsPalindrome: true,
					Type:               "palindromic_decimal",
				})
			}
		}
	}

	return results
}

// BatchFindPalindromicCosts processes multiple fuel prices
func BatchFindPalindromicCosts(prices []float64, maxLitres int) map[float64][]Result {
	results := make(map[float64][]Result)

	for _, price := range prices {
		results[price] = FindPalindromicFuelCosts(price, maxLitres)
	}

	return results
}

// Web server types and handlers
type CalculateRequest struct {
	PricePerLitre float64 `json:"pricePerLitre"`
	MaxLitres     int     `json:"maxLitres"`
}

type CalculateResponse struct {
	Results []Result `json:"results"`
	Error   string   `json:"error,omitempty"`
}

type TemplateData struct {
	Results []Result
	Error   string
	Request CalculateRequest
}

// handleAPI handles the REST API endpoint
func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(CalculateResponse{Error: "Invalid JSON"})
			return
		}
	} else {
		// GET request - parse query parameters
		priceStr := r.URL.Query().Get("price")
		maxStr := r.URL.Query().Get("max")

		if priceStr == "" || maxStr == "" {
			json.NewEncoder(w).Encode(CalculateResponse{Error: "Missing price or max parameters"})
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			json.NewEncoder(w).Encode(CalculateResponse{Error: "Invalid price parameter"})
			return
		}

		max, err := strconv.Atoi(maxStr)
		if err != nil {
			json.NewEncoder(w).Encode(CalculateResponse{Error: "Invalid max parameter"})
			return
		}

		req = CalculateRequest{PricePerLitre: price, MaxLitres: max}
	}

	results := FindPalindromicFuelCosts(req.PricePerLitre, req.MaxLitres)
	json.NewEncoder(w).Encode(CalculateResponse{Results: results})
}

// handleWebUI handles the web interface
func handleWebUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "" {
		http.NotFound(w, r)
		return
	}

	data := TemplateData{}

	if r.Method == "POST" {
		r.ParseForm()
		priceStr := r.FormValue("price")
		maxStr := r.FormValue("max")

		if priceStr != "" && maxStr != "" {
			price, err1 := strconv.ParseFloat(priceStr, 64)
			max, err2 := strconv.Atoi(maxStr)

			if err1 == nil && err2 == nil {
				data.Request = CalculateRequest{PricePerLitre: price, MaxLitres: max}
				data.Results = FindPalindromicFuelCosts(price, max)
			} else {
				data.Error = "Invalid input values"
			}
		}
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Palindromic Fuel Calculator</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .container { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        input, button { padding: 8px; margin: 5px; font-size: 16px; }
        button { background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0056b3; }
        .result { background: white; padding: 15px; margin: 10px 0; border-radius: 4px; border-left: 4px solid #28a745; }
        .error { color: #dc3545; background: #f8d7da; padding: 10px; border-radius: 4px; }
        .api-info { background: #e9ecef; padding: 15px; border-radius: 4px; margin-top: 20px; }
        pre { background: #f8f9fa; padding: 10px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>🚗✨ Palindromic Fuel Calculator</h1>
    <p>Find fuel costs that read the same forwards and backwards!</p>

    <div class="container">
        <h2>Calculate</h2>
        <form method="POST">
            <input type="number" name="price" step="0.01" placeholder="Price per litre (e.g., 128.9)" required>
            <input type="number" name="max" placeholder="Max litres (e.g., 100)" required>
            <button type="submit">Calculate</button>
        </form>
    </div>

    {{if .Error}}
    <div class="error">{{.Error}}</div>
    {{end}}

    {{if .Results}}
    <div class="container">
        <h2>Results</h2>
        <p>Fuel Price: {{.Request.PricePerLitre}}p/litre</p>
        <p>Found {{len .Results}} palindromic cost(s):</p>

        {{range .Results}}
        <div class="result">
            <strong>{{if eq .Litres .Litres|printf "%.0f"|atoi|printf "%.0f"}}{{.Litres|printf "%.0f"}}{{else}}{{.Litres|printf "%.2f"}}{{end}} litres = £{{.CostPounds}}</strong>
            <br><small>
                {{if .LitresIsPalindrome}}
                    {{if eq .Type "palindromic_decimal"}}(palindromic decimal litres){{else}}(palindromic whole litres){{end}}
                {{else}}
                    (whole number litres)
                {{end}}
            </small>
        </div>
        {{end}}
    </div>
    {{end}}

    <div class="api-info">
        <h3>API Usage</h3>
        <p>This calculator also provides a REST API:</p>

        <h4>GET Request:</h4>
        <pre>curl "http://localhost:8080/api/calculate?price=128.9&max=100"</pre>

        <h4>POST Request:</h4>
        <pre>curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"pricePerLitre": 128.9, "maxLitres": 100}'</pre>

        <p><a href="/api/calculate?price=128.9&max=50" target="_blank">Try the API</a></p>
    </div>
</body>
</html>`

	t, err := template.New("webui").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

func main() {
	pricePtr := flag.Float64("price", 0, "Price per litre in pence (required)")
	maxLitresPtr := flag.Int("max", 10000, "Maximum litres to check")
	reverseLitresPtr := flag.Float64("reverse-litres", 0, "Find nearest palindrome to this litre amount")
	reversePricePtr := flag.Float64("reverse-price", 0, "Find palindromes near this target price in pounds")
	searchRadiusPtr := flag.Int("radius", 100, "Search radius for reverse lookup")
	batchPtr := flag.String("batch", "", "Comma-separated list of prices for batch processing")
	csvPtr := flag.String("csv", "", "Export results to CSV file (e.g., results.csv)")
	webPtr := flag.Bool("web", false, "Start web server on port 8080")
	portPtr := flag.String("port", "8080", "Port for web server")

	flag.Parse()

	// Web server mode
	if *webPtr {
		fmt.Printf("Starting web server on port %s\n", *portPtr)
		fmt.Printf("Web UI: http://localhost:%s\n", *portPtr)
		fmt.Printf("API: http://localhost:%s/api/calculate\n", *portPtr)

		http.HandleFunc("/", handleWebUI)
		http.HandleFunc("/api/calculate", handleAPI)

		log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
	}

	if *pricePtr == 0 && *batchPtr == "" && !*webPtr {
		fmt.Println("Palindromic Fuel Cost Calculator")
		fmt.Println("================================")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  Normal mode:")
		fmt.Println("    ./palindromic-fuel -price=128.9 -max=100")
		fmt.Println()
		fmt.Println("  With CSV export:")
		fmt.Println("    ./palindromic-fuel -price=128.9 -max=100 -csv=results.csv")
		fmt.Println()
		fmt.Println("  Reverse lookup (find nearest to target litres):")
		fmt.Println("    ./palindromic-fuel -price=128.9 -reverse-litres=50 -radius=100")
		fmt.Println()
		fmt.Println("  Reverse lookup (find palindromes near target price):")
		fmt.Println("    ./palindromic-fuel -price=128.9 -reverse-price=50.00 -radius=500")
		fmt.Println()
		fmt.Println("  Batch mode:")
		fmt.Println("    ./palindromic-fuel -batch=128.9,135.7,142.3 -max=1000")
		fmt.Println()
		fmt.Println("  Batch mode with CSV export:")
		fmt.Println("    ./palindromic-fuel -batch=128.9,135.7,142.3 -max=1000 -csv=batch.csv")
		fmt.Println()
		fmt.Println("  Web server mode:")
		fmt.Println("    ./palindromic-fuel -web")
		fmt.Println("    ./palindromic-fuel -web -port=3000")
		fmt.Println()
		return
	}

	// Batch mode
	if *batchPtr != "" {
		priceStrs := strings.Split(*batchPtr, ",")
		var prices []float64

		for _, priceStr := range priceStrs {
			price, err := strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
			if err != nil {
				fmt.Printf("Error parsing price '%s': %v\n", priceStr, err)
				return
			}
			prices = append(prices, price)
		}

		fmt.Printf("\n=== Batch Processing %d Fuel Prices ===\n", len(prices))
		start := time.Now()
		results := BatchFindPalindromicCosts(prices, *maxLitresPtr)
		elapsed := time.Since(start)

		fmt.Printf("\nTotal batch time: %.3fms\n", float64(elapsed.Microseconds())/1000.0)
		fmt.Printf("Average per price: %.3fms\n\n", float64(elapsed.Microseconds())/1000.0/float64(len(prices)))

		for _, price := range prices {
			printResults(results[price], price)
		}

		// Export to CSV if requested
		if *csvPtr != "" {
			if err := exportBatchToCSV(*csvPtr, results, prices); err != nil {
				fmt.Printf("\nError exporting to CSV: %v\n", err)
			} else {
				fmt.Printf("\nResults exported to %s\n", *csvPtr)
			}
		}

		return
	}

	// Reverse lookup by litres
	if *reverseLitresPtr > 0 {
		fmt.Printf("\nFinding nearest palindromic cost to %.2f litres at %.1fp/litre\n", *reverseLitresPtr, *pricePtr)
		fmt.Printf("Search radius: ±%d litres\n", *searchRadiusPtr)

		start := time.Now()
		result := FindNearestPalindromicCost(*pricePtr, *reverseLitresPtr, *searchRadiusPtr)
		elapsed := time.Since(start)

		if result != nil {
			fmt.Printf("\nNearest palindromic cost:\n")
			printResult(*result)
			diff := math.Abs(result.Litres - *reverseLitresPtr)
			fmt.Printf("Difference: %.2f litres\n", diff)
		} else {
			fmt.Println("\nNo palindromic costs found in search radius")
		}

		fmt.Printf("\nSearch completed in %.3fms\n", float64(elapsed.Microseconds())/1000.0)
		return
	}

	// Reverse lookup by price
	if *reversePricePtr > 0 {
		fmt.Printf("\nFinding palindromic costs near £%.2f at %.1fp/litre\n", *reversePricePtr, *pricePtr)
		fmt.Printf("Search radius: ±%dp\n", *searchRadiusPtr)

		start := time.Now()
		results := FindPalindromicCostForTarget(*pricePtr, *reversePricePtr, *searchRadiusPtr)
		elapsed := time.Since(start)

		if len(results) > 0 {
			fmt.Printf("\nFound %d palindromic cost(s):\n\n", len(results))
			for _, result := range results {
				printResult(result)
				targetDiff := math.Abs(parseFloat(result.CostPounds) - *reversePricePtr)
				fmt.Printf("  Price difference: £%.2f\n", targetDiff)
			}
		} else {
			fmt.Println("\nNo palindromic costs found in search radius")
		}

		fmt.Printf("\nSearch completed in %.3fms\n", float64(elapsed.Microseconds())/1000.0)
		return
	}

	// Normal mode
	start := time.Now()
	results := FindPalindromicFuelCosts(*pricePtr, *maxLitresPtr)
	elapsed := time.Since(start)

	fmt.Printf("\nPerformance: Found %d results in %.3fms\n", len(results), float64(elapsed.Microseconds())/1000.0)
	fmt.Printf("Effective range checked: 1-%d litres\n", *maxLitresPtr)

	printResults(results, *pricePtr)

	// Export to CSV if requested
	if *csvPtr != "" {
		if err := exportToCSV(*csvPtr, results, *pricePtr); err != nil {
			fmt.Printf("\nError exporting to CSV: %v\n", err)
		} else {
			fmt.Printf("\nResults exported to %s\n", *csvPtr)
		}
	}
}

func printResults(results []Result, price float64) {
	fmt.Printf("\nFuel Price: %.1fp/litre\n", price)
	fmt.Printf("Found %d palindromic costs:\n\n", len(results))

	maxShow := 50
	toShow := results
	if len(results) > maxShow {
		toShow = results[:maxShow]
	}

	for _, result := range toShow {
		printResult(result)
	}

	if len(results) > maxShow {
		fmt.Printf("\n... and %d more results\n", len(results)-maxShow)
	}
}

func printResult(result Result) {
	litresStatus := "(whole number litres)"
	if result.LitresIsPalindrome {
		if result.Type == "palindromic_decimal" {
			litresStatus = "(palindromic decimal litres)"
		} else {
			litresStatus = "(palindromic whole litres)"
		}
	}

	if result.Litres == math.Floor(result.Litres) {
		fmt.Printf("%.0f litres = £%s %s\n", result.Litres, result.CostPounds, litresStatus)
	} else {
		fmt.Printf("%.2f litres = £%s %s\n", result.Litres, result.CostPounds, litresStatus)
	}
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// exportToCSV exports results to a CSV file
func exportToCSV(filename string, results []Result, price float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Price per Litre (p)", "Litres", "Cost (£)", "Litres is Palindrome", "Type"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, result := range results {
		litresStr := fmt.Sprintf("%.2f", result.Litres)
		if result.Litres == math.Floor(result.Litres) {
			litresStr = fmt.Sprintf("%.0f", result.Litres)
		}

		litresPalindrome := "No"
		if result.LitresIsPalindrome {
			litresPalindrome = "Yes"
		}

		row := []string{
			fmt.Sprintf("%.1f", price),
			litresStr,
			result.CostPounds,
			litresPalindrome,
			result.Type,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// exportBatchToCSV exports batch results to a CSV file
func exportBatchToCSV(filename string, batchResults map[float64][]Result, prices []float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Price per Litre (p)", "Litres", "Cost (£)", "Litres is Palindrome", "Type"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data for each price
	for _, price := range prices {
		results := batchResults[price]
		for _, result := range results {
			litresStr := fmt.Sprintf("%.2f", result.Litres)
			if result.Litres == math.Floor(result.Litres) {
				litresStr = fmt.Sprintf("%.0f", result.Litres)
			}

			litresPalindrome := "No"
			if result.LitresIsPalindrome {
				litresPalindrome = "Yes"
			}

			row := []string{
				fmt.Sprintf("%.1f", price),
				litresStr,
				result.CostPounds,
				litresPalindrome,
				result.Type,
			}

			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
	}

	return nil
}
