package i18n_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/i18n"

	. "github.com/getfider/fider/app/pkg/assert"
)

var emptyContext = context.Background()
var enContext = context.WithValue(emptyContext, app.LocaleCtxKey, "en")
var ptBRContext = context.WithValue(emptyContext, app.LocaleCtxKey, "pt-BR")

func TestValidMessage(t *testing.T) {
	RegisterT(t)

	key := "error.pagenotfound.title"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("Page not found")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("Page not found")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("Página não encontrada")
}

func TestEnglish_ValidMessage_WithParams(t *testing.T) {
	RegisterT(t)

	key := "email.greetings_name"

	translated := i18n.T(emptyContext, key, i18n.Params{"name": "Jon"})
	Expect(translated).Equals("Hello, Jon!")

	translated = i18n.T(enContext, key, i18n.Params{"name": "Jon"})
	Expect(translated).Equals("Hello, Jon!")

	translated = i18n.T(ptBRContext, key, i18n.Params{"name": "Jon"})
	Expect(translated).Equals("Olá, Jon!")
}

func TestEnglish_InvalidMessage(t *testing.T) {
	RegisterT(t)

	key := "This message does not exist"

	translated := i18n.T(emptyContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")

	translated = i18n.T(enContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")

	translated = i18n.T(ptBRContext, key)
	Expect(translated).Equals("⚠️ Missing Translation: This message does not exist")
}

func TestGetLocale(t *testing.T) {
	RegisterT(t)

	Expect(i18n.GetLocale(emptyContext)).Equals("en")
	Expect(i18n.GetLocale(enContext)).Equals("en")
	Expect(i18n.GetLocale(ptBRContext)).Equals("pt-BR")
}

func TestIsValidLocale(t *testing.T) {
	RegisterT(t)

	//Valid locales
	Expect(i18n.IsValidLocale("en")).IsTrue()
	Expect(i18n.IsValidLocale("pt-BR")).IsTrue()

	//Invalid locales
	Expect(i18n.IsValidLocale("EN")).IsFalse()
	Expect(i18n.IsValidLocale("PT-BR")).IsFalse()
	Expect(i18n.IsValidLocale("pt_BR")).IsFalse()
	Expect(i18n.IsValidLocale("")).IsFalse()
	Expect(i18n.IsValidLocale("xx")).IsFalse()
}
