/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type GetTorrentDownloadingStatusRequest = {
  torrentId?: string
}

export type GetTorrentDownloadingStatusResponse = {
  downloading?: number
  finished?: number
  seeding?: number
}

export type GetTrackerRequest = {
  key?: string
  info_hash?: string
  peer_id?: string
  port?: number
  uploaded?: string
  downloaded?: string
  left?: string
  event?: string
  ip?: string
  num_want?: number
  compact?: number
  no_peer_id?: number
  corrupt?: number
  support_crypto?: number
  redundant?: number
}

export type GetTrackerResponse = {
  response?: string
}

export type GetTrackerResponseDetail = {
  failure_reason?: string
  warning_message?: string
  interval?: number
  min_interval?: number
  tracker_id?: string
  complete?: number
  incomplete?: number
  peers?: PeerInfo[]
  peers_compact?: Uint8Array
}

export type PeerInfo = {
  peer_id?: Uint8Array
  ip?: string
  port?: number
}

export class TrackerService {
  static GetTorrentDownloadingStatus(req: GetTorrentDownloadingStatusRequest, initReq?: fm.InitReq): Promise<GetTorrentDownloadingStatusResponse> {
    return fm.fetchReq<GetTorrentDownloadingStatusRequest, GetTorrentDownloadingStatusResponse>(`/gapi/trackerV1/v1/status`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}