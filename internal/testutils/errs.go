package testutils

import (
	"github.com/etc-sudonters/zootler/errs"
)

func Expected(frag errs.Fragment) errs.Fragment {
	return errs.WriteAfter("\n", errs.WriteBefore("expected\t\t", frag))
}

func Actual(frag errs.Fragment) errs.Fragment {
	return errs.WriteAfter("\n", errs.WriteBefore("received\t\t", frag))
}
