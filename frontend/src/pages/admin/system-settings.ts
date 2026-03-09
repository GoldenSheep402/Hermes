import type { Setting } from '@/lib/proto/system/v1/system.pb'
import { SystemService } from '@/services/grpc'

export const settingKeys = {
  siteName: 'site.name',
  maintenanceMode: 'site.maintenance_mode',
  marquee: 'site.marquee',
  openRegistration: 'invite.open_registration',
  inviteOnly: 'invite.only',
  emailVerificationRequired: 'invite.email_verification_required',
  smtpEnable: 'auth.smtp_enable',
  globalMessage: 'invite.global_message',
  announceInterval: 'tracker.announce_interval',
  flushInterval: 'tracker.flush_interval',
  flushBatchSize: 'tracker.flush_batch_size',
  globalFreeleech: 'tracker.global_freeleech',
  freeleechCountdown: 'tracker.freeleech_countdown_hours',
  bonusFormula: 'tracker.bonus_formula',
  trackerList: 'tracker.list',
} as const

export function parseBoolean(input: string | undefined, fallback: boolean): boolean {
  if (!input) {
    return fallback
  }
  const normalized = input.trim().toLowerCase()
  if (['1', 'true', 'yes', 'on'].includes(normalized)) {
    return true
  }
  if (['0', 'false', 'no', 'off'].includes(normalized)) {
    return false
  }
  return fallback
}

export function parseNumber(input: string | undefined, fallback: number): number {
  if (!input) {
    return fallback
  }
  const value = Number(input)
  return Number.isFinite(value) ? value : fallback
}

export function boolToValue(input: boolean): string {
  return input ? 'true' : 'false'
}

export function clampNumber(input: number, min: number, max: number, fallback: number): number {
  if (!Number.isFinite(input)) {
    return fallback
  }
  if (input < min) {
    return min
  }
  if (input > max) {
    return max
  }
  return input
}

export function pickChangedSettings(settingsMap: Record<string, string>, items: Setting[]): Setting[] {
  return items.filter((item) => {
    const key = item.key || ''
    if (!key) {
      return false
    }
    if (!Object.prototype.hasOwnProperty.call(settingsMap, key)) {
      return true
    }
    return (item.value || '') !== settingsMap[key]
  })
}

export async function loadSettingsMap(): Promise<Record<string, string>> {
  const response = await SystemService.GetSettings({})
  const map: Record<string, string> = {}
  for (const item of response.settings || []) {
    if (!item.key) {
      continue
    }
    map[item.key] = item.value || ''
  }
  return map
}

export async function saveSettings(items: Setting[]): Promise<void> {
  await SystemService.SetSettings({ settings: items })
}
