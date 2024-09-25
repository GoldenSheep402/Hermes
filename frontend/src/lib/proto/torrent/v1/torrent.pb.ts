/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type Torrent = {
  data?: Uint8Array
}

export type TorrentMetaData = {
  id?: string
  categoryId?: string
  torrentId?: string
  key?: string
  order?: number
  description?: string
  type?: string
  value?: string
}

export type GetTorrentV1Request = {
  id?: string
}

export type GetTorrentV1Response = {
  metadata?: TorrentMetaData[]
}

export type CreateTorrentV1Request = {
  categoryId?: string
  metadata?: TorrentMetaData[]
  torrent?: Torrent
}

export type CreateTorrentV1Response = {
  id?: string
}

export type DownloadTorrentV1Request = {
  id?: string
}

export type DownloadTorrentV1Response = {
  data?: Uint8Array
}

export class TorrentService {
  static GetTorrentV1(req: GetTorrentV1Request, initReq?: fm.InitReq): Promise<GetTorrentV1Response> {
    return fm.fetchReq<GetTorrentV1Request, GetTorrentV1Response>(`/gapi/torrent/v1/info/v1`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateTorrentV1(req: CreateTorrentV1Request, initReq?: fm.InitReq): Promise<CreateTorrentV1Response> {
    return fm.fetchReq<CreateTorrentV1Request, CreateTorrentV1Response>(`/gapi/torrent/v1/create/v1`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DownloadTorrentV1(req: DownloadTorrentV1Request, initReq?: fm.InitReq): Promise<DownloadTorrentV1Response> {
    return fm.fetchReq<DownloadTorrentV1Request, DownloadTorrentV1Response>(`/gapi/torrent/v1/download/v1`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}