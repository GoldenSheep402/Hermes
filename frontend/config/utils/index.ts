/**
 * Whether to generate package preview
 * 是否生成打包报告
 */
export default {}

export function isReportMode(): boolean {
  try {
    return import.meta.env?.VITE_REPORT === 'true'
  }
  catch {
    return false
  }
}
