/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type Category = {
  id?: string
  name?: string
  slug?: string
  description?: string
  parentId?: string
  icon?: string
  sortOrder?: number
  isEnabled?: boolean
  metaTemplates?: CategoryMetaTemplate[]
  children?: Category[]
}

export type CategoryMetaTemplate = {
  id?: string
  categoryId?: string
  key?: string
  label?: string
  type?: string
  required?: boolean
  options?: string
  sortOrder?: number
  defaultValue?: string
}

export type MetaTemplatePreset = {
  value?: string
  label?: string
  description?: string
  templates?: CategoryMetaTemplate[]
}

export type CreateCategoryRequest = {
  name?: string
  slug?: string
  description?: string
  parentId?: string
  icon?: string
  sortOrder?: number
}

export type CreateCategoryResponse = {
  id?: string
}

export type GetCategoryRequest = {
  id?: string
}

export type GetCategoryResponse = {
  category?: Category
}

export type ListCategoriesRequest = {
}

export type ListCategoriesResponse = {
  categories?: Category[]
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

export type CreateMetaTemplateRequest = {
  template?: CategoryMetaTemplate
}

export type CreateMetaTemplateResponse = {
  id?: string
}

export type UpdateMetaTemplateRequest = {
  template?: CategoryMetaTemplate
}

export type UpdateMetaTemplateResponse = {
}

export type DeleteMetaTemplateRequest = {
  id?: string
}

export type DeleteMetaTemplateResponse = {
}

export type ListMetaTemplatePresetsRequest = {
}

export type ListMetaTemplatePresetsResponse = {
  presets?: MetaTemplatePreset[]
}

export class CategoryService {
  static CreateCategory(req: CreateCategoryRequest, initReq?: fm.InitReq): Promise<CreateCategoryResponse> {
    return fm.fetchReq<CreateCategoryRequest, CreateCategoryResponse>(`/gapi/category/v1/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetCategory(req: GetCategoryRequest, initReq?: fm.InitReq): Promise<GetCategoryResponse> {
    return fm.fetchReq<GetCategoryRequest, GetCategoryResponse>(`/gapi/category/v1/info`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListCategories(req: ListCategoriesRequest, initReq?: fm.InitReq): Promise<ListCategoriesResponse> {
    return fm.fetchReq<ListCategoriesRequest, ListCategoriesResponse>(`/gapi/category/v1/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateCategory(req: UpdateCategoryRequest, initReq?: fm.InitReq): Promise<UpdateCategoryResponse> {
    return fm.fetchReq<UpdateCategoryRequest, UpdateCategoryResponse>(`/gapi/category/v1/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteCategory(req: DeleteCategoryRequest, initReq?: fm.InitReq): Promise<DeleteCategoryResponse> {
    return fm.fetchReq<DeleteCategoryRequest, DeleteCategoryResponse>(`/gapi/category/v1/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static CreateMetaTemplate(req: CreateMetaTemplateRequest, initReq?: fm.InitReq): Promise<CreateMetaTemplateResponse> {
    return fm.fetchReq<CreateMetaTemplateRequest, CreateMetaTemplateResponse>(`/gapi/category/v1/meta/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateMetaTemplate(req: UpdateMetaTemplateRequest, initReq?: fm.InitReq): Promise<UpdateMetaTemplateResponse> {
    return fm.fetchReq<UpdateMetaTemplateRequest, UpdateMetaTemplateResponse>(`/gapi/category/v1/meta/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteMetaTemplate(req: DeleteMetaTemplateRequest, initReq?: fm.InitReq): Promise<DeleteMetaTemplateResponse> {
    return fm.fetchReq<DeleteMetaTemplateRequest, DeleteMetaTemplateResponse>(`/gapi/category/v1/meta/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListMetaTemplatePresets(req: ListMetaTemplatePresetsRequest, initReq?: fm.InitReq): Promise<ListMetaTemplatePresetsResponse> {
    return fm.fetchReq<ListMetaTemplatePresetsRequest, ListMetaTemplatePresetsResponse>(`/gapi/category/v1/meta/presets`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}