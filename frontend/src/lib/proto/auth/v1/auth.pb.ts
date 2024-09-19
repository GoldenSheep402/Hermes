/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type RegisterSendEmailRequest = {
  email?: string
}

export type RegisterSendEmailResponse = {
}

export type RegisterWithEmailRequest = {
  username?: string
  email?: string
  password?: string
  emailToken?: string
}

export type RegisterWithEmailResponse = {
}

export type LoginRequest = {
  email?: string
  password?: string
}

export type LoginResponse = {
  accessToken?: string
  refreshToken?: string
}

export type RefreshTokenRequest = {
  refreshToken?: string
}

export type RefreshTokenResponse = {
  accessToken?: string
  refreshToken?: string
}

export class AuthService {
  static RegisterSendEmail(req: RegisterSendEmailRequest, initReq?: fm.InitReq): Promise<RegisterSendEmailResponse> {
    return fm.fetchReq<RegisterSendEmailRequest, RegisterSendEmailResponse>(`/gapi/auth/v1/register/send/email`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static RegisterWithEmail(req: RegisterWithEmailRequest, initReq?: fm.InitReq): Promise<RegisterWithEmailResponse> {
    return fm.fetchReq<RegisterWithEmailRequest, RegisterWithEmailResponse>(`/gapi/auth/v1/register/email`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Login(req: LoginRequest, initReq?: fm.InitReq): Promise<LoginResponse> {
    return fm.fetchReq<LoginRequest, LoginResponse>(`/gapi/auth/v1/login`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static RefreshToken(req: RefreshTokenRequest, initReq?: fm.InitReq): Promise<RefreshTokenResponse> {
    return fm.fetchReq<RefreshTokenRequest, RefreshTokenResponse>(`/gapi/auth/v1/refreshToken`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}