package translate_test

import (
	"reflect"
	"testing"

	"github.com/karlovskiy/vocabacov/internal/translate"
)

func TestFindPhrase(t *testing.T) {
	tests := []struct {
		text        string
		lang        string
		phrase      string
		translation string
	}{
		{"/en\ncruel\nжестокий", "en", "cruel", "жестокий"},
		{"/fr\nlabiba\nлябиба", "fr", "labiba", "лябиба"},
		{"/es\nhola chica\nпривет девушка", "es", "hola chica", "привет девушка"},
	}
	for i, test := range tests {
		actualPhrase, err := translate.FindPhrase(test.text)
		if err != nil {
			t.Fatalf("error finding phrase: %v", err)
		}
		if actualPhrase == nil {
			t.Fatalf("phrase not found")
		}
		if actualPhrase.Lang != test.lang {
			t.Errorf("%d: Got lang: %q, want: %q", i, actualPhrase.Lang, test.lang)
		}
		if actualPhrase.Phrase != test.phrase {
			t.Errorf("%d: Got phrase: %q, want: %q", i, actualPhrase.Phrase, test.phrase)
		}
		if actualPhrase.Translation != test.translation {
			t.Errorf("%d: Got translation: %q, want: %q", i, actualPhrase.Translation, test.translation)
		}
	}
}

func TestFindCommand(t *testing.T) {
	tests := []struct {
		text    string
		command string
		args    []string
	}{
		{"/export en", "export", []string{"en"}},
		{"/reset es", "reset", []string{"es"}},
	}
	for i, test := range tests {
		actualCommand, err := translate.FindCommand(test.text)
		if err != nil {
			t.Fatalf("error finding command: %v", err)
		}
		if actualCommand == nil {
			t.Fatalf("command not found")
		}
		if actualCommand.Name != test.command {
			t.Errorf("%d: Got command: %q, want: %q", i, actualCommand.Name, test.command)
		}
		if !reflect.DeepEqual(actualCommand.Args, test.args) {
			t.Errorf("%d: Got args: %v, want: %v", i, actualCommand.Args, test.args)
		}
	}

}
