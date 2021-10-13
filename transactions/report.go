package transactions

import (
  "fmt"
  "time"
  "sort"
)

type report struct {
  titles string
  startDate time.Time
  transactions []transaction
}

func NewReport(startDate time.Time) *report {
  r := new(report)
  r.titles = "ID\t\tReal Date\tCompleted\tAmount\tDetails\tSource"
  r.startDate = startDate
  return r
}

func (r report) AddTransactions(rawData[]string, isRevolut bool) *report {
  if isRevolut {
    for i := len(rawData) - 1; i > 0; i-- {
      transaction := NewRevolutTransaction(rawData[i])
      if transaction.GetCompletedDate().After(r.startDate) {
        r.transactions = append(r.transactions, *transaction)
      }
    }
  } else {
    for i := len(rawData) -1; i > 0; i-- {
      transaction := NewAIBTransaction(rawData[i])
      if transaction.GetCompletedDate().After(r.startDate) {
        r.transactions = append(r.transactions, *transaction)
      }
    }
  }
  return &r
}

func (r report) Sort() {
  sort.Sort(byDate(r.transactions))
}

func (r report) Output() {
  fmt.Println(r.titles)
  for i := 0; i < len(r.transactions); i++ {
    r.transactions[i].Output()
  }
}
