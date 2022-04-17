package vocabacov

import (
	"reflect"
	"testing"
)

func TestParseChannelsSuccess(t *testing.T) {
	tests := []struct {
		value    string
		channels map[int64]struct{}
	}{
		{"123", map[int64]struct{}{123: {}}},
		{"1,2,3", map[int64]struct{}{1: {}, 2: {}, 3: {}}},
		{"1, 2, 3", map[int64]struct{}{1: {}, 2: {}, 3: {}}},
		{"1 , 2 , 3 ", map[int64]struct{}{1: {}, 2: {}, 3: {}}},
		{"123, -1234, 12345", map[int64]struct{}{123: {}, -1234: {}, 12345: {}}},
	}
	for i, test := range tests {
		channels, err := parseChannels(test.value)
		if err != nil {
			t.Errorf("%d: Got error: %v, want: %v", i, err, test.channels)
		}
		if !reflect.DeepEqual(channels, test.channels) {
			t.Errorf("%d: Got channels: %v, want: %v", i, channels, test.channels)
		}
	}
}

func TestParseChannelsError(t *testing.T) {
	tests := []struct {
		value         string
		expectedError string
	}{
		{"", "channels not found"},
		{"test", "channel \"test\" int parsing error strconv.ParseInt: parsing \"test\": invalid syntax"},
		{"1, test", "channel \"test\" int parsing error strconv.ParseInt: parsing \"test\": invalid syntax"},
	}
	for i, test := range tests {
		channels, err := parseChannels(test.value)
		if channels != nil {
			t.Errorf("%d: Got channels: %v, want: nil", i, channels)
		}
		if err == nil {
			t.Errorf("%d: Got error: nil, want: %q", i, test.expectedError)
		}
		if err.Error() != test.expectedError {
			t.Errorf("%d: Got error: %q, want: %q", i, err, test.expectedError)
		}
	}
}

func TestFindPhrase(t *testing.T) {
	tests := []struct {
		text   string
		lang   string
		phrase string
	}{
		{text: "/en cruel", lang: "en", phrase: "cruel"},
		{text: "/fr labiba", lang: "fr", phrase: "labiba"},
		{text: "/es ola chica", lang: "es", phrase: "ola chica"},
		{text: "///", lang: "", phrase: ""},
		{text: "//en cruel", lang: "en", phrase: "cruel"},
		{text: "////////////es ola chica", lang: "es", phrase: "ola chica"},
	}
	for i, test := range tests {
		actualLang, actualPhrase := findPhrase(test.text)
		if actualLang != test.lang {
			t.Errorf("%d: Got lang: %q, want: %q", i, actualLang, test.lang)
		}
		if actualPhrase != test.phrase {
			t.Errorf("%d: Got phrase: %q, want: %q", i, actualPhrase, test.phrase)
		}
	}
}
