package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func stressTest(cmd *cobra.Command, args []string) {
	url := cmd.Flag("url").Value.String()
	if url == "" {
		cmd.Help()
		os.Exit(1)
	}

	requests := cmd.Flag("requests").Value.String()
	concurrency := cmd.Flag("concurrency").Value.String()

	if requests == "" || concurrency == "" {
		cmd.Help()
		os.Exit(1)
	}

	concurrencyInt, err := strconv.Atoi(concurrency)
	if err != nil {
		fmt.Println("Invalid concurrency value")
		os.Exit(1)
	}

	requestsInt, err := strconv.Atoi(requests)
	if err != nil {
		fmt.Println("Invalid requests value")
		os.Exit(1)
	}

	startTime := time.Now()

	successRequests := 0
	successRequestsMutex := sync.Mutex{}
	invalidRequestsMap := make(map[int]int)
	invalidRequestsMapMutex := sync.Mutex{}
	connectionErrors := 0
	connectionErrorsMutex := sync.Mutex{}

	var wg sync.WaitGroup
	for i := 0; i < requestsInt; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println("Error creating request")
				return
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				connectionErrorsMutex.Lock()
				connectionErrors++
				connectionErrorsMutex.Unlock()
				return
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				invalidRequestsMapMutex.Lock()
				invalidRequestsMap[res.StatusCode] = invalidRequestsMap[res.StatusCode] + 1
				invalidRequestsMapMutex.Unlock()

				return
			}

			successRequestsMutex.Lock()
			successRequests++
			successRequestsMutex.Unlock()
		}()

		if i%concurrencyInt == 0 {
			wg.Wait()
		}
	}

	wg.Wait()

	elapsed := time.Since(startTime)

	fmt.Println("Stress test complete")

	fmt.Println("=====================================")
	fmt.Println("Total requests:", requestsInt)

	fmt.Println("=====================================")
	fmt.Printf("Total time: %.2fs \n", elapsed.Seconds())
	fmt.Println("Successful requests:", successRequests)

	fmt.Println("=====================================")
	if len(invalidRequestsMap) > 0 {
		fmt.Println("Invalid requests:")

		for k, v := range invalidRequestsMap {
			fmt.Printf("Status code: %d - Count: %d\n", k, v)
		}
	} else {
		fmt.Println("No invalid requests")
	}

	fmt.Println("=====================================")
	fmt.Println("Connection errors:", connectionErrors)
}

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "Stress test your API",
	Long:  "Stress test your API",
	Run:   stressTest,
}

func init() {
	rootCmd.Flags().StringP("url", "r", "", "URL to stress test")
	rootCmd.Flags().StringP("requests", "n", "", "Number of requests to make")
	rootCmd.Flags().StringP("concurrency", "c", "", "Number of concurrent requests to make")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
