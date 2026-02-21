/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type User = {
  id?: string
  username?: string
  email?: string
  avatar?: string
  isAdmin?: boolean
  isEnabled?: boolean
  groupId?: string
  bonusPoints?: string
  uploaded?: string
  downloaded?: string
  seedTime?: string
  inviteCount?: number
  passkey?: string
  lastLogin?: string
  createdAt?: string
}

export type UserGroup = {
  id?: string
  name?: string
  description?: string
  level?: number
  minUpload?: string
  minRatio?: number
  minSeedTime?: string
  maxDownloads?: number
  canUpload?: boolean
  canInvite?: boolean
  isImmuneToRatio?: boolean
  color?: string
  icon?: string
}

export type GetUserRequest = {
}

export type GetUserResponse = {
  user?: User
}

export type GetUserProfileRequest = {
  id?: string
}

export type GetUserProfileResponse = {
  user?: User
  realUpload?: string
  realDownload?: string
  ratio?: number
  publishedCount?: number
  seedingCount?: number
  downloadCount?: number
}

export type UpdateUserRequest = {
  username?: string
  avatar?: string
}

export type UpdateUserResponse = {
  user?: User
}

export type UpdatePasswordRequest = {
  oldPassword?: string
  newPassword?: string
}

export type UpdatePasswordResponse = {
}

export type ResetPasskeyRequest = {
}

export type ResetPasskeyResponse = {
  passkey?: string
}

export type GetUserPasskeyRequest = {
}

export type GetUserPasskeyResponse = {
  passkey?: string
}

export type ListUsersRequest = {
  page?: number
  pageSize?: number
}

export type ListUsersResponse = {
  users?: User[]
  total?: string
}

export type CreateUserGroupRequest = {
  group?: UserGroup
}

export type CreateUserGroupResponse = {
}

export type GetUserGroupRequest = {
  id?: string
}

export type GetUserGroupResponse = {
  group?: UserGroup
}

export type ListUserGroupsRequest = {
}

export type ListUserGroupsResponse = {
  groups?: UserGroup[]
}

export type UpdateUserGroupRequest = {
  group?: UserGroup
}

export type UpdateUserGroupResponse = {
}

export type DeleteUserGroupRequest = {
  id?: string
}

export type DeleteUserGroupResponse = {
}

export class UserService {
  static GetUser(req: GetUserRequest, initReq?: fm.InitReq): Promise<GetUserResponse> {
    return fm.fetchReq<GetUserRequest, GetUserResponse>(`/gapi/user/v1/info`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetUserProfile(req: GetUserProfileRequest, initReq?: fm.InitReq): Promise<GetUserProfileResponse> {
    return fm.fetchReq<GetUserProfileRequest, GetUserProfileResponse>(`/gapi/user/v1/profile`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateUser(req: UpdateUserRequest, initReq?: fm.InitReq): Promise<UpdateUserResponse> {
    return fm.fetchReq<UpdateUserRequest, UpdateUserResponse>(`/gapi/user/v1/info`, {...initReq, method: "PUT", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdatePassword(req: UpdatePasswordRequest, initReq?: fm.InitReq): Promise<UpdatePasswordResponse> {
    return fm.fetchReq<UpdatePasswordRequest, UpdatePasswordResponse>(`/gapi/user/v1/password`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ResetPasskey(req: ResetPasskeyRequest, initReq?: fm.InitReq): Promise<ResetPasskeyResponse> {
    return fm.fetchReq<ResetPasskeyRequest, ResetPasskeyResponse>(`/gapi/user/v1/passkey/reset`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetUserPasskey(req: GetUserPasskeyRequest, initReq?: fm.InitReq): Promise<GetUserPasskeyResponse> {
    return fm.fetchReq<GetUserPasskeyRequest, GetUserPasskeyResponse>(`/gapi/user/v1/passkey`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListUsers(req: ListUsersRequest, initReq?: fm.InitReq): Promise<ListUsersResponse> {
    return fm.fetchReq<ListUsersRequest, ListUsersResponse>(`/gapi/user/v1/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateUserGroup(req: CreateUserGroupRequest, initReq?: fm.InitReq): Promise<CreateUserGroupResponse> {
    return fm.fetchReq<CreateUserGroupRequest, CreateUserGroupResponse>(`/gapi/user/v1/group/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetUserGroup(req: GetUserGroupRequest, initReq?: fm.InitReq): Promise<GetUserGroupResponse> {
    return fm.fetchReq<GetUserGroupRequest, GetUserGroupResponse>(`/gapi/user/v1/group/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListUserGroups(req: ListUserGroupsRequest, initReq?: fm.InitReq): Promise<ListUserGroupsResponse> {
    return fm.fetchReq<ListUserGroupsRequest, ListUserGroupsResponse>(`/gapi/user/v1/group/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateUserGroup(req: UpdateUserGroupRequest, initReq?: fm.InitReq): Promise<UpdateUserGroupResponse> {
    return fm.fetchReq<UpdateUserGroupRequest, UpdateUserGroupResponse>(`/gapi/user/v1/group/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteUserGroup(req: DeleteUserGroupRequest, initReq?: fm.InitReq): Promise<DeleteUserGroupResponse> {
    return fm.fetchReq<DeleteUserGroupRequest, DeleteUserGroupResponse>(`/gapi/user/v1/group/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}