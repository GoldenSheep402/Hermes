/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as CategoryV1Category from "../../category/v1/category.pb"
import * as fm from "../../fetch.pb"
export type Torrent = {
  data?: Uint8Array
}

export type GetTorrentV1Request = {
  id?: string
}

export type GetTorrentV1Response = {
  category?: CategoryV1Category.Category
  torrent?: Torrent
}

export type CreateTorrentV1Request = {
  category?: CategoryV1Category.Category
  torrent?: Torrent
}

export type CreateTorrentV1Response = {
  id?: string
}

export class TorrentService {
  static GetTorrentV1(req: GetTorrentV1Request, initReq?: fm.InitReq): Promise<GetTorrentV1Response> {
    return fm.fetchReq<GetTorrentV1Request, GetTorrentV1Response>(`/gapi/torrent/v1/info/v1`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateTorrentV1(req: CreateTorrentV1Request, initReq?: fm.InitReq): Promise<CreateTorrentV1Response> {
    return fm.fetchReq<CreateTorrentV1Request, CreateTorrentV1Response>(`/gapi/torrent/v1/create/v1`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}