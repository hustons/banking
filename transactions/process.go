package transactions

import (
  "fmt"
  "time"
  "github.com/shuston/banking/utils"
)

func Process(revolutFile string, aibFile string, startDate time.Time) {
  fmt.Println("Excluding transactions on and before ", startDate)
  report := NewReport(startDate)

  revolutLines := utils.ReadInputFile(revolutFile)
  aibLines := utils.ReadInputFile(aibFile)

  report = report.AddTransactions(revolutLines)
  report = report.AddTransactions(aibLines)

  report.Sort()
  report.Output()
}
