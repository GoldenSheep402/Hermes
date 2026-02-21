/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type MessageInfo = {
  id?: string
  senderId?: string
  senderName?: string
  receiverId?: string
  receiverName?: string
  subject?: string
  body?: string
  isRead?: boolean
  type?: number
  createdAt?: string
}

export type SendMessageRequest = {
  receiverId?: string
  subject?: string
  body?: string
}

export type SendMessageResponse = {
  id?: string
}

export type GetMessageRequest = {
  id?: string
}

export type GetMessageResponse = {
  message?: MessageInfo
}

export type ListInboxRequest = {
  page?: number
  pageSize?: number
}

export type ListInboxResponse = {
  messages?: MessageInfo[]
  total?: string
}

export type ListOutboxRequest = {
  page?: number
  pageSize?: number
}

export type ListOutboxResponse = {
  messages?: MessageInfo[]
  total?: string
}

export type MarkAsReadRequest = {
  ids?: string[]
}

export type MarkAsReadResponse = {
}

export type DeleteMessageRequest = {
  id?: string
}

export type DeleteMessageResponse = {
}

export type GetUnreadCountRequest = {
}

export type GetUnreadCountResponse = {
  count?: string
}

export class MessageService {
  static SendMessage(req: SendMessageRequest, initReq?: fm.InitReq): Promise<SendMessageResponse> {
    return fm.fetchReq<SendMessageRequest, SendMessageResponse>(`/gapi/message/v1/send`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetMessage(req: GetMessageRequest, initReq?: fm.InitReq): Promise<GetMessageResponse> {
    return fm.fetchReq<GetMessageRequest, GetMessageResponse>(`/gapi/message/v1/info`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListInbox(req: ListInboxRequest, initReq?: fm.InitReq): Promise<ListInboxResponse> {
    return fm.fetchReq<ListInboxRequest, ListInboxResponse>(`/gapi/message/v1/inbox`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListOutbox(req: ListOutboxRequest, initReq?: fm.InitReq): Promise<ListOutboxResponse> {
    return fm.fetchReq<ListOutboxRequest, ListOutboxResponse>(`/gapi/message/v1/outbox`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static MarkAsRead(req: MarkAsReadRequest, initReq?: fm.InitReq): Promise<MarkAsReadResponse> {
    return fm.fetchReq<MarkAsReadRequest, MarkAsReadResponse>(`/gapi/message/v1/read`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteMessage(req: DeleteMessageRequest, initReq?: fm.InitReq): Promise<DeleteMessageResponse> {
    return fm.fetchReq<DeleteMessageRequest, DeleteMessageResponse>(`/gapi/message/v1/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetUnreadCount(req: GetUnreadCountRequest, initReq?: fm.InitReq): Promise<GetUnreadCountResponse> {
    return fm.fetchReq<GetUnreadCountRequest, GetUnreadCountResponse>(`/gapi/message/v1/unread`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}