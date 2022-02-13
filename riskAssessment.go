package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TransactionRecords struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID            int `json:"id"`
	UserID        int `json:"user_id"`
	AmountUsCents int `json:"amount_us_cents"`
	CardID        int `json:"card_id"`
}

type ResponseOutput struct {
	RiskRatings []string `json:"riskRatings"`
}

func main() {
	handler := http.HandlerFunc(riskAssessmentHandler)
	http.Handle("/risk_assessment", handler)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func convertCentsToDollars(cents int) int {
	return cents / 100
}

func riskAssessmentHandler(w http.ResponseWriter, r *http.Request) {
	var transactionRecords TransactionRecords
	err := json.NewDecoder(r.Body).Decode(&transactionRecords)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(riskAssessment(transactionRecords))
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}

func riskAssessment(transactionData TransactionRecords) ResponseOutput {
	//fmt.Printf("test")
	var riskRatings []string
	for transaction := 0; transaction < len(transactionData.Transactions); transaction++ {
		var totalAmountSpent int
		var totalMoneyInDollars int
		var transactionMoneyInDollars int
		var cardsBeingUsed []int

		if transaction > 0 {
			for value := 0; value <= transaction; value++ {
				if transactionData.Transactions[transaction].UserID == transactionData.Transactions[value].UserID {
					if len(cardsBeingUsed) == 0 {
						cardsBeingUsed = append(cardsBeingUsed, transactionData.Transactions[transaction].CardID)
					}
					totalAmountSpent += transactionData.Transactions[value].AmountUsCents
					if transactionData.Transactions[transaction].CardID != transactionData.Transactions[value].CardID {
						cardsBeingUsed = append(cardsBeingUsed, transactionData.Transactions[value].CardID)
					}
				}
			}
			totalMoneyInDollars = convertCentsToDollars(totalAmountSpent)
			transactionMoneyInDollars = convertCentsToDollars(transactionData.Transactions[transaction].AmountUsCents)
		} else {
			totalMoneyInDollars = convertCentsToDollars(transactionData.Transactions[transaction].AmountUsCents)
			transactionMoneyInDollars = totalMoneyInDollars
		}
		if len(cardsBeingUsed) > 2 || totalMoneyInDollars > 20000 || transactionMoneyInDollars > 10000 {
			riskRatings = append(riskRatings, "high")
		} else if len(cardsBeingUsed) == 2 || totalMoneyInDollars > 10000 || transactionMoneyInDollars > 5000 {
			riskRatings = append(riskRatings, "medium")
		} else {
			riskRatings = append(riskRatings, "low")
		}
	}
	//fmt.Printf("eof")
	return ResponseOutput{
		riskRatings,
	}
}
