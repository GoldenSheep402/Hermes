/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type User = {
  id?: string
  nickname?: string
  role?: string
}

export type GetUserRequest = {
  id?: string
}

export type GetUserResponse = {
  user?: User
}

export type UpdateUserRequest = {
  username?: string
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
}