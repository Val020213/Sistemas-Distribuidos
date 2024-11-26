import { CSSProperties, ChangeEvent, HTMLInputTypeAttribute } from 'react'
import { FieldError } from 'react-hook-form'
import { Box, TextField, TextFieldProps, Typography } from '@mui/material'
import { tailwindColors } from '@/theme/tailwindColors'

type Props = TextFieldProps & {
  id: string
  placeholder?: string
  value: string
  onChange?: (_: ChangeEvent<HTMLInputElement>) => void
  onBlur?: () => void
  errorAlert?: FieldError
  label?: string
  underLabel?: string
  disabled?: boolean
  type?: HTMLInputTypeAttribute
  size?: 'small' | 'medium'
  width?: CSSProperties['width']
  required?: boolean
  dataTest?: string
  autoComplete?: string
}

const InputWithLabel = ({
  id,
  placeholder,
  value,
  onChange,
  onBlur,
  errorAlert,
  label,
  underLabel,
  disabled,
  type,
  required,
  size = 'small',
  width = '100%',
  dataTest,
  autoComplete,
  ...rest
}: Props) => (
  <Box sx={{ width, display: 'flex', flexDirection: 'column', gap: 0.5 }}>
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
        {label && (
          <Typography
            fontWeight={500}
            fontSize={'0.625rem'}
            noWrap
            color={
              Boolean(errorAlert)
                ? tailwindColors.red[500]
                : tailwindColors.green[500]
            }
          >
            {label} {required && '*'}
          </Typography>
        )}
        {underLabel && (
          <Typography fontWeight={400} fontSize={10}>
            {underLabel}
          </Typography>
        )}
      </Box>
      <TextField
        autoComplete={autoComplete}
        type={type ?? 'text'}
        size={size}
        id={id}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        error={Boolean(errorAlert)}
        fullWidth
        disabled={disabled}
        data-test={dataTest}
        sx={{
          '& input': {
            color: Boolean(errorAlert)
              ? tailwindColors.red[500]
              : tailwindColors.green[500],
          },
          background: tailwindColors.gray[900],
          '& .MuiOutlinedInput-notchedOutline': {
            borderColor: Boolean(errorAlert)
              ? tailwindColors.red[500]
              : tailwindColors.green[800],
          },
          '& .css-1t8l2tu-MuiInputBase-input-MuiOutlinedInput-input': {
            textAlign: 'left',
          },
        }}
        {...rest}
      />
    </Box>

    {errorAlert && (
      <Typography
        fontSize="0.825em"
        sx={{ marginLeft: 1 }}
        color={tailwindColors.red[500]}
      >
        {errorAlert.message}
      </Typography>
    )}
  </Box>
)

export default InputWithLabel
