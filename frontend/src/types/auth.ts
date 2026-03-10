export interface PtUserProfile {
  id: string
  username: string
  email: string
  className: string
  reputation: number
  uploadBytes: number
  downloadBytes: number
  uploadRateBytes: number
  downloadRateBytes: number
  bonusPoints: number
  inboxUnread: number
  inviteCount: number
  avatar: string
  roles: string[]
}

export interface AuthLoginPayload {
  email: string
  password: string
}
