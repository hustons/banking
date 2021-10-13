package transactions

import (
  "fmt"
  "log"
  "strings"
  "time"
  "strconv"
  "math"
  "regexp"
  "github.com/shuston/banking/utils"
)

type revolutTransaction struct {
  id string
  realDate string
  completedDate time.Time
  description string
  amount float64
  category string
  notes string
  details string
  source string
  isHidden bool
}

var re = regexp.MustCompile(`\"([0-9]){1,3},([0-9]){3}\.([0-9]){2}\"+`)

func NewRevolutTransaction(rawData string) *revolutTransaction {
  t := new(revolutTransaction)

  rawData = t.preProcessData(rawData)
  data := t.splitData(rawData)

  completedDate, err := time.Parse("2 Jan 2006", data[0])
  if err != nil {
    log.Fatal("Could not parse date: ", err, "\nRaw data: ", rawData)
  }

  t.completedDate = completedDate
  t.description = strings.Trim(data[1], " ")

  paidOut := strings.Replace(data[2], ",", "", -1)
  paidIn := strings.Replace(data[3], ",", "", -1)

  if paidOut != "" {
    t.amount, err = strconv.ParseFloat(paidOut, 64)
    if err != nil {
      log.Fatal("Could not parse date: ", err, "\nRaw data: ", rawData)
    }
  } else {
    t.amount, err = strconv.ParseFloat(paidIn, 32)
    if err != nil {
      log.Fatal("Could not parse amount: ", err, "\nRaw data: ", rawData)
    }
    t.amount = -1 * t.amount
  }

  t.category = data[7]
  t.notes = data[8]

  t.details = t.description
  if t.notes != "" {
    t.details = t.details + " - " + t.notes
  }

  t.source = "Revolut"
  t.isHidden = t.shouldHide()
  t.id = utils.Classify(t.details)

  return t
}

func (t revolutTransaction) preProcessData(rawData string) string {
  // Handle additional commas in Details
  if strings.Contains(rawData, "\"\"") {
    split := strings.Split(rawData, "\"\"")
    split[1] = strings.Replace(split[1], ",", ";", -1)
    rawData = strings.Join(split, "\"\"")
  }

  // Handle Fees
  rawData = strings.Replace(rawData, ", Fee: ", "; Fee: ", -1)

  // Handle amounts >= 1,000.00
  matches := re.FindAllString(rawData, -1)
  for _, match := range matches {
    split := strings.Split(rawData, match)
    match = strings.Replace(match, "\"", "", -1)
    match = strings.Replace(match, ",", "", -1)
    rawData = strings.Join(split, match)
  }

  return rawData
}

func (t revolutTransaction) splitData(rawData string) []string {
  // Split and trim raw data
  data := strings.Split(rawData, ",")
  for i := range data {
    data[i] = strings.Trim(data[i], " ")
  }

  // Check for unexpected / differently formatted data
  if len(data) > 9 {

    // Attempt to recover for balances >= 1,000.00
    if len(data) == 10 && data[9] == "" && data[8] != "" {
      _, err1 := strconv.ParseInt(data[6], 10, 32)
      _, err2 := strconv.ParseFloat(data[7], 32)
      _, err3 := strconv.ParseFloat(data[8], 32) // This should fail

      if err1 != nil || err2 != nil || err3 == nil {
        log.Fatal("Unexpected format for input data.\nRaw data:", rawData)
      }

      data[6] = data[6] + data[7]
      data[7] = data[8]
      data[8] = ""
    } else {
      log.Fatal("Unexpected format for input data.\nRaw data:", rawData)
    }
  }

  return data
}

func (t revolutTransaction) GetCompletedDate() time.Time {
  return t.completedDate
}

func (t revolutTransaction) Output() {
  if t.isHidden {
    return
  }

  roundedAmount := math.Round(t.amount*100)/100
  strAmount := fmt.Sprintf("%f", roundedAmount)
  fmt.Println(t.id + "\t" +
    t.realDate + "\t" +
    t.completedDate.Format("2 Jan") + "\t" +
    strAmount + "\t" +
    t.details + "\t" +
    t.source)
}

func (t revolutTransaction) shouldHide() bool {
  return t.details == "Top-Up by *5278" ||
         strings.HasPrefix(t.details, "Money added via ··5278") ||
         t.details == "Money added via Google Pay"
}
