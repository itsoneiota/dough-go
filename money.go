// Package dough provides arithmetic for monetary amounts.
package dough

import (
	"fmt"
	"golang.org/x/text/currency"
	"math"
	"regexp"
	"strconv"
)

// Money is a value object representing a monetary amount.
type Money struct {
	// Currency
	c currency.Unit
	// Atoms, the amount in the smallest unit of the given currency.
	a int
}

// New returns a new Money instance for the given currency and amount.
// cur is an 3-letter ISO 4217 currency code.
// amt is a string representation of the amount, e.g. "123.45".
// It returns an error if cur is not well formed or not recognised,
// or if amt cannot be parsed.
func New(cur, amt string) (Money, error) {
	c, err := currency.ParseISO(cur)
	if err != nil {
		return Money{}, fmt.Errorf("coudn't parse currency: %v", err)
	}

	a, err := strToInt(c, amt)
	if err != nil {
		return Money{}, fmt.Errorf("couldn't parse amount: %v", err)
	}
	return Money{
		c: c,
		a: a,
	}, nil
}

func strToInt(c currency.Unit, amt string) (int, error) {
	// TODO: Capture sub-units based on currency exponent.
	// https://en.wikipedia.org/wiki/ISO_4217#Treatment_of_minor_currency_units_.28the_.22exponent.22.29
	re := regexp.MustCompile("^(-)?(\\d+)(\\.([\\d]{2}))?$")
	m := re.FindStringSubmatch(amt)
	if len(m) == 0 {
		return 0, fmt.Errorf("unable to parse amount: %s", amt)
	}
	digits := m[2] + m[4]
	a, err := strconv.Atoi(digits)
	if err != nil {
		return 0, fmt.Errorf("unable to parse amount: %v", err)
	}
	if m[1] == "-" {
		a *= -1
	}
	return a, nil
}

// Currency gets the currency of the Money.
func (x Money) Currency() string {
	return x.c.String()
}

// Amount gets the currency of the Money.
func (x Money) Amount() string {
	neg := ""
	a := x.a
	if a < 0 {
		neg = "-"
		a *= -1
	}
	maj := strconv.Itoa(a / 100) // TODO: Variable
	min := fmt.Sprintf("%02d", a%100)

	return neg + maj + "." + min
}

// Add returns a new Money with the value of the given Money added.
func (x Money) Add(y Money) (Money, error) {
	return addSub(x, y, true)
}

// Sub returns a new Money with the value of the given Money added.
func (x Money) Sub(y Money) (Money, error) {
	return addSub(x, y, false)
}

func addSub(x, y Money, add bool) (Money, error) {
	if x.Currency() != y.Currency() {
		var op string
		if add {
			op = "add"
		} else {
			op = "subtract"
		}
		err := fmt.Errorf("Can't %s different currencies. Attempting to add %s and %s", op, x.Currency(), y.Currency())
		return Money{}, err
	}
	var z int
	if add {
		z = x.a + y.a
	} else {
		z = x.a - y.a
	}
	return Money{
		x.c,
		z,
	}, nil
}

// Mul returns a new Money with the value of m multiplied by factor.
func (x Money) Mul(f int) (Money, error) {
	return Money{
		x.c,
		x.a * f,
	}, nil
}

// Cmp compares x and y and returns:
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x Money) Cmp(y Money) (c int, err error) {
	if x.Currency() != y.Currency() {
		err := fmt.Errorf("Can't compare different currencies (%s and %s)", x.Currency(), y.Currency())
		return 0, err
	}
	if x.a < y.a {
		c = -1
	} else if x.a == y.a {
		c = 0
	} else {
		c = 1
	}

	return
}

// Share allocates portions of a Money's value between parties based on weightings given.
// Spare pennies are distributed among parties evenly, from first to last.
func (x Money) Share(weightings []uint) []Money {
	n := len(weightings)
	var sum uint
	for _, w := range weightings {
		sum += w
	}
	if sum == 0 {
		for i := range weightings {
			weightings[i] = 1
		}
		sum = uint(n)
	}
	ratios := make([]float64, n)
	for i := range weightings {
		ratios[i] = float64(weightings[i]) / float64(sum)
	}

	allocations := make([]int, n)
	fa := float64(x.a)
	rem := x.a
	for i := range ratios {
		a := int(math.Trunc(ratios[i] * fa))
		allocations[i] = a
		rem -= a
	}
	d := 1
	if rem < 0 {
		d = -1
	}
	for i := 0; rem != 0; i++ {
		ind := i % n
		if weightings[ind] == 0 {
			continue
		}
		allocations[ind] += d
		rem += (-d)
	}

	// Double-check allocation to make sure we haven't made or lost pennies.
	// It would be _very_ bad to get this wrong.
	total := 0
	for i := range allocations {
		total += allocations[i]
	}
	if total != x.a {
		panic(fmt.Sprintf("dough package: bad allocation. Started with %d atoms, allocated %d. Weightings=%v", x.a, total, weightings))
	}

	res := make([]Money, len(allocations))
	for i := range allocations {
		res[i] = Money{
			x.c,
			allocations[i],
		}
	}

	return res
}
