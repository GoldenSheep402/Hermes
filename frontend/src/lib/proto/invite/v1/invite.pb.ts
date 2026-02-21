/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type InviteCodeInfo = {
  id?: string
  code?: string
  senderId?: string
  senderName?: string
  receiverId?: string
  receiverName?: string
  isUsed?: boolean
  expiredAt?: string
  usedAt?: string
  createdAt?: string
}

export type CreateInviteCodeRequest = {
  count?: number
}

export type CreateInviteCodeResponse = {
  codes?: string[]
}

export type ListInviteCodesRequest = {
  page?: number
  pageSize?: number
}

export type ListInviteCodesResponse = {
  codes?: InviteCodeInfo[]
  total?: string
}

export type ConsumeInviteCodeRequest = {
  code?: string
}

export type ConsumeInviteCodeResponse = {
  success?: boolean
}

export class InviteService {
  static CreateInviteCode(req: CreateInviteCodeRequest, initReq?: fm.InitReq): Promise<CreateInviteCodeResponse> {
    return fm.fetchReq<CreateInviteCodeRequest, CreateInviteCodeResponse>(`/gapi/invite/v1/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListInviteCodes(req: ListInviteCodesRequest, initReq?: fm.InitReq): Promise<ListInviteCodesResponse> {
    return fm.fetchReq<ListInviteCodesRequest, ListInviteCodesResponse>(`/gapi/invite/v1/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ConsumeInviteCode(req: ConsumeInviteCodeRequest, initReq?: fm.InitReq): Promise<ConsumeInviteCodeResponse> {
    return fm.fetchReq<ConsumeInviteCodeRequest, ConsumeInviteCodeResponse>(`/gapi/invite/v1/consume`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}