package dough

import "testing"

func TestCanCreate(t *testing.T) {
	var cases = []struct {
		cur string
		amt string
	}{
		{"GBP", "0.00"},
		{"GBP", "0.01"},
		{"GBP", "-0.01"},
		{"GBP", "123.45"},
		{"AUD", "0.00"},
		{"AUD", "0.01"},
		{"AUD", "-0.01"},
		{"AUD", "123.45"},
	}
	for _, c := range cases {
		sut, err := New(c.cur, c.amt)
		if err != nil {
			t.Errorf("error received from New(\"%s\",\"%s\"), none expected %v", c.cur, c.amt, err)
		}
		if sut.Currency() != c.cur {
			t.Errorf("wanted %v, got %v", c.cur, sut.Currency())
		}
		if sut.Amount() != c.amt {
			t.Errorf("wanted %v, got %v", c.amt, sut.Amount())
		}
	}
}

func TestCanRejectBadCurrency(t *testing.T) {
	var cases = []struct {
		cur string
	}{
		{"FOO"},
		{"boogaloo"},
		{"$"},
		{"GB"},
	}
	for _, c := range cases {
		_, err := New(c.cur, "123.45")
		if err == nil {
			t.Errorf("error expected from New(\"%s\",\"123.45\"), none received", c.cur)
		}
	}
}

func TestCanRejectBadAmount(t *testing.T) {
	var cases = []struct {
		amt string
	}{
		{"Z"},
		{"10.S0"},
		{"ONE"},
		{"10 EUR"},
		{"1f.00"},
	}
	for _, c := range cases {
		_, err := New("GBP", c.amt)
		if err == nil {
			t.Errorf("error expected from New(\"GBP\",\"%s\"), none received", c.amt)
		}
	}
}

func TestCanAdd(t *testing.T) {
	var cases = []struct {
		a    string
		b    string
		want string
	}{
		{"0.00", "0.00", "0.00"},
		{"0.00", "0.01", "0.01"},
		{"0.01", "0.00", "0.01"},
		{"123.45", "234.56", "358.01"},
		{"-12.34", "12.34", "0.00"},
		{"45.67", "-45.67", "0.00"},
		{"-10.11", "30.33", "20.22"},
	}
	for _, c := range cases {
		a, _ := New("GBP", c.a)
		b, _ := New("GBP", c.b)
		if got, _ := a.Add(b); got.Amount() != c.want {
			t.Errorf("wanted %s, got %s", c.want, got.Amount())
		}
	}
}

func TestCanRejectMismatchedCurrencyWhenAdding(t *testing.T) {
	var cases = []struct {
		ac string
		bc string
	}{
		{"GBP", "AUD"},
		{"GBP", "EUR"},
		{"AUD", "GBP"},
	}
	for _, c := range cases {
		a, _ := New(c.ac, "1.00")
		b, _ := New(c.bc, "1.00")
		if _, err := a.Add(b); err == nil {
			t.Errorf("error expected when adding %s to %s, none received", c.ac, c.bc)
		}
	}
}

func TestCanSubtract(t *testing.T) {
	var cases = []struct {
		a    string
		b    string
		want string
	}{
		{"0.00", "0.00", "0.00"},
		{"0.01", "0.01", "0.00"},
		{"0.01", "0.00", "0.01"},
		{"358.01", "234.56", "123.45"},
		{"0.00", "12.34", "-12.34"},
		{"0.00", "-45.67", "45.67"},
		{"-10.11", "30.33", "-40.44"},
	}
	for _, c := range cases {
		a, _ := New("GBP", c.a)
		b, _ := New("GBP", c.b)
		if got, _ := a.Sub(b); got.Amount() != c.want {
			t.Errorf("wanted %s, got %s", c.want, got.Amount())
		}
	}
}

func TestCanRejectMismatchedCurrencyWhenSubtracting(t *testing.T) {
	var cases = []struct {
		ac string
		bc string
	}{
		{"GBP", "AUD"},
		{"GBP", "EUR"},
		{"AUD", "GBP"},
	}
	for _, c := range cases {
		a, _ := New(c.ac, "1.00")
		b, _ := New(c.bc, "1.00")
		if _, err := a.Sub(b); err == nil {
			t.Errorf("error expected when adding %s to %s, none received", c.ac, c.bc)
		}
	}
}

func TestCanMultiply(t *testing.T) {
	var cases = []struct {
		a    string
		f    int
		want string
	}{
		{"123.45", 1, "123.45"},
		{"123.45", 2, "246.90"},
		{"123.45", 10, "1234.50"},
		{"123.45", 100, "12345.00"},
		{"123.45", -1, "-123.45"},
		{"-123.45", -1, "123.45"},
		{"123.45", 0, "0.00"},
	}
	for _, c := range cases {
		sut, _ := New("GBP", c.a)
		if got, _ := sut.Mul(c.f); got.Amount() != c.want {
			t.Errorf("multiplying %s by %d. wanted %s, got %s", c.a, c.f, c.want, got.Amount())
		}
	}
}

func TestCanCompare(t *testing.T) {
	var cases = []struct {
		a    string
		b    string
		want int
	}{
		{"0.00", "0.00", 0},
		{"0.00", "0.01", -1},
		{"0.01", "0.00", 1},
		{"123.45", "234.56", -1},
		{"-12.34", "12.34", -1},
		{"45.67", "-45.67", 1},
		{"-10.11", "-30.33", 1},
		{"-40.11", "-30.33", -1},
	}
	for _, c := range cases {
		a, _ := New("GBP", c.a)
		b, _ := New("GBP", c.b)
		if got, _ := a.Cmp(b); got != c.want {
			t.Errorf("wanted %d, got %d", c.want, got)
		}
	}
}

func TestCanRejectMismatchedCurrencyWhenComparing(t *testing.T) {
	var cases = []struct {
		ac string
		bc string
	}{
		{"GBP", "AUD"},
		{"GBP", "EUR"},
		{"AUD", "GBP"},
	}
	for _, c := range cases {
		a, _ := New(c.ac, "1.00")
		b, _ := New(c.bc, "1.00")
		if _, err := a.Cmp(b); err == nil {
			t.Errorf("error expected when comparing %s to %s, none received", c.ac, c.bc)
		}
	}
}

func TestCanAllocate(t *testing.T) {
	var cases = []struct {
		a      string
		ratios []uint
		want   []string
	}{
		{"0.00", []uint{1, 1, 1}, []string{"0.00", "0.00", "0.00"}},
		{"0.01", []uint{1, 1, 1}, []string{"0.01", "0.00", "0.00"}},
		{"0.02", []uint{1, 1, 1}, []string{"0.01", "0.01", "0.00"}},
		{"0.03", []uint{1, 1, 1}, []string{"0.01", "0.01", "0.01"}},
		{"0.04", []uint{1, 1, 1}, []string{"0.02", "0.01", "0.01"}},
		{"0.05", []uint{1, 1, 1}, []string{"0.02", "0.02", "0.01"}},
		{"1.00", []uint{0, 1, 0}, []string{"0.00", "1.00", "0.00"}},
		{"0.03", []uint{0, 5, 0}, []string{"0.00", "0.03", "0.00"}},
		{"3.00", []uint{1, 1, 1}, []string{"1.00", "1.00", "1.00"}},
		{"1.00", []uint{1, 1, 1}, []string{"0.34", "0.33", "0.33"}},
		{"0.03", []uint{0, 5, 0}, []string{"0.00", "0.03", "0.00"}},
		{"0.03", []uint{0, 4, 2}, []string{"0.00", "0.02", "0.01"}},

		// Copied from MoneyTest.php
		{"1.05", []uint{3, 7}, []string{"0.32", "0.73"}},
		{"0.05", []uint{1, 1}, []string{"0.03", "0.02"}},
		{"300.00", []uint{122, 878}, []string{"36.60", "263.40"}},
		{"300.00", []uint{122, 0, 878}, []string{"36.60", "0.00", "263.40"}},
		{"120.00", []uint{20, 100}, []string{"20.00", "100.00"}},

		// One deviation from the PHP version.
		// If weightings are equal, the amount will be shared.
		{"300.00", []uint{0}, []string{"300.00"}},
		{"300.00", []uint{0, 0, 0}, []string{"100.00", "100.00", "100.00"}},

		// Repeat all of the above with negatives.
		{"-0.00", []uint{1, 1, 1}, []string{"0.00", "0.00", "0.00"}},
		{"-0.01", []uint{1, 1, 1}, []string{"-0.01", "0.00", "0.00"}},
		{"-0.02", []uint{1, 1, 1}, []string{"-0.01", "-0.01", "0.00"}},
		{"-0.03", []uint{1, 1, 1}, []string{"-0.01", "-0.01", "-0.01"}},
		{"-0.04", []uint{1, 1, 1}, []string{"-0.02", "-0.01", "-0.01"}},
		{"-0.05", []uint{1, 1, 1}, []string{"-0.02", "-0.02", "-0.01"}},
		{"-1.00", []uint{0, 1, 0}, []string{"0.00", "-1.00", "0.00"}},
		{"-0.03", []uint{0, 5, 0}, []string{"0.00", "-0.03", "0.00"}},
		{"-3.00", []uint{1, 1, 1}, []string{"-1.00", "-1.00", "-1.00"}},
		{"-1.00", []uint{1, 1, 1}, []string{"-0.34", "-0.33", "-0.33"}},
		{"-0.03", []uint{0, 5, 0}, []string{"0.00", "-0.03", "0.00"}},
		{"-0.03", []uint{0, 4, 2}, []string{"0.00", "-0.02", "-0.01"}},
		{"-1.05", []uint{3, 7}, []string{"-0.32", "-0.73"}},
		{"-0.05", []uint{1, 1}, []string{"-0.03", "-0.02"}},
		{"-300.00", []uint{122, 878}, []string{"-36.60", "-263.40"}},
		{"-300.00", []uint{122, 0, 878}, []string{"-36.60", "0.00", "-263.40"}},
		{"-120.00", []uint{20, 100}, []string{"-20.00", "-100.00"}},
		{"-300.00", []uint{0}, []string{"-300.00"}},
		{"-300.00", []uint{0, 0, 0}, []string{"-100.00", "-100.00", "-100.00"}},
	}
	for ci, c := range cases {
		a, _ := New("GBP", c.a)
		res := a.Share(c.ratios)
		if len(c.ratios) != len(res) {
			t.Errorf("Case %d. Incorrect number of allocations returned. Expected %d, got %d: %v", ci, len(c.ratios), len(res), res)
			return
		}
		for i := range c.want {
			if c.want[i] != res[i].Amount() {
				t.Errorf("Case %d: Sharing %s into (%v), portion %d: Expected %s, got %s", ci, c.a, c.ratios, i, c.want[i], res[i].Amount())
			}
		}
	}
}
