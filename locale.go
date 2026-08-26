package verax

import "github.com/verax-validation/validation/internal/messages"

// Locale constants for commonly used languages.
// Language codes follow BCP 47; only a high-frequency subset is listed here to avoid typos;
// for languages not listed, pass the string directly, e.g. RegisterLocale("pt-BR", messages).
const (
	LocaleEn   = "en"    // English
	LocaleZhCN = "zh-CN" // Simplified Chinese
	LocaleZhTW = "zh-TW" // Traditional Chinese
	LocaleJa   = "ja"    // Japanese
	LocaleKo   = "ko"    // Korean
	LocaleFr   = "fr"    // French
	LocaleDe   = "de"    // German
	LocaleEs   = "es"    // Spanish
	LocalePtBR = "pt-BR" // Brazilian Portuguese
	LocaleRu   = "ru"    // Russian
)

// RegisterZhCN registers the Simplified Chinese translation table; idempotent, later registration overwrites the same code.
func RegisterZhCN() {
	RegisterLocale(LocaleZhCN, messages.ZhCN)
}

// RegisterZhTW registers the Traditional Chinese translation table.
func RegisterZhTW() {
	RegisterLocale(LocaleZhTW, messages.ZhTW)
}

// RegisterJa registers the Japanese translation table.
func RegisterJa() {
	RegisterLocale(LocaleJa, messages.Ja)
}

// RegisterFr registers the French translation table.
func RegisterFr() {
	RegisterLocale(LocaleFr, messages.Fr)
}

// RegisterDe registers the German translation table.
func RegisterDe() {
	RegisterLocale(LocaleDe, messages.De)
}

// RegisterEn registers the English translation table.
// English is also the default when no language is registered; this table stays consistent with it entry by entry.
func RegisterEn() {
	RegisterLocale(LocaleEn, messages.En)
}
