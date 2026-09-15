package repull

import "testing"

func TestRenamedEnumAliasesKeepValues(t *testing.T) {
	cases := map[string]string{
		string(FullAccess):     "full_access",
		string(Messaging):      "messaging",
		string(ReadOnly):       "read_only",
		string(Availability):   "availability",
		string(DerivedPricing): "derived-pricing",
		string(Rates):          "rates",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("alias value %q, want %q", got, want)
		}
	}
	var _ CreateConnectionJSONBodyAccessType = FullAccess
	var _ BookingAvailabilityUpdateRequestType = Rates
}
