import pinia, { useAuthStore } from '@/store'

const DEFAULT_PATH_PREFIX = import.meta.env.VITE_API_BASE_URL || ''

function buildDownloadUrl(torrentId: string): string {
  const path = `/api/torrent/download/${encodeURIComponent(torrentId)}`
  return DEFAULT_PATH_PREFIX ? `${DEFAULT_PATH_PREFIX}${path}` : path
}

function sanitizeFileName(name: string): string {
  return name
    .replace(/[\\/:*?"<>|]+/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

function buildTorrentFileName(name: string, fallback: string): string {
  const normalized = sanitizeFileName(name || fallback || 'download')
  const baseName = normalized || 'download'
  return baseName.toLowerCase().endsWith('.torrent') ? baseName : `${baseName}.torrent`
}

function triggerBlobDownload(blob: Blob, fileName: string) {
  const objectUrl = window.URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = objectUrl
  anchor.download = fileName
  anchor.rel = 'noopener'
  document.body.appendChild(anchor)
  anchor.click()
  document.body.removeChild(anchor)
  window.URL.revokeObjectURL(objectUrl)
}

function parseDownloadFileName(contentDisposition: string | null, fallbackName: string): string {
  const header = contentDisposition || ''
  const utf8Match = header.match(/filename\*\s*=\s*UTF-8''([^;]+)/i)
  if (utf8Match && utf8Match[1]) {
    try {
      return buildTorrentFileName(decodeURIComponent(utf8Match[1]), fallbackName)
    } catch {
      return buildTorrentFileName(utf8Match[1], fallbackName)
    }
  }

  const plainMatch = header.match(/filename\s*=\s*\"?([^\";]+)\"?/i)
  if (plainMatch && plainMatch[1]) {
    return buildTorrentFileName(plainMatch[1], fallbackName)
  }

  return buildTorrentFileName('', fallbackName)
}

async function readDownloadErrorMessage(response: Response): Promise<string> {
  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    try {
      const payload = (await response.json()) as { message?: string }
      if (payload.message) {
        return payload.message
      }
    } catch {
      return `下载失败 (${response.status})`
    }
  } else {
    const text = (await response.text()).trim()
    if (text) {
      return text
    }
  }
  return `下载失败 (${response.status})`
}

export async function downloadTorrentFile(torrentId: string, displayName: string = ''): Promise<void> {
  const torrentIdText = torrentId.trim()
  if (!torrentIdText) {
    throw new Error('无效的种子 ID')
  }

  const authStore = useAuthStore(pinia)
  const headers = new Headers()
  if (authStore.accessToken) {
    headers.set('Authorization', `Bearer ${authStore.accessToken}`)
  }

  const response = await fetch(buildDownloadUrl(torrentIdText), {
    method: 'GET',
    headers,
    credentials: 'include',
  })
  if (!response.ok) {
    const message = await readDownloadErrorMessage(response)
    throw new Error(message)
  }

  const blob = await response.blob()
  if (blob.size === 0) {
    throw new Error('种子文件为空')
  }

  const fileName = parseDownloadFileName(
    response.headers.get('content-disposition'),
    displayName || torrentIdText,
  )
  triggerBlobDownload(blob, fileName)
}
