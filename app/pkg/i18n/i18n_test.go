package i18n

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// ロケールJSONと model.Localization の宣言がずれると、翻訳は「黙って空文字」になり、
// 画面にはキー名がそのまま出る。ビルドもテストも通ってしまい実機でしか気づけないため、
// ここで機械的に塞ぐ。
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

			// JSON にあるが構造体に宣言が無い → 読み込まれず黙って捨てられる
			for key := range keys {
				if _, ok := declared[key]; !ok {
					t.Errorf("locale key %q is not declared in model.Localization; it would be silently dropped", key)
				}
			}
			// 構造体にあるが JSON に無い → 空文字が画面に出る
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
