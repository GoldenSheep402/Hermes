/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type PeerInfo = {
  peerId?: Uint8Array
  ip?: string
  port?: number
}

export type SnatchInfo = {
  id?: string
  torrentId?: string
  userId?: string
  uploaded?: string
  downloaded?: string
  seedTime?: string
  isActive?: boolean
  finishedAt?: string
  lastAction?: string
}

export type GetTorrentPeersRequest = {
  torrentId?: string
}

export type GetTorrentPeersResponse = {
  seederCount?: number
  leecherCount?: number
  seeders?: PeerInfo[]
  leechers?: PeerInfo[]
}

export type ListSnatchesRequest = {
  torrentId?: string
  page?: number
  pageSize?: number
}

export type ListSnatchesResponse = {
  snatches?: SnatchInfo[]
  total?: string
}

export type GetUserSnatchesRequest = {
  userId?: string
  page?: number
  pageSize?: number
}

export type GetUserSnatchesResponse = {
  snatches?: SnatchInfo[]
  total?: string
}

export type AnnounceRequest = {
  passkey?: string
  infoHash?: string
  peerId?: string
  port?: number
  uploaded?: string
  downloaded?: string
  left?: string
  event?: string
  ip?: string
  numWant?: number
  compact?: number
}

export type AnnounceResponse = {
  failureReason?: string
  warningMessage?: string
  interval?: number
  minInterval?: number
  complete?: number
  incomplete?: number
  peers?: PeerInfo[]
  peersCompact?: Uint8Array
}

export type ScrapeRequest = {
  passkey?: string
  infoHashes?: string[]
}

export type ScrapeFile = {
  complete?: number
  incomplete?: number
  downloaded?: number
}

export type ScrapeResponse = {
  files?: {[key: string]: ScrapeFile}
}

export class TrackerService {
  static GetTorrentPeers(req: GetTorrentPeersRequest, initReq?: fm.InitReq): Promise<GetTorrentPeersResponse> {
    return fm.fetchReq<GetTorrentPeersRequest, GetTorrentPeersResponse>(`/gapi/tracker/v1/peers`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListSnatches(req: ListSnatchesRequest, initReq?: fm.InitReq): Promise<ListSnatchesResponse> {
    return fm.fetchReq<ListSnatchesRequest, ListSnatchesResponse>(`/gapi/tracker/v1/snatches`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetUserSnatches(req: GetUserSnatchesRequest, initReq?: fm.InitReq): Promise<GetUserSnatchesResponse> {
    return fm.fetchReq<GetUserSnatchesRequest, GetUserSnatchesResponse>(`/gapi/tracker/v1/user/snatches`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}