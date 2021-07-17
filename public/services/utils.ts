import { Fider } from "."

export const delay = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export const classSet = (input?: any): string => {
  let classes = ""
  if (input) {
    for (const key in input) {
      if (key && !!input[key]) {
        classes += ` ${key}`
      }
    }
    return classes.trim()
  }
  return ""
}

type DateFormat = "full" | "short"
const shortOpts: Intl.DateTimeFormatOptions = { month: "short", year: "2-digit" }
const fullOpts: Intl.DateTimeFormatOptions = { day: "2-digit", month: "long", year: "numeric", hour: "numeric", minute: "numeric" }

export const formatDate = (locale: string, input: Date | string, format: DateFormat = "full"): string => {
  const date = input instanceof Date ? input : new Date(input)

  try {
    const opts = format === "short" ? shortOpts : fullOpts
    return new Intl.DateTimeFormat(locale, opts).format(date)
  } catch {
    return date.toLocaleString(locale)
  }
}

export const timeSince = (locale: string, now: Date, date: Date): string => {
  try {
    const seconds = Math.round((now.getTime() - date.getTime()) / 1000)
    const minutes = Math.round(seconds / 60)
    const hours = Math.round(minutes / 60)
    const days = Math.round(hours / 24)
    const months = Math.round(days / 30)
    const years = Math.round(days / 365)

    const rtf = new Intl.RelativeTimeFormat(locale, { numeric: "auto" })
    return (
      (seconds < 60 && rtf.format(-1 * seconds, "seconds")) ||
      (minutes < 60 && rtf.format(-1 * minutes, "minutes")) ||
      (hours < 24 && rtf.format(-1 * hours, "hours")) ||
      (days < 30 && rtf.format(-1 * days, "days")) ||
      (days < 365 && rtf.format(-1 * months, "months")) ||
      rtf.format(-1 * years, "years")
    )
  } catch {
    return formatDate(locale, date, "short")
  }
}

export const fileToBase64 = async (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.addEventListener(
      "load",
      () => {
        const parts = (reader.result as string).split("base64,")
        resolve(parts[1])
      },
      false
    )

    reader.addEventListener(
      "error",
      () => {
        reject(reader.error)
      },
      false
    )

    reader.readAsDataURL(file)
  })
}

export const timeAgo = (date: string | Date): number => {
  const d = date instanceof Date ? date : new Date(date)
  return (new Date().getTime() - d.getTime()) / 1000
}

export const isCookieEnabled = (): boolean => {
  try {
    document.cookie = "cookietest=1"
    const ret = document.cookie.indexOf("cookietest=") !== -1
    document.cookie = "cookietest=1; expires=Thu, 01-Jan-1970 00:00:01 GMT"
    return ret
  } catch (e) {
    return false
  }
}

export const uploadedImageURL = (bkey: string | undefined, size?: number): string | undefined => {
  if (bkey) {
    if (size) {
      return `${Fider.settings.assetsURL}/static/images/${bkey}?size=${size}`
    }
    return `${Fider.settings.assetsURL}/static/images/${bkey}`
  }
  return undefined
}

export const truncate = (input: string, maxLength: number): string => {
  if (input && input.length > maxLength) {
    return `${input.substr(0, maxLength)}...`
  }
  return input
}
