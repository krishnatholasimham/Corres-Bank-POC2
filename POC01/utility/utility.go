/*
Copyright (C) Australia and New Zealand Banking Group Limited (ANZ)
100 Queen Street, Melbourne 3000, ABN 11 005 357 522.
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by Chris T'en <chris.ten@anz.com> March 2016
*/

package utility

import (
	// "crypto/rand"
	"crypto/sha256"
	"math/rand"
	"strings"
	"time"

	"github.com/op/go-logging"
)

// GenerateKey generates a unique ID by hashing the input message.
// TODO: Will overwrite existing record if two identical MT103 messages
// are processed.
func GenerateKey(message string) []byte {
	hash := sha256.New()
	hash.Write([]byte(message))
	h := hash.Sum(nil)

	return h
}

// GenerateRandomNumber generates a random number. Used to populate fields
// with dummy data.
func GenerateRandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("1234567890")

	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// IdentifyLedgerType takes a value entry from a ledger and returns a string
// specifying the type of entry (i.e. Request, Confirmation or Funding).
func IdentifyLedgerType(value []byte) string {
	ledgerType := ""
	if strings.Contains(string(value), "REQUEST-RECORD") {
		ledgerType = "REQUEST-RECORD"
	} else if strings.Contains(string(value), "CONFIRMATION-RECORD") {
		ledgerType = "CONFIRMATION-RECORD"
	} else if strings.Contains(string(value), "FUNDING-RECORD") {
		ledgerType = "FUNDING-RECORD"
	} else if strings.Contains(string(value), "FEE-RECORD") {
		ledgerType = "FEE-RECORD"
	} else if strings.Contains(string(value), "DIRECT-CREDIT-RECORD") {
		ledgerType = "DIRECT-CREDIT-RECORD"
	}

	return ledgerType
}

func ToDigits(num int) []int {
	n := Abs(num)
	s := make([]int, 0)
	for n > 0 {
		s = append(s, n%10)
		n /= 10
	}
	return s
}

func CountDigitNoMatches(num int, digit int) int { return CountNoMatches(ToDigits(num), digit) }
func CountDigitMatches(num int, digit int) int   { return CountMatches(ToDigits(num), digit) }
func CountMatches(nums []int, num int) int       { return countImpl(nums, num, true) }
func CountNoMatches(nums []int, num int) int     { return countImpl(nums, num, false) }
func countImpl(nums []int, num int, exist bool) int {
	c := 0
	for i := 0; i < len(nums); i++ {
		if (exist && nums[i] == num) || (!exist && nums[i] != num) {
			c += 1
		}
	}
	return c
}

// Abs for int oddly missing from golang: http://osdir.com/ml/go-language-discuss/2013-03/msg02848.html
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func InitDefaultLogging() {
	InitLogging("%{color}%{shortfunc:-10.10s} %{level:-3.3s} %{id:03x}%{color:reset} - %{message}")
}

func InitLogging(f string) {
	var format = logging.MustStringFormatter(f)
	logging.SetFormatter(format)
}
