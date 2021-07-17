import React from "react"
import { renderToStaticMarkup } from "react-dom/server"
import { Fider, FiderContext } from "./services/fider"
import { Header } from "./components"
import { resolveRootComponent, route } from "./router"

import HomePage from "./pages/Home/Home.page"
import ShowPostPage from "./pages/ShowPost/ShowPost.page"
import SignInPage from "./pages/SignIn/SignIn.page"
import SignUpPage from "./pages/SignUp/SignUp.page"
import DesignSystemPage from "./pages/DesignSystem/DesignSystem.page"
import LegalPage from "./pages/Legal/Legal.page"

import { activateI18NSync } from "./services"
import { I18nProvider } from "@lingui/react"

// Locale files must be bundled for SSR to work synchronously
const messages: { [key: string]: any } = {
  en: require(`../locale/en/client`),
  "pt-BR": require(`../locale/pt-BR/client`),
  fr: require(`../locale/fr/client`),
}

// Only public routes should be here
// Routes behind authentication are not crawled
const routes = [
  route("", HomePage),
  route("/posts/:number*", ShowPostPage),
  route("/signin", SignInPage, false),
  route("/signup", SignUpPage, false),
  route("/terms", LegalPage, false),
  route("/privacy", LegalPage, false),
  route("/_design", DesignSystemPage),
]

function ssrRender(url: string, pathname: string, args: any) {
  const fider = Fider.initialize({ ...args })
  const i18n = activateI18NSync(fider.currentLocale, messages[fider.currentLocale].messages)
  const config = resolveRootComponent(pathname, routes)
  window.location.href = url

  return renderToStaticMarkup(
    <I18nProvider i18n={i18n}>
      <FiderContext.Provider value={fider}>
        {config.showHeader && <Header />}
        {React.createElement(config.component, args.props)}
      </FiderContext.Provider>
    </I18nProvider>
  )
}

;(globalThis as any).ssrRender = ssrRender
