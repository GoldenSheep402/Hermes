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
  phone?: string;
  projectLimit?: number;
  role: RoleType;
}
