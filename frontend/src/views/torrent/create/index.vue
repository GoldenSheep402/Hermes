<script lang="ts" setup>
import type { Category as CategoryBase } from '@/lib/proto/category/v1/category.pb.ts'
import type { CreateTorrentV1Request } from '@/lib/proto/torrent/v1/torrent.pb.ts'
import type { FileItem } from '@arco-design/web-vue'
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
const fileList = ref<FileItem[]>([])
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

function handleFileChange(fileList: FileItem[], file: FileItem) {
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
    <div class="bg-[--color-bg-2] p-5">
      <div class="mb-5 p-0.5 text-20px text-[--color-text-1] font-500 leading-[1.4]">
        发布种子
      </div>

      <div>
        <div class="mb-5 p-0.5 text-20px text-[--color-text-1] font-500 leading-[1.4]">
          选择类别
        </div>

        <div>
          <a-radio-group v-model="checkedID" class="mb-5" type="button">
            <div v-for="category in categoryList" :key="category.id">
              <a-radio :value="category.id">
                {{ category.name }}
              </a-radio>
            </div>
          </a-radio-group>

          <div v-if="checkedID !== ''">
            <div v-for="metas in categoryFullInfo" :key="metas.id">
              <div v-if="metas.id === checkedID">
                <div v-for="meta in metas.metaData!!" :key="meta.id" class="mb-4 flex items-center">
                  <div class="mr-4 flex items-center">
                    <div class="flex items-center justify-center text-[--color-text-1] font-semibold">
                      {{ meta.key }}:
                    </div>
                  </div>
                  <div class="flex-grow">
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
                    />
                    <a-textarea
                      v-else-if="meta.type === 'textarea'"
                      v-model="meta.value"
                      class="w-full"
                    />
                    <a-select
                      v-else-if="meta.type === 'select' && meta.defaultValue"
                      v-model="meta.value"
                      class="w-full"
                    >
                      <a-option v-for="item in meta.defaultValue.split(',')" :key="item">
                        {{ item }}
                      </a-option>
                    </a-select>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex flex-col space-y-4">
        <div class="flex items-center">
          <div class="w-25 text-[--color-text-1] font-semibold">
            Comment:
          </div>
          <a-input v-model="torrentComment" class="flex-grow" placeholder="请输入Comment" />
          <a-button type="primary" class="ml-5" @click="sendFile">
            发布
          </a-button>
        </div>

        <a-upload
          :auto-upload="false"
          :file-list="fileList"
          @change="handleFileChange"
        />
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>

</style>
