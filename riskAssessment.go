package main

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

	//fmt.Printf("hey")
}

func convertCentsToDollars(cents int) int {
	return cents / 100
}

func riskAssessment(transactionData TransactionRecords) ResponseOutput {
	var riskRatings []string
	for transaction := 0; transaction < len(transactionData.Transactions); transaction++ {
		var totalAmountSpent int
		var totalMoneyInDollars int
		var transactionMoneyInDollars int
		var cardsBeingUsed []int

		if transaction > 0 {
			for value := 0; value < transaction; value++ {
				if transactionData.Transactions[transaction].UserID == transactionData.Transactions[value].UserID {
					if len(cardsBeingUsed) == 0 {
						cardsBeingUsed = append(cardsBeingUsed, transactionData.Transactions[transaction].CardID)
					}
					totalAmountSpent += transactionData.Transactions[value].AmountUsCents
					if transactionData.Transactions[transaction].CardID != transactionData.Transactions[value].CardID {
						cardsBeingUsed = append(cardsBeingUsed, transactionData.Transactions[transaction].CardID)
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
	return ResponseOutput{
		riskRatings,
	}
}
