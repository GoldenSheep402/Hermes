import type { RouteRecordRaw } from 'vue-router'
import { NOT_FOUND_ROUTE, publicRoutes } from './base'

type RouteModule = {
  default: RouteRecordRaw | RouteRecordRaw[]
}

function formatModules(modules: Record<string, RouteModule>): RouteRecordRaw[] {
  const result: RouteRecordRaw[] = []
  Object.keys(modules).forEach((key) => {
    const routeModule = modules[key]
    if (!routeModule || !routeModule.default) {
      return
    }

    const moduleRoutes = Array.isArray(routeModule.default)
      ? routeModule.default
      : [routeModule.default]

    result.push(...moduleRoutes)
  })
  return result
}

const modules = import.meta.glob('./modules/*.ts', { eager: true }) as Record<
  string,
  RouteModule
>

export const appRoutes: RouteRecordRaw[] = formatModules(modules)

export const routes: RouteRecordRaw[] = [...publicRoutes, ...appRoutes, NOT_FOUND_ROUTE]
