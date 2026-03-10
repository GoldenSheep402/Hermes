/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
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

export type StreamSiteTrafficRequest = {
  intervalSeconds?: number
  smoothingFactor?: number
}

export type SiteTrafficRatePoint = {
  timestamp?: string
  uploadRate?: string
  downloadRate?: string
  rawUploadRate?: string
  rawDownloadRate?: string
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
  static ListTransferHistory(req: ListTransferHistoryRequest, initReq?: fm.InitReq): Promise<ListTransferHistoryResponse> {
    return fm.fetchReq<ListTransferHistoryRequest, ListTransferHistoryResponse>(`/gapi/traffic/v1/history`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetTorrentStats(req: GetTorrentStatsRequest, initReq?: fm.InitReq): Promise<GetTorrentStatsResponse> {
    return fm.fetchReq<GetTorrentStatsRequest, GetTorrentStatsResponse>(`/gapi/traffic/v1/torrent`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static StreamSiteTraffic(req: StreamSiteTrafficRequest, entityNotifier?: fm.NotifyStreamEntityArrival<SiteTrafficRatePoint>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<StreamSiteTrafficRequest, SiteTrafficRatePoint>(`/gapi/traffic/v1/site/stream?${fm.renderURLSearchParams(req, [])}`, entityNotifier, {...initReq, method: "GET"})
  }
}