import { Controller, useFormContext } from 'react-hook-form'
import { TextFieldProps } from '@mui/material/TextField'
import { CSSProperties, HTMLInputTypeAttribute } from 'react'
import InputWithLabel from '../InputWithLabel'

type Props = TextFieldProps & {
  name: string
  type?: HTMLInputTypeAttribute
  disabled?: boolean
  size?: 'small' | 'medium'
  label?: string
  underLabel?: string
  width?: CSSProperties['width']
  required?: boolean
  dataTest?: string
  isNotAccountant?: boolean
}

export default function RHFInputWithLabel({
  name,
  type,
  disabled,
  size,
  label,
  underLabel,
  width,
  required,
  dataTest,
  autoComplete,
  ...rest
}: Props) {
  const { control } = useFormContext()

  return (
    <Controller
      name={name}
      control={control}
      render={({
        field: { onBlur, onChange, value },
        fieldState: { error },
      }) => (
        <InputWithLabel
          errorAlert={error}
          {...rest}
          required={required}
          underLabel={underLabel}
          label={label}
          onBlur={onBlur}
          onChange={onChange}
          value={value}
          id={name}
          disabled={disabled}
          type={type ?? 'text'}
          size={size}
          width={width}
          dataTest={dataTest}
          autoComplete={autoComplete}
        />
      )}
    />
  )
}
