package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/machinebox/graphql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	graphQLEndpoint    string
	port               int
	processingInterval time.Duration
	lastProcessedBlock int
)

var (
	totalTransactions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gno_transactions_total",
			Help: "Total number of transactions processed",
		},
	)

	senderActivity = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gno_sender_activity_total",
			Help: "Number of transactions by sender address",
		},
		[]string{"address"},
	)

	packageActivity = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gno_package_activity_total",
			Help: "Number of calls to packages/realms",
		},
		[]string{"package_path"},
	)

	processedBlocks = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gno_blocks_processed_total",
			Help: "Total number of blocks processed",
		},
	)

	latestProcessedHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gno_latest_processed_height",
			Help: "Tracks the latest processed block height",
		},
	)

	successfulTransactions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gno_transactions_success_total",
			Help: "Total number of successful transactions",
		},
	)

	failedTransactions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gno_transactions_failed_total",
			Help: "Total number of failed transactions",
		},
	)
)

func registerMetrics() {
	prometheus.MustRegister(
		totalTransactions,
		senderActivity,
		packageActivity,
		processedBlocks,
		successfulTransactions,
		failedTransactions,
		latestProcessedHeight,
	)
}

// Main function
func main() {
	registerMetrics()
	graphQLEndpoint = getEnv("INDEXER_URL", "http://localhost:8546/graphql/query")
	portStr := getEnv("METRICS_PORT", "8080")
	intervalStr := getEnv("PROCESSING_INTERVAL", "5s")

	var err error
	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT: %v", err)
	}

	processingInterval, err = time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("Invalid PROCESSING_INTERVAL: %v", err)
	}

	client := graphql.NewClient(graphQLEndpoint)
	go collectMetrics(client)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting metrics server on", ":"+strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getLatestBlockHeight(client *graphql.Client) (int, error) {
	req := graphql.NewRequest(`
        query {
            latestBlockHeight
        }
    `)

	var latestBlockResp struct {
		LatestBlockHeight int `json:"latestBlockHeight"`
	}

	if err := client.Run(context.Background(), req, &latestBlockResp); err != nil {
		return 0, err
	}

	return latestBlockResp.LatestBlockHeight, nil
}

func collectMetrics(client *graphql.Client) {
	latestHeight, err := getLatestBlockHeight(client)
	if err != nil {
		log.Fatalf("Error querying latest block height: %v", err)
	}

	log.Printf("Latest block height: %d", latestHeight)

	processHistoricalBlocks(client, latestHeight)
	subscribeToNewBlocks(client)
}

func processHistoricalBlocks(client *graphql.Client, latestHeight int) {
	log.Println("Processing historical blocks...")

	for height := 1; height <= latestHeight; height++ {
		processBlockTransactions(client, height)
		processedBlocks.Inc()
		lastProcessedBlock = height
		latestProcessedHeight.Set(float64(lastProcessedBlock))
		log.Printf("Processed block: %d", height)
	}

	log.Println("Finished processing historical blocks")
}

func processBlockTransactions(client *graphql.Client, blockHeight int) {
	req := graphql.NewRequest(`
    query GetTransactions($fromIndex: Int!, $toIndex: Int!) {
        transactions(filter: { 
            from_index: $fromIndex, 
            to_index: $toIndex
        }) {
            hash
            block_height
            success
            messages {
                typeUrl
                route
                value {
                    ... on BankMsgSend {
						from_address
                    	to_address
                    }
                    ... on MsgAddPackage {
                        creator
                        package {
                        	path
                        	name
                    	}
                    }
                }
            }
        }
    }
`)

	req.Var("fromIndex", 0)
	req.Var("toIndex", 99999)

	var resp struct {
		Transactions []Transaction `json:"transactions"`
	}

	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		log.Printf("Error querying transactions: %v", err)
		return
	}

	var blockTxs []Transaction
	for _, tx := range resp.Transactions {
		if tx.BlockHeight == blockHeight {
			blockTxs = append(blockTxs, tx)
		}
	}

	log.Printf("Found %d transactions in block %d", len(blockTxs), blockHeight)

	for _, tx := range blockTxs {
		processTx(tx)
	}
}

func processTx(tx Transaction) {
	totalTransactions.Inc()

	log.Printf("Processing tx %s (success: %v) with %d messages",
		tx.Hash, tx.Success, len(tx.Messages))

	if tx.Success {
		successfulTransactions.Inc()
	} else {
		failedTransactions.Inc()
	}

	for _, msg := range tx.Messages {
		log.Printf("Processing message type: %s", msg.TypeUrl)
		handleMessage(msg)
	}
}

func handleMessage(msg TransactionMessage) {
	switch msg.TypeUrl {
	case "send":
		handleSend(msg)
	case "add_package":
		handleAddPackage(msg)
	default:
		log.Printf("Unknown message type: %s", msg.TypeUrl)
	}
}

func handleSend(msg TransactionMessage) {
	if msg.Value != nil {
		var msgSend BankMsgSend
		if err := json.Unmarshal(msg.Value, &msgSend); err == nil {
			senderActivity.WithLabelValues(msgSend.FromAddress).Inc()
			log.Printf("Processed bank.MsgSend from %s", msgSend.FromAddress)
		} else {
			log.Printf("Error unmarshaling bank.MsgSend: %v", err)
		}
	}
}

func handleAddPackage(msg TransactionMessage) {
	if msg.Value != nil {
		var msgAddPkg MsgAddPackage
		if err := json.Unmarshal(msg.Value, &msgAddPkg); err == nil {
			senderActivity.WithLabelValues(msgAddPkg.Creator).Inc()
			packageActivity.WithLabelValues(msgAddPkg.Package.Path).Inc()
			log.Printf("Processed m_addpkg from %s to %s", msgAddPkg.Creator, msgAddPkg.Package.Path)
		} else {
			log.Printf("Error unmarshaling m_addpkg: %v", err)
		}
	}
}

func subscribeToNewBlocks(client *graphql.Client) {
	log.Println("Subscribing to new blocks...")

	for {
		latestHeight, err := getLatestBlockHeight(client)
		if err != nil {
			log.Printf("Error fetching latest block height: %v", err)
			time.Sleep(processingInterval)
			continue
		}

		for block := lastProcessedBlock + 1; block <= latestHeight; block++ {
			processBlockTransactions(client, block)
			processedBlocks.Inc()
			lastProcessedBlock = block
			latestProcessedHeight.Set(float64(block))
			log.Printf("Processed new block: %d", block)
		}

		time.Sleep(processingInterval)
	}
}
