<script setup lang="ts">
import type { CreateCategoryRequest } from '@/lib/proto/category/v1/category.pb.ts'
import RichTextEditor from '@/components/Editor/index.vue'
import { CategoryService } from '@/services/grpc.ts'
import { Modal, Notification } from '@arco-design/web-vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

interface Category {
  id: string
  name: string
  description: string
}

interface CategoryMetadata {
  order: number
  type: string
  key: string
  description: string
  DefaultValue: string
  richTextValue?: string
}

const categoryNew = ref<Category>({
  id: '',
  name: '',
  description: '',
})

const categoryMetadatas = ref<CategoryMetadata[]>([])
const categoryList = ref<Category[]>([])
const activeKey = ref<number>(1)
const showAddCategory = ref<boolean>(false)
const tagDefaultValue = ref<string[]>([])

const router = useRouter()

function fetchCategoryList() {
  categoryList.value = []
  CategoryService.GetCategoryList({}).then((res) => {
    for (let i = 0; i < res.category!.length; i++) {
      categoryList.value.push({
        id: res.category![i].id!,
        name: res.category![i].name!,
        description: res.category![i].description!,
      })
    }
  })
}

function addCategory() {
  showAddCategory.value = true
}

function handleAdd() {
  categoryMetadatas.value.push({
    order: categoryMetadatas.value.length + 1,
    type: 'string',
    key: '默认key',
    description: '类别名称',
    DefaultValue: '',
  })
}

function handleDelete(id: string) {
  Modal.warning({
    title: '确认删除',
    content: '确定要删除这个类别吗？删除后无法恢复。',
    okText: '确认',
    cancelText: '取消',
    onOk: () => {
      CategoryService.DeleteCategory({ id }).then(() => {
        handleNotification('success', '成功', '删除类别成功')
        fetchCategoryList()
      }).catch((err) => {
        console.error('删除类别失败:', err)
        handleNotification('error', '失败', '删除类别失败')
      })
    },
  })
}

function clear() {
  categoryMetadatas.value = [{
    order: 1,
    type: 'string',
    key: '默认key',
    description: '类别名称',
    DefaultValue: '',
  }]
  activeKey.value = 1
}

function checkEmptyInput() {
  if (categoryNew.value.name === '') {
    handleNotification('warning', '警告', '类别名称不能为空')
    return false
  }
  return true
}

function createCategory() {
  const req = ref<CreateCategoryRequest>({ category: { metaData: [] } } as CreateCategoryRequest)

  req.value.category!.name = categoryNew.value.name
  req.value.category!.description = categoryNew.value.description

  for (let i = 0; i < categoryMetadatas.value.length; i++) {
    if (categoryMetadatas.value[i].type === 'select') {
      req.value.category!.metaData!.push({
        key: categoryMetadatas.value[i].key,
        order: categoryMetadatas.value[i].order,
        type: categoryMetadatas.value[i].type,
        description: categoryMetadatas.value[i].description,
        defaultValue: tagDefaultValue.value.join(','),
      })
    }
    else if (categoryMetadatas.value[i].type === 'richText') {
      req.value.category!.metaData!.push({
        key: categoryMetadatas.value[i].key,
        order: categoryMetadatas.value[i].order,
        type: categoryMetadatas.value[i].type,
        description: categoryMetadatas.value[i].description,
        defaultValue: categoryMetadatas.value[i].richTextValue || '',
      })
    }
    else {
      req.value.category!.metaData!.push({
        key: categoryMetadatas.value[i].key,
        order: categoryMetadatas.value[i].order,
        type: categoryMetadatas.value[i].type,
        description: categoryMetadatas.value[i].description,
        defaultValue: categoryMetadatas.value[i].DefaultValue,
      })
    }
  }

  CategoryService.CreateCategory(req.value).then(() => {
    handleNotification('success', '成功', '添加类别成功')
  }).catch(() => {
    handleNotification('error', '失败', '添加类别失败')
  }).finally(() => {
    showAddCategory.value = false
    fetchCategoryList()
  })
}

function viewDetail(record: Category) {
  router.push({
    name: 'CategoryDetail',
    params: { id: record.id },
  })
}

categoryMetadatas.value.push({
  order: 1,
  type: 'string',
  key: '默认key',
  description: '类别名称',
  DefaultValue: '',
})

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
      <div class="flex justify-between">
        <div class="mb-5 p-0.5 text-20px text-[--color-text-1] font-500 leading-[1.4]">
          类别列表
        </div>

        <div>
          <a-button type="primary" @click="addCategory()">
            添加类别
          </a-button>
        </div>
      </div>

      <a-table :data="categoryList">
        <template #columns>
          <a-table-column title="名称" data-index="name" :width="200" />
          <a-table-column title="描述" data-index="description" :width="300" />
          <a-table-column title="操作" align="center" :width="150">
            <template #cell="{ record }">
              <div class="w-full flex justify-center">
                <div class="w-fit flex flex-col items-center gap-2 md:flex-row">
                  <a-button @click="viewDetail(record)">
                    详情
                  </a-button>
                  <a-button type="primary" status="danger" @click="handleDelete(record.id)">
                    删除
                  </a-button>
                </div>
              </div>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <a-modal v-model:visible="showAddCategory" width="800px" :on-before-ok="checkEmptyInput" @ok="createCategory()" @close="clear()">
        <template #title>
          添加类别
        </template>
        <div class="p-5">
          <div class="flex flex-row gap-5">
            <a-input v-model="categoryNew.name" placeholder="类别名称" />
            <a-input v-model="categoryNew.description" placeholder="类别描述" />
          </div>

          <div class="mt-4">
            <a-tabs
              v-model:active-key="activeKey" :editable="true" type="card-gutter" show-add-button auto-switch @add="handleAdd"
              @delete="handleDelete"
            >
              <a-tab-pane v-for="meta in categoryMetadatas" :key="meta.order" :title="(meta.order).toString()">
                <a-form :model="meta" class="p-5">
                  <a-form-item field="order" label="序号">
                    <a-input-number v-model="meta.order" />
                  </a-form-item>
                  <a-form-item field="type" label="类型">
                    <a-radio-group v-model="meta.type" type="button">
                      <a-radio value="string">
                        字符串
                      </a-radio>
                      <a-radio value="number">
                        数字
                      </a-radio>
                      <a-radio value="switch">
                        开关
                      </a-radio>
                      <a-radio value="textarea">
                        多行文本
                      </a-radio>
                      <a-radio value="select">
                        下拉框
                      </a-radio>
                      <a-radio value="richText">
                        富文本
                      </a-radio>
                    </a-radio-group>
                  </a-form-item>
                  <a-form-item field="key" label="键">
                    <a-input v-model="meta.key" />
                  </a-form-item>
                  <a-form-item field="description" label="描述">
                    <a-input v-model="meta.description" />
                  </a-form-item>
                  <a-form-item v-if="meta.type !== 'select'" field="DefaultValue" label="默认值">
                    <a-input v-model="meta.DefaultValue" />
                  </a-form-item>
                  <a-form-item v-if="meta.type === 'select'" field="values" label="多个值">
                    <a-input-tag v-model:model-value="tagDefaultValue" allow-clear />
                  </a-form-item>
                  <a-form-item v-if="meta.type === 'richText'" field="richTextValue" label="富文本内容">
                    <RichTextEditor v-model="meta.richTextValue" min-height="200px" />
                  </a-form-item>
                </a-form>
              </a-tab-pane>
            </a-tabs>
          </div>
        </div>
      </a-modal>
    </div>
  </div>
</template>

<style scoped lang="less"></style>
