import dayjs from 'dayjs'

export type DomainTimePreset = '24h' | '7d' | '30d' | 'live'

export interface DomainTimeRange {
  start: number
  end: number
}

export function resolveDomainTimeRange(preset: DomainTimePreset): DomainTimeRange {
  const end = dayjs()
  switch (preset) {
    case '7d':
      return { start: end.subtract(7, 'day').unix(), end: end.unix() }
    case '30d':
      return { start: end.subtract(30, 'day').unix(), end: end.unix() }
    case 'live':
      return { start: end.subtract(15, 'minute').unix(), end: end.unix() }
    case '24h':
    default:
      return { start: end.subtract(24, 'hour').unix(), end: end.unix() }
  }
}

export function timePresetLabel(preset: DomainTimePreset): string {
  switch (preset) {
    case '7d':
      return $gettext('7 Days')
    case '30d':
      return $gettext('30 Days')
    case 'live':
      return $gettext('Live')
    case '24h':
    default:
      return $gettext('24 Hours')
  }
}
