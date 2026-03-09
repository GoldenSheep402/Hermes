export type PromotionKey = 'normal' | 'free' | '2xfree' | '50down'

export interface TorrentCategoryOption {
  id: string
  name: string
  icon: string
}

export interface TorrentRecord {
  id: string
  resourceId: string
  torrentId: string
  categoryId: string
  categoryLabel: string
  categoryIcon: string
  title: string
  subtitle: string
  tags: string[]
  comments: number
  files: number
  addedAt: string
  size: number
  seeders: number
  leechers: number
  completed: number
  uploader: string
  anonymous: boolean
  imdbId: string
  imdbRating: number
  poster: string
  source: string
  resolution: string
  promotion: PromotionKey
  status: 'alive' | 'dead'
  downloaded: boolean
}

export interface TorrentFilterState {
  keyword: string
  categories: string[]
  source: string
  resolution: string
  promotion: '' | PromotionKey
  status: 'all' | 'alive' | 'dead'
  undownloadedOnly: boolean
}

export interface TorrentFileRow {
  id: string
  path: string
  size: number
}
