/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type Settings = {
  peerExpireTime?: number
  allowedNets?: string[]
  smtpEnable?: boolean
  smtpHost?: string
  smtpPort?: number
  smtpUser?: string
  smtpPassword?: string
  registerEnable?: boolean
  loginEnable?: boolean
  publishEnable?: boolean
  innetTrackerAddrs?: string[]
}

export type GetSettingsRequest = {
}

export type GetSettingsResponse = {
  settings?: Settings
}

export type SetSettingsRequest = {
  settings?: Settings
}

export type SetSettingsResponse = {
}

export class SystemService {
  static GetSettings(req: GetSettingsRequest, initReq?: fm.InitReq): Promise<GetSettingsResponse> {
    return fm.fetchReq<GetSettingsRequest, GetSettingsResponse>(`/gapi/system/v1/settings/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static SetSettings(req: SetSettingsRequest, initReq?: fm.InitReq): Promise<SetSettingsResponse> {
    return fm.fetchReq<SetSettingsRequest, SetSettingsResponse>(`/gapi/system/v1/settings/set`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}