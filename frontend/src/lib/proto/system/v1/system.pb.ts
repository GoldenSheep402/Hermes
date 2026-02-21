/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type Setting = {
  key?: string
  value?: string
  type?: string
  desc?: string
}

export type GetSettingsRequest = {
}

export type GetSettingsResponse = {
  settings?: Setting[]
}

export type GetSettingRequest = {
  key?: string
}

export type GetSettingResponse = {
  setting?: Setting
}

export type SetSettingsRequest = {
  settings?: Setting[]
}

export type SetSettingsResponse = {
}

export type DeleteSettingRequest = {
  key?: string
}

export type DeleteSettingResponse = {
}

export type GetSiteStatsRequest = {
}

export type GetSiteStatsResponse = {
  totalUsers?: string
  totalTorrents?: string
  totalResources?: string
  totalTraffic?: string
  totalSeeders?: string
  totalLeechers?: string
}

export class SystemService {
  static GetSettings(req: GetSettingsRequest, initReq?: fm.InitReq): Promise<GetSettingsResponse> {
    return fm.fetchReq<GetSettingsRequest, GetSettingsResponse>(`/gapi/system/v1/settings/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetSetting(req: GetSettingRequest, initReq?: fm.InitReq): Promise<GetSettingResponse> {
    return fm.fetchReq<GetSettingRequest, GetSettingResponse>(`/gapi/system/v1/setting/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static SetSettings(req: SetSettingsRequest, initReq?: fm.InitReq): Promise<SetSettingsResponse> {
    return fm.fetchReq<SetSettingsRequest, SetSettingsResponse>(`/gapi/system/v1/settings/set`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteSetting(req: DeleteSettingRequest, initReq?: fm.InitReq): Promise<DeleteSettingResponse> {
    return fm.fetchReq<DeleteSettingRequest, DeleteSettingResponse>(`/gapi/system/v1/setting/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetSiteStats(req: GetSiteStatsRequest, initReq?: fm.InitReq): Promise<GetSiteStatsResponse> {
    return fm.fetchReq<GetSiteStatsRequest, GetSiteStatsResponse>(`/gapi/system/v1/stats`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}