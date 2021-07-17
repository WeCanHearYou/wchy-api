interface Locale {
  text: string
  translated: number
}

const locales: { [key: string]: Locale } = {
  en: {
    text: "🇺🇸 English",
    translated: 100,
  },
  "pt-BR": {
    text: "🇧🇷 Portuguese (Brazilian)",
    translated: 100,
  },
  fr: {
    text: "🇫🇷 French",
    translated: 100,
  },
}

export default locales
