<script lang="ts" setup>
import { CategoryService, TorrentService, TrackerService, UserService } from '@/services/grpc.ts'
import { Notification } from '@arco-design/web-vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
interface TorrentMessage {
  id: string
  name: string
  description: string
  categoryId: string
  categoryName: string
  downloading: number
  seeding: number
  finished: number
}

interface Category {
  id: string
  name: string
  description: string
}

const torrentList = ref<TorrentMessage[]>([])
const categoryList = ref<Category[]>([])
const selectedCategory = ref<string>('')
const searchKeyword = ref<string>('')
const allTorrents = ref<TorrentMessage[]>([])

async function fetchCategoryList() {
  CategoryService.GetCategoryList({}).then((res) => {
    for (let i = 0; i < res.category!.length; i++) {
      categoryList.value.push({
        id: res.category![i].id!,
        name: res.category![i].name!,
        description: res.category![i].description!,
      })
    }
  }).catch((err) => {
    console.error('Failed to fetch category list', err)
  })
}

async function fetchTorrentData() {
  allTorrents.value = []
  const req = {
    categoryId: selectedCategory.value || undefined,
  }
  TorrentService.GetTorrentV1List(req).then(async (res) => {
    for (let i = 0; i < res.torrents!.length; i++) {
      const downloadingCount = ref<number>(0)
      const seedingCount = ref<number>(0)
      const finishedCount = ref<number>(0)
      await TrackerService.GetTorrentDownloadingStatus({ torrentId: res.torrents![i].id })
        .then((statusRes) => {
          downloadingCount.value = statusRes.downloading!
          seedingCount.value = statusRes.seeding!
          finishedCount.value = statusRes.finished!
        })
        .catch((err) => {
          console.error('Failed to get torrent status', err)
        })

      allTorrents.value.push({
        id: res.torrents![i].id!,
        name: res.torrents![i].name!,
        description: res.torrents![i].description!,
        categoryId: res.torrents![i].categoryId!,
        categoryName: res.torrents![i].categoryName!,
        downloading: downloadingCount.value,
        seeding: seedingCount.value,
        finished: finishedCount.value,
      })
    }
    filterTorrents()
  }).catch((err) => {
    console.error('Failed to fetch torrent list', err)
  })
}

function filterTorrents() {
  torrentList.value = allTorrents.value.filter((torrent) => {
    if (searchKeyword.value) {
      return torrent.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
    }
    return true
  })
}

function handleSearch() {
  fetchTorrentData()
}

function handleReset() {
  selectedCategory.value = ''
  searchKeyword.value = ''
  fetchTorrentData()
}

function base64ToUint8Array(base64: string): Uint8Array {
  const binaryString = window.atob(base64)
  const len = binaryString.length
  const bytes = new Uint8Array(len)
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }
  return bytes
}

function handleDetail(id: string) {
  router.push(`/torrent/detail?id=${id}`)
}

function downloadTorrent(id: string, name: string) {
  TorrentService.DownloadTorrentV1({ id }).then((res) => {
    if (res.data) {
      const uint8Array = base64ToUint8Array(res.data)
      const blob = new Blob([uint8Array], { type: 'application/octet-stream' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `${name}.torrent`
      document.body.appendChild(link)
      link.click()

      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    }
    handleNotification('success', '下载成功', '种子下载成功')
  }).catch((err) => {
    handleNotification('error', '下载失败', '种子下载失败')
    console.error('Error downloading torrent:', err)
  })
}

function handleNotification(type: string, title: string, content: string) {
  switch (type) {
    case 'success':
      Notification.success({
        title,
        content,
      })
      break
    case 'error':
      Notification.error({
        title,
        content,
      })
      break
    case 'warning':
      Notification.warning({
        title,
        content,
      })
      break
    default:
      Notification.info({
        title,
        content,
      })
  }
}

const passkey = ref<string>('')

function getPasskey() {
  UserService.GetUserPassKey({}).then((res) => {
    passkey.value = res.passKey!
  }).catch(() => {
    console.log('Get Passkey Fail')
  })
}

function genUrl(id: string) {
  const baseUrl = import.meta.env.VITE_GAPI_URL
  return `${baseUrl}/api/torrent/download/${passkey.value}?id=${id}`
}

onMounted(() => {
  getPasskey()
  fetchCategoryList()
  fetchTorrentData()
})
</script>

<template>
  <div class="p-5">
    <div class="bg-[--color-bg-2] p-5">
      <div class="mb-5 p-0.5 text-20px text-[--color-text-1] font-500 leading-[1.4]">
        种子列表
      </div>

      <a-card class="search-card mb-4" :bordered="true">
        <template #title>
          <div class="flex items-center">
            <icon-search class="mr-2 text-[var(--color-text-3)]" />
            <span class="text-[var(--color-text-1)]">搜索条件</span>
          </div>
        </template>
        <a-form layout="inline" :model="{ category: selectedCategory, keyword: searchKeyword }" @submit="handleSearch">
          <a-form-item field="category" label="类别" class="!mb-0">
            <a-select
              v-model="selectedCategory"
              placeholder="选择类别"
              allow-clear
              style="width: 200px"
            >
              <a-option
                v-for="category in categoryList"
                :key="category.id"
                :value="category.id"
              >
                {{ category.name }}
              </a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="keyword" label="关键词" class="!mb-0">
            <a-input
              v-model="searchKeyword"
              placeholder="搜索种子名称"
              allow-clear
              style="width: 300px"
            />
          </a-form-item>
          <a-form-item class="!mb-0">
            <a-space>
              <a-button type="primary" html-type="submit">
                <template #icon>
                  <icon-search />
                </template>
                搜索
              </a-button>
              <a-button @click="handleReset">
                <template #icon>
                  <icon-refresh />
                </template>
                重置
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>

      <a-table :data="torrentList">
        <template #columns>
          <a-table-column key="name" data-index="name" title="名称" />
          <a-table-column key="categoryName" data-index="categoryName" title="类别名称" />
          <!--          <a-table-column key="finished" dataIndex="finished" title="下载"></a-table-column> -->
          <a-table-column key="status" title="状态">
            <template #cell="{ record }">
              <div class="flex flex-row gap-2">
                <a-statistic :value="record.downloading">
                  <template #suffix>
                    <icon-arrow-down />
                  </template>
                </a-statistic>
                <a-statistic :value="record.seeding">
                  <template #suffix>
                    <icon-arrow-up />
                  </template>
                </a-statistic>

                <a-statistic :value="record.finished">
                  <template #suffix>
                    <icon-check />
                  </template>
                </a-statistic>
              </div>
            </template>
          </a-table-column>
          <a-table-column key="action" title="操作">
            <template #cell="{ record }">
              <div class="flex gap-2">
                <a-button type="primary" @click="handleDetail(record.id)">
                  查看
                </a-button>
                <a-button
                  type="primary"
                  :href="genUrl(record.id)"
                  @click.prevent="downloadTorrent(record.id, record.name)"
                >
                  下载
                </a-button>
              </div>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>
  </div>
</template>

<style lang="less" scoped>
.search-card {
  border: 1px solid var(--color-border);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}
</style>
