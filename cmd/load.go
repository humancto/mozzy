package cmd

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/humancto/mozzy/internal/httpclient"
)

var (
	loadRequests   int
	loadConcurrent int
	loadDuration   string
)

var loadCmd = &cobra.Command{
	Use:   "load <url>",
	Short: "Run load tests against an endpoint",
	Long: `Run load tests to measure API performance under load.

Examples:
  mozzy load https://api.example.com/users --requests 1000 --concurrent 10
  mozzy load https://api.example.com/users --duration 30s --concurrent 5
  mozzy load https://api.example.com/users --requests 100 --auth "Bearer token"`,
	Args: cobra.ExactArgs(1),
	RunE: runLoad,
}

func init() {
	loadCmd.Flags().IntVar(&loadRequests, "requests", 100, "Total number of requests to send")
	loadCmd.Flags().IntVar(&loadConcurrent, "concurrent", 10, "Number of concurrent workers")
	loadCmd.Flags().StringVar(&loadDuration, "duration", "", "Run for duration (e.g. 30s) instead of fixed requests")
	rootCmd.AddCommand(loadCmd)
}

func runLoad(cmd *cobra.Command, args []string) error {
	url := args[0]

	// Prepare request
	req := httpclient.Request{
		Method:  "GET",
		URL:     url,
		Headers: headers,
		Token:   authToken,
	}

	fmt.Printf("🔥 Load Testing\n")
	fmt.Printf("   Target: %s\n", url)
	fmt.Printf("   Concurrent workers: %d\n", loadConcurrent)

	var totalRequests int64
	var successCount int64
	var errorCount int64
	var totalDuration time.Duration
	var minDuration time.Duration = time.Hour
	var maxDuration time.Duration

	var wg sync.WaitGroup
	var mu sync.Mutex

	startTime := time.Now()
	ctx := context.Background()

	// Determine mode: fixed requests or duration-based
	if loadDuration != "" {
		duration, err := time.ParseDuration(loadDuration)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		fmt.Printf("   Duration: %s\n\n", loadDuration)

		// Duration-based load test
		ctx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()

		for i := 0; i < loadConcurrent; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						res, _, reqDuration, err := httpclient.Do(ctx, req)
						if err != nil {
							atomic.AddInt64(&errorCount, 1)
						} else {
							res.Body.Close()
							if res.StatusCode >= 200 && res.StatusCode < 300 {
								atomic.AddInt64(&successCount, 1)
							} else {
								atomic.AddInt64(&errorCount, 1)
							}

							mu.Lock()
							totalDuration += reqDuration
							if reqDuration < minDuration {
								minDuration = reqDuration
							}
							if reqDuration > maxDuration {
								maxDuration = reqDuration
							}
							mu.Unlock()
						}
						atomic.AddInt64(&totalRequests, 1)
					}
				}
			}()
		}
	} else {
		// Fixed requests load test
		fmt.Printf("   Total requests: %d\n\n", loadRequests)

		requestChan := make(chan struct{}, loadRequests)
		for i := 0; i < loadRequests; i++ {
			requestChan <- struct{}{}
		}
		close(requestChan)

		for i := 0; i < loadConcurrent; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range requestChan {
					res, _, reqDuration, err := httpclient.Do(ctx, req)
					if err != nil {
						atomic.AddInt64(&errorCount, 1)
					} else {
						res.Body.Close()
						if res.StatusCode >= 200 && res.StatusCode < 300 {
							atomic.AddInt64(&successCount, 1)
						} else {
							atomic.AddInt64(&errorCount, 1)
						}

						mu.Lock()
						totalDuration += reqDuration
						if reqDuration < minDuration {
							minDuration = reqDuration
						}
						if reqDuration > maxDuration {
							maxDuration = reqDuration
						}
						mu.Unlock()
					}
					atomic.AddInt64(&totalRequests, 1)
				}
			}()
		}
	}

	// Show progress
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				current := atomic.LoadInt64(&totalRequests)
				fmt.Printf("\r   Progress: %d requests completed...", current)
			}
		}
	}()

	wg.Wait()
	elapsed := time.Since(startTime)

	// Print results
	fmt.Printf("\r                                              \r")
	fmt.Println("\n" + "============================================================")
	fmt.Printf("📊 Load Test Results\n\n")

	fmt.Printf("Requests:\n")
	fmt.Printf("  Total:        %d\n", totalRequests)
	color.Green("  Successful:   %d\n", successCount)
	if errorCount > 0 {
		color.Red("  Failed:       %d\n", errorCount)
	} else {
		fmt.Printf("  Failed:       %d\n", errorCount)
	}

	fmt.Printf("\nTiming:\n")
	fmt.Printf("  Total time:   %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  Requests/sec: %.2f\n", float64(totalRequests)/elapsed.Seconds())

	if totalRequests > 0 {
		avgDuration := time.Duration(int64(totalDuration) / totalRequests)
		fmt.Printf("\nResponse Times:\n")
		fmt.Printf("  Min:          %s\n", minDuration.Round(time.Millisecond))
		fmt.Printf("  Max:          %s\n", maxDuration.Round(time.Millisecond))
		fmt.Printf("  Average:      %s\n", avgDuration.Round(time.Millisecond))
	}

	fmt.Println("============================================================")

	return nil
}
