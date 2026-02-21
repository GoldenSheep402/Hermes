/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type Resource = {
  id?: string
  title?: string
  subtitle?: string
  description?: string
  categoryId?: string
  torrentId?: string
  uploaderId?: string
  status?: number
  isSticky?: boolean
  isFree?: boolean
  freeUntil?: string
  doubleUpload?: boolean
  doubleUntil?: string
  viewCount?: number
  commentCount?: number
  thankCount?: number
  createdAt?: string
  updatedAt?: string
  uploaderName?: string
  categoryName?: string
  torrentSize?: string
  seedCount?: number
  leechCount?: number
  snatchCount?: number
  metadata?: ResourceMeta[]
  tags?: Tag[]
  screenshots?: ResourceScreenshot[]
}

export type ResourceMeta = {
  id?: string
  key?: string
  value?: string
}

export type Tag = {
  id?: string
  name?: string
}

export type ResourceScreenshot = {
  id?: string
  url?: string
  sortOrder?: number
}

export type CreateResourceRequest = {
  title?: string
  subtitle?: string
  description?: string
  categoryId?: string
  torrentData?: Uint8Array
  metadata?: ResourceMeta[]
  tagNames?: string[]
  screenshotUrls?: string[]
}

export type CreateResourceResponse = {
  id?: string
}

export type GetResourceRequest = {
  id?: string
}

export type GetResourceResponse = {
  resource?: Resource
}

export type ListResourcesRequest = {
  categoryId?: string
  keyword?: string
  status?: number
  page?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: string
}

export type ListResourcesResponse = {
  resources?: Resource[]
  total?: string
}

export type UpdateResourceRequest = {
  id?: string
  title?: string
  subtitle?: string
  description?: string
  categoryId?: string
  metadata?: ResourceMeta[]
  tagNames?: string[]
  screenshotUrls?: string[]
}

export type UpdateResourceResponse = {
}

export type DeleteResourceRequest = {
  id?: string
}

export type DeleteResourceResponse = {
}

export type ThankResourceRequest = {
  resourceId?: string
}

export type ThankResourceResponse = {
}

export type SetPromotionRequest = {
  resourceId?: string
  isFree?: boolean
  freeUntil?: string
  doubleUpload?: boolean
  doubleUntil?: string
}

export type SetPromotionResponse = {
}

export class ResourceService {
  static CreateResource(req: CreateResourceRequest, initReq?: fm.InitReq): Promise<CreateResourceResponse> {
    return fm.fetchReq<CreateResourceRequest, CreateResourceResponse>(`/gapi/resource/v1/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetResource(req: GetResourceRequest, initReq?: fm.InitReq): Promise<GetResourceResponse> {
    return fm.fetchReq<GetResourceRequest, GetResourceResponse>(`/gapi/resource/v1/info`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListResources(req: ListResourcesRequest, initReq?: fm.InitReq): Promise<ListResourcesResponse> {
    return fm.fetchReq<ListResourcesRequest, ListResourcesResponse>(`/gapi/resource/v1/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateResource(req: UpdateResourceRequest, initReq?: fm.InitReq): Promise<UpdateResourceResponse> {
    return fm.fetchReq<UpdateResourceRequest, UpdateResourceResponse>(`/gapi/resource/v1/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteResource(req: DeleteResourceRequest, initReq?: fm.InitReq): Promise<DeleteResourceResponse> {
    return fm.fetchReq<DeleteResourceRequest, DeleteResourceResponse>(`/gapi/resource/v1/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ThankResource(req: ThankResourceRequest, initReq?: fm.InitReq): Promise<ThankResourceResponse> {
    return fm.fetchReq<ThankResourceRequest, ThankResourceResponse>(`/gapi/resource/v1/thank`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static SetPromotion(req: SetPromotionRequest, initReq?: fm.InitReq): Promise<SetPromotionResponse> {
    return fm.fetchReq<SetPromotionRequest, SetPromotionResponse>(`/gapi/resource/v1/promotion`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}