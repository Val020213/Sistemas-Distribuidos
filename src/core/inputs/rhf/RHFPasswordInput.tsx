import { TextFieldProps } from '@mui/material/TextField'
import { useState } from 'react'
import { IconButton, InputAdornment } from '@mui/material'
import RHFInputWithLabel from './RHFInputWithLabel'
import { gray, green } from '@/theme/tailwindColors'
import { Eye, EyeOff } from 'lucide-react'

type Props = Omit<TextFieldProps, 'type'> & {
  name: string
  size?: 'small' | 'medium'
  label?: string
  required?: boolean
}

export default function RHFPasswordInput({
  name,
  size,
  label,
  required,
  ...rest
}: Props) {
  const [viewPassword, setViewPassword] = useState<boolean>(false)
  return (
    <RHFInputWithLabel
      label={label}
      name={name}
      placeholder="••••••••••"
      type={viewPassword ? 'text' : 'password'}
      size={size}
      required={required}
      {...rest}
      InputProps={{
        endAdornment: (
          <InputAdornment position="end">
            <IconButton
              onClick={() => setViewPassword(!viewPassword)}
              edge="end"
              disableRipple
              sx={{
                color: gray[600],
                '&:hover': { color: green[500] },
              }}
            >
              {viewPassword ? <Eye /> : <EyeOff />}
            </IconButton>
          </InputAdornment>
        ),
      }}
    />
  )
}
