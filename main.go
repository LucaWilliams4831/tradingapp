package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"database/sql"
	_ "github.com/lib/pq"
)

func UnmarshalDexSpyResponse(data []byte) (DexSpyResponse, error) {
	var r DexSpyResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DexSpyResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DexSpyResponse struct {
	Items      []Item `json:"items"`
	PageIndex  int64  `json:"pageIndex"`
	PageSize   int64  `json:"pageSize"`
	TotalItems int64  `json:"totalItems"`
}

type Item struct {
	ID                   string        `json:"_id"`
	Address              string        `json:"address"`
	Block                int64         `json:"block"`
	CreateDate           string        `json:"createDate"`
	Creator              string        `json:"creator"`
	DeployerWalletAge    float64       `json:"deployerWalletAge"`
	HexCodeHash          string        `json:"hexCodeHash"`
	Name                 string        `json:"name"`
	Symbol               string        `json:"symbol"`
	Decimals             int64         `json:"decimals"`
	TotalSupply          string        `json:"totalSupply"`
	StartLiquidity       float64       `json:"startLiquidity"`
	StartTokenReserve    string        `json:"startTokenReserve"`
	CurrentLiquidity     float64       `json:"currentLiquidity"`
	CurrentTokenReserve  string        `json:"currentTokenReserve"`
	StartMcap            float64         `json:"startMcap"`
	CurrentMcap          float64         `json:"currentMcap"`
	CurrentGain          float64       `json:"currentGain"`
	AthMcap              float64         `json:"athMcap"`
	AthGain              float64       `json:"athGain"`
	LiquidityAdded       bool          `json:"liquidityAdded"`
	LiquidityRemoved     bool          `json:"liquidityRemoved"`
	CouldBuy             bool          `json:"couldBuy"`
	CouldSell            bool          `json:"couldSell"`
	DelayedHoneypot      bool          `json:"delayedHoneypot"`
	BuyTax               int64         `json:"buyTax"`
	SellTax              int64         `json:"sellTax"`
	MaxBuy               string        `json:"maxBuy"`
	UpdateBlock          int64         `json:"updateBlock"`
	V                    int64         `json:"__v"`
	LiquidityInUsd       float64       `json:"liquidityInUsd"`
	MaxBuyPercentage     float64       `json:"maxBuyPercentage"`
	MaxBuyTokens         string        `json:"maxBuyTokens"`
	Owner                *string       `json:"owner,omitempty"`
	PairAddress          string        `json:"pairAddress"`
	DeadTokens           string        `json:"deadTokens"`
	DeadTokensPercentage float64         `json:"deadTokensPercentage"`
	LiquidityRatio       float64         `json:"liquidityRatio"`
	PairToken            string     	`json:"pairToken"`
	PairTokenName        string 		`json:"pairTokenName"`
	InitMaxBuy           *string       `json:"initMaxBuy,omitempty"`
	InitMaxBuyPercentage *float64      `json:"initMaxBuyPercentage,omitempty"`
	InitMaxBuyTokens     *string       `json:"initMaxBuyTokens,omitempty"`
	LaunchTime           *string       `json:"launchTime,omitempty"`
	Launched             *bool         `json:"launched,omitempty"`
	LiquidityBurnDesc    *string       `json:"liquidityBurnDesc,omitempty"`
	LiquidityBurnPercent *float64      `json:"liquidityBurnPercent,omitempty"`
	LiquidityBurnTxHash  *string       `json:"liquidityBurnTxHash,omitempty"`
	LiquidityBurned      *bool         `json:"liquidityBurned,omitempty"`
	LiquidityLockDays    *float64      `json:"liquidityLockDays,omitempty"`
	LiquidityLockDesc    *string       `json:"liquidityLockDesc,omitempty"`
	LiquidityLockTxHash  *string       `json:"liquidityLockTxHash,omitempty"`
	LiquidityLocked      *bool         `json:"liquidityLocked,omitempty"`
	LiquidityUnlockDate  *string       `json:"liquidityUnlockDate,omitempty"`
	Renounced            *bool         `json:"renounced,omitempty"`
}
func main() {
    client := &http.Client{}
    var data = strings.NewReader(`{"name":"","pageIndex":0,"pageSize":25,"liquidityAdded":true,"liquidityRemoved":false}`)
    pageIndex := 0
    pageSize := 25
    hasMorePages := true
	var tokenCount = 0
    
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/tokendb")
    if err != nil {
        log.Fatal(err)
		return 

    }

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tokentb (
		id VARCHAR(36) ,
		address TEXT NOT NULL,
		block BIGINT NOT NULL,
		create_date TEXT NOT NULL,
		creator TEXT NOT NULL,
		deployer_wallet_age FLOAT NOT NULL,
		hex_code_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		symbol TEXT NOT NULL,
		decimals BIGINT NOT NULL,
		total_supply TEXT NOT NULL,
		start_liquidity FLOAT NOT NULL,
		start_token_reserve TEXT NOT NULL,
		current_liquidity FLOAT NOT NULL,
		current_token_reserve TEXT NOT NULL,
		start_mcap FLOAT NOT NULL,
		current_mcap FLOAT NOT NULL,
		current_gain FLOAT NOT NULL,
		ath_mcap FLOAT NOT NULL,
		ath_gain FLOAT NOT NULL,
		liquidity_added BOOLEAN NOT NULL,
		liquidity_removed BOOLEAN NOT NULL,
		could_buy BOOLEAN NOT NULL,
		could_sell BOOLEAN NOT NULL,
		delayed_honeypot BOOLEAN NOT NULL,
		buy_tax BIGINT NOT NULL,
		sell_tax BIGINT NOT NULL,
		max_buy TEXT NOT NULL,
		update_block BIGINT NOT NULL,
		v BIGINT NOT NULL,
		liquidity_in_usd FLOAT NOT NULL,
		max_buy_percentage FLOAT NOT NULL,
		max_buy_tokens TEXT NOT NULL,
		owner TEXT,
		pair_address TEXT NOT NULL,
		dead_tokens TEXT NOT NULL,
		dead_tokens_percentage FLOAT NOT NULL,
		liquidity_ratio FLOAT NOT NULL,
		pair_token TEXT NOT NULL,
		pair_token_name TEXT NOT NULL,
		init_max_buy TEXT,
		init_max_buy_percentage FLOAT,
		init_max_buy_tokens TEXT,
		launch_time TEXT,
		launched BOOLEAN,
		liquidity_burn_desc TEXT,
		liquidity_burn_percent FLOAT,
		liquidity_burn_tx_hash TEXT,
		liquidity_burned BOOLEAN,
		liquidity_lock_days FLOAT,
		liquidity_lock_desc TEXT,
		liquidity_lock_tx_hash TEXT,
		liquidity_locked BOOLEAN,
		liquidity_unlock_date TEXT,
		renounced BOOLEAN
	);
`
 
    _, err = db.Exec(createTableQuery)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Table created successfully!")

	sqlStatement := `
	INSERT INTO tokentb (id, address, block, create_date, creator, deployer_wallet_age, hex_code_hash, name, symbol, decimals, total_supply, start_liquidity, start_token_reserve, current_liquidity, current_token_reserve, start_mcap, current_mcap, current_gain, ath_mcap, ath_gain, liquidity_added, liquidity_removed, could_buy, could_sell, delayed_honeypot, buy_tax, sell_tax, max_buy, update_block, v, liquidity_in_usd, max_buy_percentage, max_buy_tokens, owner, pair_address, dead_tokens, dead_tokens_percentage, liquidity_ratio, pair_token, pair_token_name, init_max_buy, init_max_buy_percentage, init_max_buy_tokens, launch_time, launched, liquidity_burn_desc, liquidity_burn_percent, liquidity_burn_tx_hash, liquidity_burned, liquidity_lock_days, liquidity_lock_desc, liquidity_lock_tx_hash, liquidity_locked, liquidity_unlock_date, renounced)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55);
`
    for hasMorePages {
        req, err := http.NewRequest("POST", "https://dexspy.io/tokens", data)
        if err != nil {
            log.Fatal(err)
        }
        req.Header.Set("authority", "dexspy.io")
        req.Header.Set("accept", "application/json")
        req.Header.Set("accept-language", "en-US,en;q=0.9,es-ES;q=0.8,es;q=0.7")
        req.Header.Set("authorization", "undefined")
        req.Header.Set("content-type", "application/json")
        req.Header.Set("cookie", "_ga=GA1.1.2015534493.1683312341; _ga_8ZBN5RRG8E=GS1.1.1683312340.1.1.1683312576.0.0.0")
        req.Header.Set("origin", "https://dexspy.io")
        req.Header.Set("referer", "https://dexspy.io/eth/tokens-db")
        req.Header.Set("sec-ch-ua", `"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`)
        req.Header.Set("sec-ch-ua-mobile", "?0")
        req.Header.Set("sec-ch-ua-platform", `"macOS"`)
        req.Header.Set("sec-fetch-dest", "empty")
        req.Header.Set("sec-fetch-mode", "cors")
        req.Header.Set("sec-fetch-site", "same-origin")
        req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
        
        // Set the page index and size in the request data
        requestData := fmt.Sprintf(`{"name":"","pageIndex":%d,"pageSize":%d,"liquidityAdded":true,"liquidityRemoved":false}`, pageIndex, pageSize)
        data = strings.NewReader(requestData)
        resp, err := client.Do(req)
        if err != nil {
            log.Fatal(err)
        }
        defer resp.Body.Close()
        bodyText, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }

        var dexSpyResponse DexSpyResponse
        err = json.Unmarshal(bodyText, &dexSpyResponse)
        if err != nil {
            log.Fatal(err)
        }
        
        // If there are no more items in the response, break out of the loop
        if len(dexSpyResponse.Items) == 0 {
            break
        }
		if len(dexSpyResponse.Items) == 0 {
            hasMorePages = false
        }
        
        for _, item := range dexSpyResponse.Items {
			tokenCount ++
			_, err = db.Exec(sqlStatement, item.ID, item.Address, item.Block, item.CreateDate, item.Creator, item.DeployerWalletAge, 
				item.HexCodeHash, item.Name, item.Symbol, item.Decimals, item.TotalSupply, item.StartLiquidity, item.StartTokenReserve,
				item.CurrentLiquidity, item.CurrentTokenReserve, item.StartMcap, item.CurrentMcap, item.CurrentGain, item.AthMcap, item.AthGain, 
				item.LiquidityAdded, item.LiquidityRemoved, item.CouldBuy, item.CouldSell, item.DelayedHoneypot, item.BuyTax, item.SellTax, 
				item.MaxBuy, item.UpdateBlock, item.V, item.LiquidityInUsd, item.MaxBuyPercentage, item.MaxBuyTokens, item.Owner, item.PairAddress, 
				item.DeadTokens, item.DeadTokensPercentage, item.LiquidityRatio, item.PairToken, item.PairTokenName, item.InitMaxBuy, 
				item.InitMaxBuyPercentage, item.InitMaxBuyTokens, item.LaunchTime, item.Launched, item.LiquidityBurnDesc, item.LiquidityBurnPercent, 
				item.LiquidityBurnTxHash, item.LiquidityBurned, item.LiquidityLockDays, item.LiquidityLockDesc, item.LiquidityLockTxHash, 
				item.LiquidityLocked, item.LiquidityUnlockDate, item.Renounced)

			if err != nil {
				panic(err)
			}
            // fmt.Println(item[0], " and " , item[1])
        }
        
        pageIndex++
    }
	defer db.Close()
	fmt.Println("token count = ", tokenCount)
}