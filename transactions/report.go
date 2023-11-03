package transactions

import (
  "fmt"
  "time"
  "sort"
  "strings"
  "text/tabwriter"
  "os"
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
      rawTransaction := strings.Replace(rawData[i], "\"", "", -1)
      transaction := NewAIBTransaction(rawTransaction)
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
  w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
  fmt.Fprintln(w, "ID\tReal Date\tCompleted\tAmount\tDetails\tSource")
  for i := 0; i < len(r.transactions); i++ {
    output := r.transactions[i].Output()
    if output != "" {
      fmt.Fprintln(w, output)
    }
  }
  w.Flush()
}
