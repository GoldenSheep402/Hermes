<script lang="ts" setup>
import type { Category as CategoryBase } from '@/lib/proto/category/v1/category.pb.ts'
import type { CreateTorrentV1Request } from '@/lib/proto/torrent/v1/torrent.pb.ts'
import type { FileItem } from '@arco-design/web-vue'
import RichTextEditor from '@/components/Editor/index.vue'
import { CategoryService, TorrentService } from '@/services/grpc.ts'
import { Notification } from '@arco-design/web-vue'
import { onMounted, ref } from 'vue'

interface Category {
  id: string
  name: string
  description: string
}

const checkedID = ref('')
const categoryFullInfo = ref<CategoryBase[]>([])
const torrentComment = ref('')
const categoryList = ref<Category[]>([] as Category[])
const fileList = ref<typeof FileItem[]>([])
const uint8Array = ref<Uint8Array | null>(null)

function fetchCategoryList() {
  CategoryService.GetCategoryList({}).then((res) => {
    for (let i = 0; i < res.category!.length; i++) {
      categoryList.value.push({
        id: res.category![i].id!,
        name: res.category![i].name!,
        description: res.category![i].description!,
      })

      CategoryService.GetCategory({ id: res.category![i].id! }).then((_res) => {
        for (let j = 0; j < _res.category!.metaData!.length; j++) {
          categoryFullInfo.value.push({
            id: _res.category!.id!,
            name: _res.category!.name!,
            description: _res.category!.description!,
            metaData: [
              {
                type: _res.category!.metaData![j].type!,
                id: _res.category!.metaData![j].id!,
                order: _res.category!.metaData![j].order!,
                categoryId: _res.category!.metaData![j].categoryId!,
                description: _res.category!.metaData![j].description!,
                key: _res.category!.metaData![j].key!,
                defaultValue: _res.category!.metaData![j].defaultValue!,
                value: _res.category!.metaData![j].value!,
              },
            ],
          })
        }
      }).catch((err) => {
        console.error(err)
      })
    }
  }).catch((err) => {
    console.error(err)
  })
}

function sendFile() {
  const req = ref<CreateTorrentV1Request>({})
  req.value.comment = torrentComment.value
  req.value.categoryId = checkedID.value

  if (!req.value.metadata) {
    req.value.metadata = []
  }

  for (let i = 0; i < categoryFullInfo.value.length; i++) {
    if (categoryFullInfo.value[i].id === checkedID.value) {
      if (categoryFullInfo.value[i].metaData) {
        for (let j = 0; j < categoryFullInfo.value[i].metaData!.length; j++) {
          req.value.metadata!.push({
            id: categoryFullInfo.value[i].metaData![j].id,
            categoryId: checkedID.value,
            value: categoryFullInfo.value[i].metaData![j].value?.toString(),
          })
        }
      }
    }
  }

  if (uint8Array.value) {
    req.value.torrent = {
      data: uint8Array.value,
    }
  }

  console.log('CreateTorrentV1Request:', req.value)

  TorrentService.CreateTorrentV1(req.value).then((res) => {
    console.log('CreateTorrentV1:', res)
    handleNotification('success', '发布成功', '种子发布成功')
  }).catch(() => {
    handleNotification('error', '发布失败', '种子发布失败')
  }).finally(() => {

  })
}

function handleFileChange(fileList: typeof FileItem[], file: typeof FileItem) {
  const fileObj = file.file as File
  if (fileObj) {
    const reader = new FileReader()
    reader.onload = (event) => {
      const content = event.target?.result as ArrayBuffer
      uint8Array.value = new Uint8Array(content)
      console.log('File Content as Uint8Array:', uint8Array.value)
    }
    reader.readAsArrayBuffer(fileObj)
  }
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

onMounted(() => {
  fetchCategoryList()
})
</script>

<template>
  <div class="p-5">
    <div class="bg-[--color-bg-2] p-5 shadow-sm">
      <div class="mb-8 flex items-center p-0.5 text-24px text-[--color-text-1] font-600 leading-[1.4]">
        <icon-upload class="mr-2 text-[var(--color-primary-6)]" />
        发布种子
      </div>

      <a-card class="mb-6" :bordered="true">
        <template #title>
          <div class="flex items-center text-[var(--color-text-1)]">
            <icon-folder class="mr-2 text-[var(--color-primary-6)]" />
            选择类别
          </div>
        </template>
        <div class="p-4">
          <a-select
            v-model="checkedID"
            placeholder="请选择类别"
            allow-clear
            class="mb-6 w-full"
          >
            <a-option
              v-for="category in categoryList"
              :key="category.id"
              :value="category.id"
            >
              <div class="flex items-center">
                <icon-folder class="mr-2 text-[var(--color-text-3)]" />
                <span>{{ category.name }}</span>
              </div>
            </a-option>
          </a-select>

          <div v-if="checkedID !== ''" class="mt-6">
            <a-form :model="{}" layout="vertical">
              <div v-for="metas in categoryFullInfo" :key="metas.id">
                <div v-if="metas.id === checkedID">
                  <a-form-item
                    v-for="meta in metas.metaData!!"
                    :key="meta.id"
                    :label="meta.key"
                    :field="meta.id"
                    :tooltip="meta.description"
                  >
                    <a-input
                      v-if="meta.type === 'number'"
                      v-model="meta.value"
                      :max="100"
                      :min="0"
                      :step="1"
                      class="w-full"
                    />
                    <a-switch
                      v-else-if="meta.type === 'switch'"
                      v-model="meta.value"
                      class="mr-2"
                    />
                    <a-input
                      v-else-if="meta.type === 'string'"
                      v-model="meta.value"
                      class="w-full"
                      :placeholder="meta.description"
                    />
                    <a-textarea
                      v-else-if="meta.type === 'textarea'"
                      v-model="meta.value"
                      class="w-full"
                      :placeholder="meta.description"
                    />
                    <a-select
                      v-else-if="meta.type === 'select' && meta.defaultValue"
                      v-model="meta.value"
                      class="w-full"
                      :placeholder="meta.description"
                    >
                      <template v-if="meta.defaultValue">
                        <a-option v-for="item in meta.defaultValue.split(',')" :key="item">
                          {{ item }}
                        </a-option>
                      </template>
                    </a-select>
                    <RichTextEditor
                      v-else-if="meta.type === 'richText'"
                      v-model="meta.value"
                      min-height="200px"
                      class="w-full"
                    />
                  </a-form-item>
                </div>
              </div>
            </a-form>
          </div>
        </div>
      </a-card>

      <a-card class="mb-6" :bordered="true">
        <template #title>
          <div class="flex items-center text-[var(--color-text-1)]">
            <icon-file class="mr-2 text-[var(--color-primary-6)]" />
            上传种子文件
          </div>
        </template>
        <div class="p-4">
          <a-upload
            :auto-upload="false"
            :file-list="fileList"
            class="w-full"
            draggable
            @change="handleFileChange"
          >
            <template #upload-button>
              <div class="h-48 w-full flex flex-col items-center justify-center border-2 border-[var(--color-border)] rounded-sm border-dashed transition-colors hover:border-[var(--color-primary-6)]">
                <icon-upload class="mb-4 text-32px text-[var(--color-text-3)]" />
                <div class="text-16px text-[var(--color-text-3)]">
                  点击或拖拽文件到此处上传
                </div>
                <div class="mt-2 text-14px text-[var(--color-text-4)]">
                  支持 .torrent 文件
                </div>
              </div>
            </template>
          </a-upload>
        </div>
      </a-card>

      <a-card class="mb-6" :bordered="true">
        <template #title>
          <div class="flex items-center text-[var(--color-text-1)]">
            <icon-message class="mr-2 text-[var(--color-primary-6)]" />
            评论信息
            <a-tag class="ml-2" color="gray">
              选填
            </a-tag>
          </div>
        </template>
        <div class="p-4">
          <a-form-item field="comment" label="Comment">
            <a-textarea
              v-model="torrentComment"
              placeholder="请输入评论信息（选填）"
              :max-length="500"
              show-word-limit
            />
          </a-form-item>
        </div>
      </a-card>

      <div class="flex justify-end">
        <a-button type="primary" size="large" class="!px-8" @click="sendFile">
          <template #icon>
            <icon-check />
          </template>
          发布种子
        </a-button>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
:deep(.arco-card) {
  transition: all 0.3s ease;
  max-width: 100%;

  &:hover {
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  }
}

:deep(.arco-upload) {
  .arco-upload-list {
    margin-top: 16px;
  }

  .arco-upload-trigger {
    width: 100%;
    max-width: 100%;
  }

  .arco-upload-drag {
    width: 100%;
    max-width: 100%;
  }
}

:deep(.arco-radio-button) {
  border-radius: 4px;
  margin-right: 8px;
  margin-bottom: 8px;
}
</style>
