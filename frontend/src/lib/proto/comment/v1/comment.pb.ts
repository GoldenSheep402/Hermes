/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type CommentInfo = {
  id?: string
  resourceId?: string
  authorId?: string
  authorName?: string
  authorAvatar?: string
  body?: string
  parentId?: string
  createdAt?: string
  replies?: CommentInfo[]
}

export type CreateCommentRequest = {
  resourceId?: string
  body?: string
  parentId?: string
}

export type CreateCommentResponse = {
  id?: string
}

export type ListCommentsRequest = {
  resourceId?: string
  page?: number
  pageSize?: number
}

export type ListCommentsResponse = {
  comments?: CommentInfo[]
  total?: string
}

export type UpdateCommentRequest = {
  id?: string
  body?: string
}

export type UpdateCommentResponse = {
}

export type DeleteCommentRequest = {
  id?: string
}

export type DeleteCommentResponse = {
}

export class CommentService {
  static CreateComment(req: CreateCommentRequest, initReq?: fm.InitReq): Promise<CreateCommentResponse> {
    return fm.fetchReq<CreateCommentRequest, CreateCommentResponse>(`/gapi/comment/v1/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListComments(req: ListCommentsRequest, initReq?: fm.InitReq): Promise<ListCommentsResponse> {
    return fm.fetchReq<ListCommentsRequest, ListCommentsResponse>(`/gapi/comment/v1/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateComment(req: UpdateCommentRequest, initReq?: fm.InitReq): Promise<UpdateCommentResponse> {
    return fm.fetchReq<UpdateCommentRequest, UpdateCommentResponse>(`/gapi/comment/v1/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteComment(req: DeleteCommentRequest, initReq?: fm.InitReq): Promise<DeleteCommentResponse> {
    return fm.fetchReq<DeleteCommentRequest, DeleteCommentResponse>(`/gapi/comment/v1/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}