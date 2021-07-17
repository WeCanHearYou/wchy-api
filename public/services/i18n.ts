import { i18n, I18n } from "@lingui/core"
import { en, pt, fr } from "make-plural/plurals"

export function activateI18NSync(locale: string, messages?: any): I18n {
  i18n._missing = (_, key) => `⚠️ Missing Translation: ${key}`
  i18n.loadLocaleData("en", { plurals: en })
  i18n.loadLocaleData("pt-BR", { plurals: pt })
  i18n.loadLocaleData("fr", { plurals: fr })
  i18n.load(locale, messages)
  i18n.activate(locale)
  return i18n
}

export async function activateI18N(locale: string): Promise<I18n> {
  try {
    const content = await import(
      /* webpackChunkName: "locale-[request]" */
      `@locale/${locale}/client.json`
    )
    return activateI18NSync(locale, content.messages)
  } catch (err) {
    console.error(err)
    return activateI18NSync(locale)
  }
}
