'use client'
import { useMemo } from 'react'
import CssBaseline from '@mui/material/CssBaseline'
import {
  ThemeProvider as MuiThemeProvider,
  ThemeOptions,
  createTheme,
} from '@mui/material/styles'
import { typography } from './typography'
import {
  background,
  error,
  primary,
  secondary,
  success,
  warning,
} from './colors'
import shadows from '@mui/material/styles/shadows'
import { shadows as customizedShadows } from './shadows'
import {
  black,
  white,
  slate,
  gray,
  zinc,
  neutral,
  stone,
  red,
  orange,
  amber,
  yellow,
  lime,
  green,
  emerald,
  teal,
  cyan,
  sky,
  blue,
  indigo,
  violet,
  purple,
  fuchsia,
  pink,
  rose,
} from './tailwindColors'
type Props = {
  children: React.ReactNode
}

export default function ThemeProvider({ children }: Readonly<Props>) {
  const memoizedValue = useMemo(
    () => ({
      palette: {
        primary: { main: primary.main, contrastText: primary.contrastText },
        secondary: {
          main: secondary.main,
          contrastText: secondary.contrastText,
        },
        success: { main: success.main, contrastText: success.contrastText },
        error: { main: error.main, contrastText: error.contrastText },
        warning: { main: warning.main, contrastText: warning.contrastText },
        text: { primary: green[500], main: green[400] },
        black: black,
        white: white,
        slate: slate,
        gray: gray,
        zinc: zinc,
        neutral: neutral,
        stone: stone,
        red: red,
        orange: orange,
        amber: amber,
        yellow: yellow,
        lime: lime,
        green: green,
        emerald: emerald,
        teal: teal,
        cyan: cyan,
        sky: sky,
        blue: blue,
        indigo: indigo,
        violet: violet,
        purple: purple,
        fuchsia: fuchsia,
        pink: pink,
        rose: rose,
      },
      background: background,
      typography: typography,
      shape: { borderRadius: 8 },
      shadows: { ...shadows, ...customizedShadows() },
      spacing: (factor: number) => `${factor * 8}px`,
    }),
    []
  )

  const theme = createTheme(memoizedValue as ThemeOptions)

  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      {children}
    </MuiThemeProvider>
  )
}
