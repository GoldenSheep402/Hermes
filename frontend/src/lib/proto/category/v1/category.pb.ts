/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type Category = {
  id?: string
  name?: string
  description?: string
  metaData?: CategoryMetaData[]
}

export type CategoryMetaData = {
  id?: string
  categoryId?: string
  key?: string
  order?: number
  description?: string
  type?: string
  value?: string
  defaultValue?: string
}

export type GetCategoryRequest = {
  id?: string
}

export type GetCategoryResponse = {
  category?: Category
}

export type CreateCategoryRequest = {
  category?: Category
}

export type CreateCategoryResponse = {
}

export type UpdateCategoryRequest = {
  category?: Category
}

export type UpdateCategoryResponse = {
}

export type DeleteCategoryRequest = {
  id?: string
}

export type DeleteCategoryResponse = {
}

export class CategoryService {
  static CreateCategory(req: CreateCategoryRequest, initReq?: fm.InitReq): Promise<CreateCategoryResponse> {
    return fm.fetchReq<CreateCategoryRequest, CreateCategoryResponse>(`/gapi/category/v1/info`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetCategory(req: GetCategoryRequest, initReq?: fm.InitReq): Promise<GetCategoryResponse> {
    return fm.fetchReq<GetCategoryRequest, GetCategoryResponse>(`/gapi/category/v1/info?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static UpdateCategory(req: UpdateCategoryRequest, initReq?: fm.InitReq): Promise<UpdateCategoryResponse> {
    return fm.fetchReq<UpdateCategoryRequest, UpdateCategoryResponse>(`/gapi/category/v1/info`, {...initReq, method: "PUT", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteCategory(req: DeleteCategoryRequest, initReq?: fm.InitReq): Promise<DeleteCategoryResponse> {
    return fm.fetchReq<DeleteCategoryRequest, DeleteCategoryResponse>(`/gapi/category/v1/info`, {...initReq, method: "DELETE"})
  }
}