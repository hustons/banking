package transactions

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shuston/banking/utils"
)

type revolutTransaction struct {
	id              string
	transactionType string
	product         string
	startedDate     time.Time
	completedDate   time.Time
	description     string
	amount          float64
	fee             float64
	currency        string
	state           string
	source          string
	balance         float64
	isHidden        bool
}

func NewRevolutTransaction(rawData string) *revolutTransaction {
	data := strings.Split(rawData, ",")

	t := new(revolutTransaction)
	t.transactionType = data[0]
	t.product = data[1]

	startedDate, err := time.Parse("2006-01-02 15:04:05", data[2])
	if err != nil {
		log.Fatal("Could not parse date: ", err, "\nRaw data: ", rawData)
	}
	t.startedDate = startedDate

	completedDateStr := data[3]
	if completedDateStr != "" {
		completedDate, err := time.Parse("2006-01-02 15:04:05", completedDateStr)
		if err != nil {
			log.Fatal("Could not parse date: ", err, "\nRaw data: ", rawData)
		}
		t.completedDate = completedDate
	}

	t.description = strings.Trim(data[4], " ")

	amount, err := strconv.ParseFloat(data[5], 64)
	if err != nil {
		log.Fatal("Could not parse amount: ", err, "\nRaw data: ", rawData)
	}

	fee, err := strconv.ParseFloat(data[6], 64)
	if err != nil {
		log.Fatal("Could not parse fee: ", err, "\nRaw data: ", rawData)
	}
	t.fee = fee
	t.amount = -amount + t.fee

	t.currency = data[7]
	t.state = data[8]

	balanceStr := data[9]
	if balanceStr != "" {
		balance, err := strconv.ParseFloat(data[9], 64)
		if err != nil {
			log.Fatal("Could not parse balance: ", err, "\nRaw data: ", rawData)
		}
		t.balance = balance
	}

	t.source = "Revolut"
	t.isHidden = t.shouldHide()
	t.id = utils.Classify(t.description)

	return t
}

func (t revolutTransaction) GetCompletedDate() time.Time {
	return time.Date(t.completedDate.Year(), t.completedDate.Month(), t.completedDate.Day(), 0, 0, 0, 0, time.UTC)
}

func (t revolutTransaction) Output() string {
	if t.isHidden {
		return ""
	}

	roundedAmount := math.Round(t.amount*100) / 100
	strAmount := fmt.Sprintf("%f", roundedAmount)
	return t.id + "\t" +
		t.startedDate.Format("2 Jan") + "\t" +
		t.completedDate.Format("2 Jan") + "\t" +
		strAmount + "\t" +
		t.description + "\t" +
		t.source
}

func (t revolutTransaction) shouldHide() bool {
	return t.amount == 0 || t.transactionType == "TOPUP" || t.state == "PENDING"
}
