/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type BonusLogItem = {
  id?: string
  userId?: string
  amount?: string
  reason?: string
  description?: string
  relatedId?: string
  createdAt?: string
}

export type ShopItem = {
  id?: string
  name?: string
  pricePoints?: string
  unit?: string
}

export type GetBonusBalanceRequest = {
}

export type GetBonusBalanceResponse = {
  balance?: string
  balancePoints?: string
}

export type ListBonusLogsRequest = {
  page?: number
  pageSize?: number
}

export type ListBonusLogsResponse = {
  logs?: BonusLogItem[]
  total?: string
}

export type ListShopItemsRequest = {
}

export type ListShopItemsResponse = {
  items?: ShopItem[]
}

export type ExchangeBonusRequest = {
  item?: string
  quantity?: string
}

export type ExchangeBonusResponse = {
  remainingBalance?: string
}

export class BonusService {
  static GetBonusBalance(req: GetBonusBalanceRequest, initReq?: fm.InitReq): Promise<GetBonusBalanceResponse> {
    return fm.fetchReq<GetBonusBalanceRequest, GetBonusBalanceResponse>(`/gapi/bonus/v1/balance`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListBonusLogs(req: ListBonusLogsRequest, initReq?: fm.InitReq): Promise<ListBonusLogsResponse> {
    return fm.fetchReq<ListBonusLogsRequest, ListBonusLogsResponse>(`/gapi/bonus/v1/logs`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListShopItems(req: ListShopItemsRequest, initReq?: fm.InitReq): Promise<ListShopItemsResponse> {
    return fm.fetchReq<ListShopItemsRequest, ListShopItemsResponse>(`/gapi/bonus/v1/shop`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ExchangeBonus(req: ExchangeBonusRequest, initReq?: fm.InitReq): Promise<ExchangeBonusResponse> {
    return fm.fetchReq<ExchangeBonusRequest, ExchangeBonusResponse>(`/gapi/bonus/v1/exchange`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}