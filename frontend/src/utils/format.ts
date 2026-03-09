export function formatNumber(value: number): string {
  return value.toLocaleString('en-US')
}

export function formatRatio(value: number): string {
  if (!Number.isFinite(value)) {
    return '0.00'
  }
  if (value >= 99.99) {
    return '99.99+'
  }
  return value.toFixed(2)
}

export function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) {
    return '0 B'
  }

  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const power = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  const value = bytes / 1024 ** power
  return `${value.toFixed(power >= 3 ? 2 : 1)} ${units[power]}`
}

export function formatDateTime(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
