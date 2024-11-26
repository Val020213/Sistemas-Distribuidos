import { black, tailwindColors, white } from './tailwindColors'

export const common = {
  black: black,
  white: white,
}

export const primary = {
  main: tailwindColors.green[500],
  light: tailwindColors.green[400],
  contrastText: common.black,
}

export const secondary = {
  main: tailwindColors.orange[500],
  light: tailwindColors.orange[400],
  contrastText: common.black,
}

export const error = {
  main: tailwindColors.red[500],
  light: tailwindColors.red[400],
  contrastText: common.black,
}

export const info = {
  main: tailwindColors.blue[500],
  light: tailwindColors.blue[400],
  contrastText: common.black,
}

export const success = {
  main: tailwindColors.green[500],
  light: tailwindColors.green[400],
  contrastText: common.black,
}

export const warning = {
  main: tailwindColors.yellow[500],
  light: tailwindColors.yellow[400],
  contrastText: common.black,
}

export const background = {
  grey: '#EEEEEE',
  light: '#4ade80',
  black: '#252525',
}
