import type {
  LoginData,
} from '@/api/user'

import type { UserState } from './types'
import { AuthService, UserService } from '@/services/grpc.ts'
import { removeRouteListener } from '@/utils/route-listener'
import { defineStore } from 'pinia'
import useAppStore from '../app'

const useUserStore = defineStore('hermes', {
  state: (): UserState => ({
    token: undefined,
    _refreshToken: undefined,
    isLogin: undefined,
    userId: undefined,
    name: undefined,
    avatar: undefined,
    location: undefined,
    email: undefined,
    role: '',
  }),

  getters: {
    userInfo(state: UserState): UserState {
      return { ...state }
    },
  },

  actions: {
    switchRoles() {
      return new Promise((resolve) => {
        this.role = this.role === 'user' ? 'admin' : 'user'
        resolve(this.role)
      })
    },
    // Set user's information
    setInfo(partial: Partial<UserState>) {
      this.$patch(partial)
    },

    // Reset user's information
    resetInfo() {
      this.$reset()
    },

    // Get user's information
    info() {
      UserService.GetUser({}).then((res) => {
        this.userId = res.user?.id
        this.name = res.user?.name
        this.email = res.user?.role
        // this.projectLimit = res.info?.limit;
        if (res.info?.isAdmin) {
          this.role = 'admin'
        }
        else {
          this.role = 'user'
        }
      })

      this.setInfo({
        avatar: 'https://www.z4a.net/images/2025/04/09/Snipaste_2025-04-09_01-17-44.png',
      })
    },

    async refreshToken() {
      try {
        const request = { refreshToken: this._refreshToken }
        const initReq = {
          pathPrefix: import.meta.env.VITE_API_BASE_URL,
        }

        await AuthService.RefreshToken(request, initReq).then((res) => {
          if (!res.accessToken || !res.refreshToken) {
            throw new Error('Invalid response from login')
          }
          this.token = res.accessToken
          this._refreshToken = res.refreshToken
        }).catch((err) => {
          console.error(err)
          throw err
        })
      }
      catch (err) {
        this.token = ''
        this._refreshToken = ''
        console.error(err)
        throw err
      }
    },

    // Login
    async login(loginForm: LoginData) {
      try {
        const request = { email: loginForm.email, password: loginForm.password }
        const initReq = {
          pathPrefix: import.meta.env.VITE_GAPI_URL,
        }

        await AuthService.Login(request, initReq).then(async (res) => {
          if (!res.accessToken || !res.refreshToken) {
            throw new Error('Invalid response from login')
          }
          this.isLogin = true
          this.token = res.accessToken
          this._refreshToken = res.refreshToken
          await UserService.GetUser({}).then((res) => {
            if (res.user?.role === 'admin') {
              this.role = 'admin'
            }
            else {
              this.role = 'user'
            }
          }).catch((err) => {
            console.error(err)
            throw err
          })
        }).catch((err) => {
          console.error(err)
          throw err
        }).finally(() => {
        })
      }
      catch (err) {
        this.token = ''
        this._refreshToken = ''
        this.isLogin = false
        throw err
      }
    },
    logoutCallBack() {
      const appStore = useAppStore()
      this.resetInfo()
      this.token = ''
      this._refreshToken = ''
      this.isLogin = false
      removeRouteListener()
      appStore.clearServerMenu()
    },
    // Logout
    async logout() {
      try {
        // await userLogout();
      }
      finally {
        this.logoutCallBack()
      }
    },
  },
  persist: true,
})

export default useUserStore
