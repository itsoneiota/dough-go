// Package dough provides arithmetic for monetary amounts.
package dough

import (
	"fmt"
	"math/big"

	"golang.org/x/text/currency"
)

// Money is a value object representing a monetary amount.
type Money struct {
	// Currency
	c currency.Unit
	// Amount
	a *big.Float
}

// New returns a new Money instance for the given currency and amount.
// cur is an 3-letter ISO 4217 currency code.
// amt is a string representation of the amount, e.g. "123.45".
// It returns an error if cur is not well formed or not recognised,
// or if amt cannot be parsed as a float.
func New(cur, amt string) (Money, error) {
	c, err := currency.ParseISO(cur)
	if err != nil {
		return Money{}, fmt.Errorf("coudn't parse currency: %v", err)
	}
	a, _, err := big.ParseFloat(amt, 10, 0, big.ToNearestAway)
	if err != nil {
		return Money{}, fmt.Errorf("couldn't parse amount: %v", err)
	}
	return Money{
		c: c,
		a: a,
	}, nil
}

// Currency gets the currency of the Money.
func (x Money) Currency() string {
	return x.c.String()
}

// Amount gets the currency of the Money.
func (x Money) Amount() string {
	return x.a.Text('f', 2)
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
	var z big.Float
	if add {
		z.Add(x.a, y.a)
	} else {
		z.Sub(x.a, y.a)
	}
	return Money{
		x.c,
		&z,
	}, nil
}

// Mul returns a new Money with the value of m multiplied by factor.
func (x Money) Mul(factor int64) (Money, error) {
	var ff, z big.Float
	ff.SetInt64(factor)
	return Money{
		x.c,
		z.Mul(x.a, &ff),
	}, nil
}

// Cmp compares x and y and returns:
//	-1 if x <  y
//	 0 if x == y (incl. -0 == 0, -Inf == -Inf, and +Inf == +Inf)
//	+1 if x >  y
func (x Money) Cmp(y Money) (int, error) {
	if x.Currency() != y.Currency() {
		err := fmt.Errorf("Can't compare different currencies (%s and %s)", x.Currency(), y.Currency())
		return 0, err
	}
	return x.a.Cmp(y.a), nil
}
