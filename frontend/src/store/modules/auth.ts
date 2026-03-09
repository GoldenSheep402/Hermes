import { defineStore } from 'pinia'
import { AuthService, UserService } from '@/services/grpc'

export interface PtUserProfile {
  id: string
  username: string
  email: string
  className: string
  reputation: number
  uploadBytes: number
  downloadBytes: number
  bonusPoints: number
  inboxUnread: number
  inviteCount: number
  avatar: string
  roles: string[]
}

interface LoginPayload {
  email: string
  password: string
}

function createDefaultProfile(): PtUserProfile {
  return {
    id: '',
    username: 'Guest',
    email: '',
    className: 'User',
    reputation: 0,
    uploadBytes: 0,
    downloadBytes: 0,
    bonusPoints: 0,
    inboxUnread: 0,
    inviteCount: 0,
    avatar: '',
    roles: ['user'],
  }
}

export const useAuthStore = defineStore('hermes/auth', {
  state: () => ({
    accessToken: '',
    refreshToken: '',
    profile: createDefaultProfile(),
  }),

  getters: {
    isLoggedIn: (state): boolean => state.accessToken.length > 0,
    isStaff: (state): boolean => state.profile.roles.includes('staff'),
    ratio: (state): number => {
      if (state.profile.downloadBytes === 0) {
        return state.profile.uploadBytes > 0 ? 99.99 : 0
      }
      return Number((state.profile.uploadBytes / state.profile.downloadBytes).toFixed(2))
    },
  },

  actions: {
    applySession(token: string, profilePatch?: Partial<PtUserProfile>) {
      this.accessToken = token
      this.profile = {
        ...createDefaultProfile(),
        ...this.profile,
        ...profilePatch,
      }
    },

    async login(payload: LoginPayload) {
      const email = payload.email.trim()
      if (!payload.password) {
        throw new Error('密码不能为空')
      }

      if (!email) {
        throw new Error('邮箱不能为空')
      }

      const response = await AuthService.Login({
        email,
        password: payload.password,
      })

      if (!response.accessToken || !response.refreshToken) {
        throw new Error('登录响应无效')
      }

      this.accessToken = response.accessToken
      this.refreshToken = response.refreshToken
      await this.fetchProfile()
    },

    async fetchProfile() {
      const [basicUserResult, userProfileResult] = await Promise.allSettled([
        UserService.GetUser({}),
        UserService.GetUserProfile({}),
      ])

      const basicUser =
        basicUserResult.status === 'fulfilled' ? basicUserResult.value.user : undefined
      const userProfile =
        userProfileResult.status === 'fulfilled' ? userProfileResult.value : undefined
      const user = basicUser || userProfile?.user

      if (!user?.id) {
        throw new Error('无法获取用户信息')
      }

      const upload = parseInt64(userProfile?.realUpload ?? user.uploaded)
      const download = parseInt64(userProfile?.realDownload ?? user.downloaded)
      const isStaff = Boolean(user.isAdmin)

      this.profile = {
        id: user.id,
        username: user.username || 'User',
        email: user.email || '',
        className: isStaff ? 'Sysop' : 'User',
        reputation: 0,
        uploadBytes: upload,
        downloadBytes: download,
        bonusPoints: parseInt64(user.bonusPoints),
        inboxUnread: 0,
        inviteCount: user.inviteCount || 0,
        avatar: user.avatar || '',
        roles: isStaff ? ['user', 'staff'] : ['user'],
      }
    },

    logout() {
      this.$reset()
    },

    handleUnauthorized() {
      this.logout()
    },
  },

  persist: {
    key: 'hermes-auth',
    pick: ['accessToken', 'refreshToken', 'profile'],
  },
})

function parseInt64(value: string | number | undefined): number {
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : 0
  }
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : 0
  }
  return 0
}
