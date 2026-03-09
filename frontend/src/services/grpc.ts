import { AuthService as originAuthService } from '@/lib/proto/auth/v1/auth.pb'
import { BonusService as originBonusService } from '@/lib/proto/bonus/v1/bonus.pb'
import { CategoryService as originCategoryService } from '@/lib/proto/category/v1/category.pb'
import { CommentService as originCommentService } from '@/lib/proto/comment/v1/comment.pb'
import type { InitReq } from '@/lib/proto/fetch.pb'
import { InviteService as originInviteService } from '@/lib/proto/invite/v1/invite.pb'
import { MessageService as originMessageService } from '@/lib/proto/message/v1/message.pb'
import { ResourceService as originResourceService } from '@/lib/proto/resource/v1/resource.pb'
import { SystemService as originSystemService } from '@/lib/proto/system/v1/system.pb'
import { TorrentService as originTorrentService } from '@/lib/proto/torrent/v1/torrent.pb'
import { TrackerService as originTrackerService } from '@/lib/proto/tracker/v1/tracker.pb'
import { TrafficService as originTrafficService } from '@/lib/proto/traffic/v1/traffic.pb'
import { UserService as originUserService } from '@/lib/proto/user/v1/user.pb'
import pinia, { useAuthStore } from '@/store'

type GrpcInitReq = InitReq & {
  headers?: HeadersInit
}

type GrpcMethod = (...args: unknown[]) => Promise<unknown>

const DEFAULT_PATH_PREFIX = import.meta.env.VITE_API_BASE_URL || ''
const UNAUTHORIZED_CODES = new Set([16, 401])

let unauthorizedRedirecting = false

export const AuthService = createProxy(originAuthService)
export const BonusService = createProxy(originBonusService)
export const CategoryService = createProxy(originCategoryService)
export const CommentService = createProxy(originCommentService)
export const InviteService = createProxy(originInviteService)
export const MessageService = createProxy(originMessageService)
export const ResourceService = createProxy(originResourceService)
export const SystemService = createProxy(originSystemService)
export const TorrentService = createProxy(originTorrentService)
export const TrackerService = createProxy(originTrackerService)
export const TrafficService = createProxy(originTrafficService)
export const UserService = createProxy(originUserService)

function createProxy<T extends object>(service: T): T {
  const handler: ProxyHandler<T> = {
    get(target, prop, receiver) {
      const targetProperty = Reflect.get(target, prop, receiver)
      if (!isGrpcMethod(targetProperty)) {
        return targetProperty
      }

      return async (...args: unknown[]) => {
        const finalArgs = args.length === 0 ? [{}] : [...args]
        const lastArg = finalArgs[finalArgs.length - 1]

        if (finalArgs.length === 1 || !isInitReqLike(lastArg)) {
          finalArgs.push(buildInitReq())
        } else {
          finalArgs[finalArgs.length - 1] = buildInitReq(lastArg)
        }

        try {
          return await targetProperty(...finalArgs)
        } catch (error: unknown) {
          const requestError = new RequestError(error)
          if (UNAUTHORIZED_CODES.has(requestError.code)) {
            handleUnauthorized()
          }
          throw requestError
        }
      }
    },
  }

  return new Proxy(service, handler)
}

function buildInitReq(rawInitReq?: unknown): GrpcInitReq {
  const initReq = isInitReqLike(rawInitReq) ? rawInitReq : {}
  const headers = normalizeHeaders(initReq.headers)
  const authStore = useAuthStore(pinia)

  if (authStore.isLoggedIn) {
    headers.Authorization = `Bearer ${authStore.accessToken}`
  }

  return {
    ...initReq,
    pathPrefix: initReq.pathPrefix || DEFAULT_PATH_PREFIX,
    headers,
  }
}

function isGrpcMethod(value: unknown): value is GrpcMethod {
  return typeof value === 'function'
}

function isInitReqLike(value: unknown): value is GrpcInitReq {
  if (!value || typeof value !== 'object') {
    return false
  }

  return 'pathPrefix' in value || 'headers' in value || 'signal' in value || 'credentials' in value
}

function normalizeHeaders(headers?: HeadersInit): Record<string, string> {
  if (!headers) {
    return {}
  }

  if (headers instanceof Headers) {
    return Object.fromEntries(headers.entries())
  }

  if (Array.isArray(headers)) {
    return headers.reduce<Record<string, string>>((record, [key, value]) => {
      record[key] = value
      return record
    }, {})
  }

  return { ...headers }
}

function handleUnauthorized() {
  const authStore = useAuthStore(pinia)
  authStore.handleUnauthorized()

  if (typeof window === 'undefined' || unauthorizedRedirecting) {
    return
  }

  if (window.location.pathname === '/login') {
    return
  }

  unauthorizedRedirecting = true
  const query = new URLSearchParams()
  query.set('redirect', `${window.location.pathname}${window.location.search}`)
  window.location.href = `/login?${query.toString()}`

  setTimeout(() => {
    unauthorizedRedirecting = false
  }, 300)
}

interface NormalizeErrorResult {
  code: number
  message: string
  details?: unknown[]
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function parseCode(value: unknown): number {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value
  }
  if (typeof value === 'string') {
    const parsed = Number(value)
    if (!Number.isNaN(parsed)) {
      return parsed
    }
  }
  return 0
}

function normalizeError(error: unknown): NormalizeErrorResult {
  if (error instanceof RequestError) {
    return {
      code: error.code,
      message: error.message,
      details: error.details,
    }
  }

  let code = 0
  let message = '非预期错误'
  let details: unknown[] | undefined

  if (error instanceof Error && error.message) {
    message = error.message
  }

  if (isRecord(error)) {
    if ('code' in error) {
      code = parseCode(error.code)
    }
    if ('message' in error && typeof error.message === 'string' && error.message) {
      message = error.message
    } else if ('error' in error && typeof error.error === 'string' && error.error) {
      message = error.error
    }
    if ('details' in error && Array.isArray(error.details)) {
      details = error.details
    }

    if ('error' in error && isRecord(error.error)) {
      const nestedError = error.error
      if (code === 0 && 'code' in nestedError) {
        code = parseCode(nestedError.code)
      }
      if ((message === '非预期错误' || !message) && 'message' in nestedError && typeof nestedError.message === 'string') {
        message = nestedError.message
      }
    }
  }

  if (message === 'Failed to fetch') {
    message = '网络连接失败'
  }

  return { code, message, details }
}

export class RequestError extends Error {
  public code: number
  public details?: unknown[]

  constructor(error: unknown) {
    const normalized = normalizeError(error)
    super(normalized.message)
    this.name = 'RequestError'
    this.code = normalized.code
    this.details = normalized.details
  }
}
