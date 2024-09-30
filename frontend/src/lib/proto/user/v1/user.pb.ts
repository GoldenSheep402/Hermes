/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type User = {
  id?: string
  name?: string
  role?: string
}

export type GetUserRequest = {
  id?: string
}

export type GetUserResponse = {
  user?: User
}

export type UpdateUserRequest = {
  user?: User
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

export type Group = {
  id?: string
  name?: string
  description?: string
  metaData?: GroupMetaData[]
}

export type GroupMetaData = {
  groupMetaDataOriginalID?: string
  id?: string
  key?: string
  order?: number
  description?: string
  type?: string
  value?: string
  defaultValue?: string
}

export type CreateGroupRequest = {
  group?: Group
}

export type CreateGroupResponse = {
}

export type GetGroupRequest = {
  id?: string
}

export type GetGroupResponse = {
  group?: Group
}

export type UpdateGroupRequest = {
  group?: Group
}

export type UpdateGroupResponse = {
}

export type GroupAddUserRequest = {
  userId?: string
  groupId?: string
  metaData?: GroupMetaData[]
}

export type GroupAddUserResponse = {
}

export type GroupRemoveUserRequest = {
  userId?: string
  groupId?: string
}

export type GroupRemoveUserResponse = {
}

export type GroupUserUpdateRequest = {
  userId?: string
  groupId?: string
  metaData?: GroupMetaData[]
}

export type GroupUserUpdateResponse = {
}

export type GetUserPassKeyRequest = {
  id?: string
}

export type GetUserPassKeyResponse = {
  passKey?: string
}

export class UserService {
  static GetUser(req: GetUserRequest, initReq?: fm.InitReq): Promise<GetUserResponse> {
    return fm.fetchReq<GetUserRequest, GetUserResponse>(`/gapi/user/v1/info?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static UpdateUser(req: UpdateUserRequest, initReq?: fm.InitReq): Promise<UpdateUserResponse> {
    return fm.fetchReq<UpdateUserRequest, UpdateUserResponse>(`/gapi/user/v1/info`, {...initReq, method: "PUT", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdatePassword(req: UpdatePasswordRequest, initReq?: fm.InitReq): Promise<UpdatePasswordResponse> {
    return fm.fetchReq<UpdatePasswordRequest, UpdatePasswordResponse>(`/gapi/user/v1/password`, {...initReq, method: "PUT", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateGroup(req: CreateGroupRequest, initReq?: fm.InitReq): Promise<CreateGroupResponse> {
    return fm.fetchReq<CreateGroupRequest, CreateGroupResponse>(`/gapi/group/v1/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetGroup(req: GetGroupRequest, initReq?: fm.InitReq): Promise<GetGroupResponse> {
    return fm.fetchReq<GetGroupRequest, GetGroupResponse>(`/gapi/group/v1/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateGroup(req: UpdateGroupRequest, initReq?: fm.InitReq): Promise<UpdateGroupResponse> {
    return fm.fetchReq<UpdateGroupRequest, UpdateGroupResponse>(`/gapi/group/v1/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GroupAddUser(req: GroupAddUserRequest, initReq?: fm.InitReq): Promise<GroupAddUserResponse> {
    return fm.fetchReq<GroupAddUserRequest, GroupAddUserResponse>(`/gapi/group/v1/user/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GroupRemoveUser(req: GroupRemoveUserRequest, initReq?: fm.InitReq): Promise<GroupRemoveUserResponse> {
    return fm.fetchReq<GroupRemoveUserRequest, GroupRemoveUserResponse>(`/gapi/group/v1/user/remove`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GroupUserUpdate(req: GroupUserUpdateRequest, initReq?: fm.InitReq): Promise<GroupUserUpdateResponse> {
    return fm.fetchReq<GroupUserUpdateRequest, GroupUserUpdateResponse>(`/gapi/group/v1/user/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetUserPassKey(req: GetUserRequest, initReq?: fm.InitReq): Promise<GetUserResponse> {
    return fm.fetchReq<GetUserRequest, GetUserResponse>(`/gapi/user/v1/passkey`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}