import { locales } from "@/constants"

export type App = {
  link: string
  author: string
  name: string
  imgSrc: string
  thumb: string
  rating: string
}

export type SearchItem = {
  locale: keyof typeof locales
  apps: App[]
}