package i18n

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// If the locale JSON and model.Localization drift apart, translations silently become
// empty and the raw key shows in the UI. Nothing else catches it, so check it here.
func TestLocalesMatchLocalizationStruct(t *testing.T) {
	languages, err := GetSupportedLanguages()
	if err != nil {
		t.Fatalf("get supported languages: %v", err)
	}
	if len(languages) == 0 {
		t.Fatal("no locales found")
	}

	localizationType := reflect.TypeOf(model.Localization{})
	declared := make(map[string]struct{}, localizationType.NumField())
	for i := range localizationType.NumField() {
		tag := localizationType.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		declared[tag] = struct{}{}
	}

	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			raw, err := localeFs.ReadFile(fmt.Sprintf("locales/%s.json", lang))
			if err != nil {
				t.Fatalf("read locale json: %v", err)
			}
			var keys map[string]any
			if err := json.Unmarshal(raw, &keys); err != nil {
				t.Fatalf("unmarshal locale json: %v", err)
			}

			// In JSON but not in the struct: silently dropped
			for key := range keys {
				if _, ok := declared[key]; !ok {
					t.Errorf("locale key %q is not declared in model.Localization; it would be silently dropped", key)
				}
			}
			// In the struct but not in JSON: shows as empty
			for key := range declared {
				if _, ok := keys[key]; !ok {
					t.Errorf("model.Localization declares %q but %s.json does not define it", key, lang)
				}
			}

			lng, err := GetTranslation(lang)
			if err != nil {
				t.Fatalf("get translation: %v", err)
			}
			value := reflect.ValueOf(*lng)
			for i := range localizationType.NumField() {
				if value.Field(i).Kind() != reflect.String {
					continue
				}
				if value.Field(i).String() == "" {
					t.Errorf("translation for %q is empty in %s.json", localizationType.Field(i).Tag.Get("json"), lang)
				}
			}
		})
	}
}
