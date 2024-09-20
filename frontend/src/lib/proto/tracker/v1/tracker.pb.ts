/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type GetTrackerRequest = {
  key?: string
  info_hash?: Uint8Array
  peer_id?: Uint8Array
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
  static GetTracker(req: GetTrackerRequest, initReq?: fm.InitReq): Promise<GetTrackerResponse> {
    return fm.fetchReq<GetTrackerRequest, GetTrackerResponse>(`/trackerV1/${req["key"]}?${fm.renderURLSearchParams(req, ["key"])}`, {...initReq, method: "GET"})
  }
}