package utils

import (
  "strings"
  "os"
  "io/ioutil"
  "encoding/json"
  "log"
)

var classifications map[string]string
var desiredLength = 8

func Classify(details string) string {
  if classifications == nil {
    err := loadClassifications()
    if err != nil {
      log.Fatal("Could not load classifications: ", err)
    }
  }

  classification, exists := classifications[details]
  if (exists) {
    return pad(classification)
  }

  if strings.HasPrefix(details, "Dominos Pizza Uk") {
    return pad("Take Away")
  }

  if strings.HasPrefix(details, "VDC-TESCO STORES") {
    return pad("Groceries")
  }

  if strings.HasPrefix(details, "Amazon Prime FX Rate") ||
     strings.HasPrefix(details, "Amazon FX Rate €1 = £") ||
     strings.HasPrefix(details, "Amazon Prime*mk2yu0v14 FX Rate €1 = £") ||
     strings.HasPrefix(details, "FEE-QTR TO") ||
     strings.HasPrefix(details, "Expressvpn") ||
     strings.HasPrefix(details, "VDP-LINKEDIN") ||
     strings.HasPrefix(details, "VDP-LinkedIn") {
    return pad("Remainder")
  }

  if strings.HasPrefix(details, "VDP-Spotify P") {
    return pad("Spotify")
  }

  if strings.HasPrefix(details, "D/D CLOSE BROTHERS") {
    return pad("Car Loan")
  }

  if strings.HasPrefix(details, "VDC-APPLEGREEN") ||
     strings.HasPrefix(details, "VDC-CIRCLE K") {
    return pad("Petrol")
  }

  return pad("UNKNOWN")
}

func loadClassifications() error {
  file, err := os.Open("./classifications.json")
  if err != nil {
    return err
  }
  defer file.Close()

  byteValue, err := ioutil.ReadAll(file)
  if err != nil {
    return err
  }

  err = json.Unmarshal(byteValue, &classifications)
  if err != nil {
    return err
  }

  return nil
}

func pad(classification string) string {
  for i := desiredLength - len(classification); i > 0; i-- {
    classification += " "
  }
  return classification
}
