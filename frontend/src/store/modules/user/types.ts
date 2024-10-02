export type RoleType = '' | '*' | 'admin' | 'user';
export interface UserState {
  token?: string;
  _refreshToken?:string;
  isLogin?: boolean;
  userId?: string;
  name?: string;
  avatar?: string;
  location?: string;
  email?: string;
  role: RoleType;
}
