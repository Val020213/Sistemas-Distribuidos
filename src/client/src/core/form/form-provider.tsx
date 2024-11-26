import {
  FieldValues,
  FormProviderProps,
  FormProvider as RHFProvider,
  SubmitHandler,
} from 'react-hook-form'

const FormProvider = <T extends FieldValues>({
  props,
  onSubmit,
  autocomplete = 'new-password',
  children,
}: {
  props: Omit<FormProviderProps<T>, 'children'>
  autocomplete?: string
  onSubmit: SubmitHandler<T>
  children: React.ReactNode
}) => {
  return (
    <RHFProvider {...props}>
      <form onSubmit={props.handleSubmit(onSubmit)} autoComplete={autocomplete}>
        {children}
      </form>
    </RHFProvider>
  )
}

export { FormProvider }
