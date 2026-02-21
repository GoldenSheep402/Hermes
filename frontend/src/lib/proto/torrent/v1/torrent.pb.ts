/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type TorrentInfo = {
  id?: string
  infoHash?: string
  uploaderId?: string
  name?: string
  size?: string
  isSingleFile?: boolean
  fileCount?: number
  seedCount?: number
  leechCount?: number
  snatchCount?: number
  isActive?: boolean
  createdAt?: string
}

export type TorrentFile = {
  id?: string
  torrentId?: string
  path?: string
  size?: string
}

export type UploadTorrentRequest = {
  torrentData?: Uint8Array
}

export type UploadTorrentResponse = {
  torrentId?: string
  infoHash?: string
}

export type DownloadTorrentRequest = {
  torrentId?: string
}

export type DownloadTorrentResponse = {
  torrentData?: Uint8Array
}

export type GetTorrentRequest = {
  id?: string
}

export type GetTorrentResponse = {
  torrent?: TorrentInfo
}

export type ListTorrentFilesRequest = {
  torrentId?: string
}

export type ListTorrentFilesResponse = {
  files?: TorrentFile[]
}

export type DeleteTorrentRequest = {
  id?: string
}

export type DeleteTorrentResponse = {
}

export class TorrentService {
  static UploadTorrent(req: UploadTorrentRequest, initReq?: fm.InitReq): Promise<UploadTorrentResponse> {
    return fm.fetchReq<UploadTorrentRequest, UploadTorrentResponse>(`/gapi/torrent/v1/upload`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DownloadTorrent(req: DownloadTorrentRequest, initReq?: fm.InitReq): Promise<DownloadTorrentResponse> {
    return fm.fetchReq<DownloadTorrentRequest, DownloadTorrentResponse>(`/gapi/torrent/v1/download`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetTorrent(req: GetTorrentRequest, initReq?: fm.InitReq): Promise<GetTorrentResponse> {
    return fm.fetchReq<GetTorrentRequest, GetTorrentResponse>(`/gapi/torrent/v1/info`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListTorrentFiles(req: ListTorrentFilesRequest, initReq?: fm.InitReq): Promise<ListTorrentFilesResponse> {
    return fm.fetchReq<ListTorrentFilesRequest, ListTorrentFilesResponse>(`/gapi/torrent/v1/files`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteTorrent(req: DeleteTorrentRequest, initReq?: fm.InitReq): Promise<DeleteTorrentResponse> {
    return fm.fetchReq<DeleteTorrentRequest, DeleteTorrentResponse>(`/gapi/torrent/v1/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}