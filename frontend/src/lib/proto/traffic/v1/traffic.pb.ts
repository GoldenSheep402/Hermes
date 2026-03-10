/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type UserTrafficInfo = {
  userId?: string
  realUpload?: string
  realDownload?: string
  bonusUpload?: string
  bonusDownload?: string
  ratio?: number
  uploadRate?: string
  downloadRate?: string
}

export type TransferHistoryItem = {
  id?: string
  torrentId?: string
  torrentName?: string
  uploaded?: string
  downloaded?: string
  seedTime?: string
  isFinished?: boolean
  isActive?: boolean
  lastAction?: string
}

export type TorrentStatsInfo = {
  torrentId?: string
  seedCount?: number
  leechCount?: number
  snatchCount?: number
  totalUpload?: string
  totalDownload?: string
}

export type GetUserTrafficRequest = {
  userId?: string
}

export type GetUserTrafficResponse = {
  traffic?: UserTrafficInfo
}

export type ListTransferHistoryRequest = {
  userId?: string
  page?: number
  pageSize?: number
}

export type ListTransferHistoryResponse = {
  items?: TransferHistoryItem[]
  total?: string
}

export type GetTorrentStatsRequest = {
  torrentId?: string
}

export type GetTorrentStatsResponse = {
  stats?: TorrentStatsInfo
}

export class TrafficService {
  static GetUserTraffic(req: GetUserTrafficRequest, initReq?: fm.InitReq): Promise<GetUserTrafficResponse> {
    return fm.fetchReq<GetUserTrafficRequest, GetUserTrafficResponse>(`/gapi/traffic/v1/user`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListTransferHistory(req: ListTransferHistoryRequest, initReq?: fm.InitReq): Promise<ListTransferHistoryResponse> {
    return fm.fetchReq<ListTransferHistoryRequest, ListTransferHistoryResponse>(`/gapi/traffic/v1/history`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetTorrentStats(req: GetTorrentStatsRequest, initReq?: fm.InitReq): Promise<GetTorrentStatsResponse> {
    return fm.fetchReq<GetTorrentStatsRequest, GetTorrentStatsResponse>(`/gapi/traffic/v1/torrent`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}