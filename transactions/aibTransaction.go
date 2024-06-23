package transactions

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/stonish/banking/utils"
)

// type creditAndDebitAmountError struct {}
//
//	func (e *creditAndDebitAmountError) Error() string {
//	 return "Transaction with both credit and debit amounts found"
//	}
var blerg = errors.New("Transaction with both credit and debit amounts found")

func errCreditAndDebitAmount() error { return blerg }

var creditAndDebitAmountError = errCreditAndDebitAmount()

type aibTransaction struct {
	id            string
	realDate      string
	completedDate time.Time
	details       string
	amount        float64
	balance       string
	source        string
	isHidden      bool
}

func NewAIBTransaction(rawData string) *aibTransaction {
	t := new(aibTransaction)

	data := t.splitData(rawData)

	// 0 Account
	// 1 TransactionDate
	// 2-3-4 Description (1-2-3)
	// 5 Debit Amount
	// 6 Credit Amount
	// 7 Balance
	// 8 Posted Currency
	// 9 Transaction Type
	// 10 Local Currency Amount
	// 11 Local Currency

	transactionDate, err := time.Parse("02/01/2006", data[1])
	if err != nil {
		log.Fatal("Could not parse date: ", err, "\nRaw data: ", rawData)
	}
	t.completedDate = transactionDate

	initialDetail := strings.Replace(data[2], "VDC-", "", -1)
	initialDetail = strings.Replace(initialDetail, "VDP-", "", -1)
	details := []string{strings.Trim(initialDetail, " "), strings.Trim(data[3], " "), strings.Trim(data[4], " ")}
	t.details = strings.Join(details, " ")
	t.details = strings.Trim(t.details, " ")

	t.amount, err = t.parseAmount(data[5], data[6])
	if err != nil {
		log.Fatal("Could not parse amount: ", err, "\nRaw data: ", rawData)
	}

	t.balance = data[7]
	t.source = "AIB"
	t.isHidden = t.shouldHide()
	t.id = utils.Classify(t.details)

	return t
}

func (t aibTransaction) splitData(rawData string) []string {
	data := strings.Split(rawData, ",")
	if len(data) == 12 {
		return data
	}

	log.Printf("AIB transaction received with unexpected number of inputs, attempting to recover")

	_, err := t.parseAmount(data[5], data[6])
	if errors.Is(err, creditAndDebitAmountError) {
		log.Printf("Additional comma on debit amount detected")
		data[5] = data[5] + data[6]
		for i := 6; i < 11; i++ {
			data[i] = data[i+1]
		}
		return data
	}

	_, errCredit := strconv.ParseFloat(data[6], 32)
	_, errBalance := strconv.ParseFloat(data[7], 32)
	_, errCurrency := strconv.ParseFloat(data[8], 32)
	if data[5] == "" && errCredit == nil && errBalance == nil && errCurrency == nil {
		log.Printf("Additional comma on credit amount detected")
		data[6] = data[6] + data[7]
		for i := 7; i < 11; i++ {
			data[i] = data[i+1]
		}
		return data
	}
	log.Fatal("Could not recover\nRaw data: ", rawData)
	return data
}

func (t aibTransaction) GetCompletedDate() time.Time {
	return t.completedDate
}

func (t aibTransaction) Output() string {
	if t.isHidden {
		return ""
	}

	roundedAmount := math.Round(t.amount*100) / 100
	strAmount := fmt.Sprintf("%f", roundedAmount)
	return t.id + "\t" +
		t.realDate + "\t" +
		t.completedDate.Format("2 Jan") + "\t" +
		strAmount + "\t" +
		t.details + "\t" +
		t.source
}

func (t aibTransaction) parseAmount(debitAmount string, creditAmount string) (float64, error) {
	var amount float64

	if debitAmount != "" && creditAmount != "" {
		return amount, creditAndDebitAmountError
	}

	if debitAmount != "" {
		amount, err := strconv.ParseFloat(debitAmount, 32)

		if err != nil {
			return amount, fmt.Errorf("Could not parse amount: %w", err)
		}

		return amount, nil
	}

	if creditAmount != "" {
		amount, err := strconv.ParseFloat(creditAmount, 32)

		if err != nil {
			return amount, fmt.Errorf("Could not parse amount: %w", err)
		}

		if amount != 0 {
			amount = -1 * amount
		}

		return amount, nil
	}

	return amount, errors.New("Transaction with neither credit nor debit amounts found")
}

func (t aibTransaction) shouldHide() bool {
	if t.amount == 0 {
		return true
	}

	if t.isRevolutTopUp() {
		return true
	}

	if strings.HasPrefix(t.details, "*INET SAVINGS ") {
		return true
	}

	if t.amount == 600 &&
		(strings.HasPrefix(t.details, "*INET RENT ") ||
			strings.HasPrefix(t.details, "*INET DAVID ")) {
		return true
	}

	if t.amount == 10 && t.details == "931365 22689017" {
		return true
	}

	return false
}

func (t aibTransaction) isRevolutTopUp() bool {
	return strings.HasPrefix(t.details, "Revolut**") ||
		strings.HasPrefix(t.details, "Revolut* - ") ||
		strings.HasPrefix(t.details, "Revolut  - ") ||
		strings.HasPrefix(t.details, "REVOLUT*") ||
		t.details == "Revolut"
}
