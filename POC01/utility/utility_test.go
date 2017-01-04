package utility

import (
	"github.com/op/go-logging"
	"testing"
	"strconv"

	"math/big"
	"math"
)

func TestLast3(t *testing.T) {
	InitDefaultLogging()

	var log = logging.MustGetLogger("test")
	log.Debugf("debug %s", "arg")
	log.Errorf("error")

	amount, _ := strconv.ParseFloat("-100546.295", 64)

	n := []float64{amount, 100547.346, 108546.742}
	log.Infof("1: %f, 2: %f, 3: %f\n", n[0], n[1], n[2])

	f := []*big.Float{big.NewFloat(n[0]), big.NewFloat(n[1]), big.NewFloat(n[2])}
	log.Infof("1: %f, 2: %f, 3: %f\n", f[0], f[1], f[2])

	number, accuracy := f[0].Int64()
	num := number
	log.Infof("1: %d, %s", num, accuracy)
	log.Infof("1: %d", num % 10)
	log.Infof("1: %d", num / 10 % 10)
	log.Infof("1: %d", num / 100 % 10)
	for num > 0 {
		log.Info("1: %d, %d", num % 10, num)
		num /= 10
	}

	log.Infof("rem: %f", math.Mod(float64(number), 1000))
	last3 := number % 1000
	log.Infof("rem: %d", last3)

	log.Infof("Now calling toDigits")
	digits := ToDigits(int(number))
	log.Infof("3: %v, len: %d, #zeros: %d, #non zeros: %d", digits, len(digits), CountMatches(digits, 0), CountNoMatches(digits, 0))
}

