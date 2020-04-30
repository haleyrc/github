package github

import "testing"

func TestParseLinkHeader(t *testing.T) {
	hdr := `<https://api.github.com/repositories/172137510/issues?page=2>; rel="next", <https://api.github.com/repositories/172137510/issues?page=40>; rel="last"`
	want := map[string]string{
		"next": "/repositories/172137510/issues?page=2",
		"last": "/repositories/172137510/issues?page=40",
	}
	got := parseLinkHeader(hdr)
	if !equalMap(want, got) {
		t.Errorf("parseLinkHeader(%s) = %v, wanted %v", hdr, got, want)
	}
}

func equalMap(first, second map[string]string) bool {
	if len(first) != len(second) {
		return false
	}
	for key, firstVal := range first {
		secondVal, ok := second[key]
		if !ok {
			return false
		}
		if secondVal != firstVal {
			return false
		}
	}
	return true
}
