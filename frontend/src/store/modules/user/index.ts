import { defineStore } from 'pinia';
import {
  login as userLogin,
  logout as userLogout,
  getUserInfo,
  LoginData,
} from '@/api/user';
import { setToken, clearToken } from '@/utils/auth';
import { removeRouteListener } from '@/utils/route-listener';
import { UserState } from './types';
import useAppStore from '../app';
import {AuthService, UserService} from "@/services/grpc.ts";

const useUserStore = defineStore('hermes', {
  state: (): UserState => ({
    token:  undefined,
    _refreshToken: undefined,
    isLogin: undefined,
    userId: undefined,
    name: undefined,
    avatar: undefined,
    location: undefined,
    email: undefined,
    phone: undefined,
    role: '',
  }),

  getters: {
    userInfo(state: UserState): UserState {
      return { ...state };
    },
  },

  actions: {
    switchRoles() {
      return new Promise((resolve) => {
        this.role = this.role === 'user' ? 'admin' : 'user';
        resolve(this.role);
      });
    },
    // Set user's information
    setInfo(partial: Partial<UserState>) {
      this.$patch(partial);
    },

    // Reset user's information
    resetInfo() {
      // this.$reset();
    },

    // Get user's information
    info() {
      // rbacValues res = await getUserInfo();
      UserService.GetUser({}).then((res) => {
        this.userId = res.info?.id;
        this.name = res.info?.name;
        this.email = res.info?.email;
        this.projectLimit = res.info?.limit;
        if (res.info?.isAdmin) {
          this.role = 'admin';
        }else {
          this.role = 'user';
        }
      })

      this.setInfo({
        // name: 'Admin',
      });
    },

    async refreshToken() {
      try {
        const request = { refreshToken: this._refreshToken };
        const initReq = {
          pathPrefix: import.meta.env.VITE_API_BASE_URL,
        };

        await AuthService.RefreshToken(request, initReq).then((res) => {
          if (!res.accessToken || !res.refreshToken) {
            throw new Error('Invalid response from login');
          }
          this.token = res.accessToken;
          this._refreshToken = res.refreshToken;
        }).catch((err) => {
          console.error(err);
          throw err;
        });
      } catch (err) {
        this.token = "";
        this._refreshToken = "";
        console.error(err);
        throw err;
      }
    },

    // Login
    async login(loginForm: LoginData) {
      try {
        const request = { email: loginForm.email, password: loginForm.password };
        const initReq = {
          pathPrefix: import.meta.env.VITE_GAPI_URL,
        };

        await AuthService.Login(request, initReq).then((res) => {
          if (!res.accessToken || !res.refreshToken) {
            throw new Error('Invalid response from login');
          }
          this.isLogin = true;
          this.token = res.accessToken;
          this._refreshToken = res.refreshToken;
          UserService.GetUser({}).then((res) => {
            if (res.user?.role) {
              this.role = 'admin';
            }else {
              this.role = 'user';
            }
          }).catch((err) => {
            console.error(err);
            throw err;
          });
        }).catch((err) => {
          console.error(err);
          throw err;
        }).finally(() => {
        });
      } catch (err) {
        this.token = "";
        this._refreshToken = "";
        this.isLogin = false;
        throw err;
      }
    },
    logoutCallBack() {
      const appStore = useAppStore();
      this.resetInfo();
      this.token = "";
      this._refreshToken = "";
      this.isLogin = false;
      removeRouteListener();
      appStore.clearServerMenu();
    },
    // Logout
    async logout() {
      try {
        // await userLogout();
      } finally {
        this.logoutCallBack();
      }
    },
  },
  persist: true
});

export default useUserStore;
