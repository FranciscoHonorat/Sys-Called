export const priorityOptions = [
  { value: 'Low', label: 'Baixa' },
  { value: 'Medium', label: 'Média' },
  { value: 'High', label: 'Alta' },
] as const

export function priorityLabel(value: string | undefined): string {
  return priorityOptions.find((option) => option.value === value)?.label ?? '—'
}
