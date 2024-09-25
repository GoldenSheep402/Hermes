<script setup lang="ts">
import {CategoryService} from "@/services/grpc.ts";
import {onMounted, ref} from "vue";
import {CreateCategoryRequest} from "@/lib/proto/category/v1/category.pb.ts";
import {Notification} from "@arco-design/web-vue";

interface Category {
  id: string;
  name: string;
  description: string;
}

interface CategoryMetadata {
  order: number;
  type: string;
  description: string;
  DefaultValue: string;
}

const categoryList = ref<Category[]>([]);

function fetchCategoryList() {
  categoryList.value = [];
  CategoryService.GetCategoryList({}).then((res) => {
    for (let i = 0; i < res.category!!.length; i++) {
      categoryList.value.push({
        id: res.category!![i].id!!,
        name: res.category!![i].name!!,
        description: res.category!![i].description!!,
      });
    }
  });
}

const showAddCategory = ref<boolean>(false);
function addCategory() {
  showAddCategory.value = true;
}

function handleAdd() {
  categoryMetadatas.value.push({
    order: categoryMetadatas.value.length + 1,
    type: "string",
    description: "类别名称",
    DefaultValue: "默认值",
  });
}

function handleDelete(order: string | number) {
  const index = categoryMetadatas.value.findIndex((meta) => meta.order === order);
  categoryMetadatas.value.splice(index, 1);
}

function createCategory() {
  const req = ref<CreateCategoryRequest>({ category: { metaData: [] } } as CreateCategoryRequest);

  req.value.category!!.name = categoryNew.value.name;
  req.value.category!!.description = categoryNew.value.description;

  for (let i = 0; i < categoryMetadatas.value.length; i++) {
    req.value.category!!.metaData!!.push({
      order: categoryMetadatas.value[i].order,
      type: categoryMetadatas.value[i].type,
      description: categoryMetadatas.value[i].description,
      defaultValue: categoryMetadatas.value[i].DefaultValue,
    });
  }

  console.log(req);

  CategoryService.CreateCategory(req.value).then((res) => {
    handleNotification("success", "成功", "添加类别成功");
  }).catch((err) => {
    handleNotification("error", "失败", "添加类别失败");
  }).finally(() => {
    showAddCategory.value = false;
    fetchCategoryList()
  });
}


const categoryNew = ref<Category>({
  id: "",
  name: "",
  description: "",
});
const categoryMetadatas = ref<CategoryMetadata[]>([]);
categoryMetadatas.value.push({
  order: 1,
  type: "string",
  description: "类别名称",
  DefaultValue: "",
});


const handleNotification = (type: string, title: string, content: string) => {
  switch (type) {
    case "success":
      Notification.success({
        title: title,
        content: content,
      });
      break;
    case "error":
      Notification.error({
        title: title,
        content: content,
      });
      break;
    case "warning":
      Notification.warning({
        title: title,
        content: content,
      });
      break;
    default:
      Notification.info({
        title: title,
        content: content,
      });
  }
}

onMounted(() => {
  fetchCategoryList()
});
</script>
<template>
  <div class="p-5">
    <div class="p-5 bg-[--color-bg-2]">
      <div class="flex justify-between">
        <div class="p-0.5 text-20px leading-[1.4] font-500 text-[--color-text-1] mb-5">
          类别列表
        </div>

        <div>
          <a-button type="primary" @click="addCategory()">添加类别</a-button>
        </div>
      </div>

      <a-table :data="categoryList">
        <template #columns>
          <a-table-column title="名称" data-index="name" :width="200"></a-table-column>
          <a-table-column title="描述" data-index="description" :width="300"></a-table-column>
          <a-table-column title="操作" align="center" :width="100">
            <template #cell="{ record }">
              <div class="w-full flex justify-center">
                <div class="w-fit flex flex-col md:flex-row items-center gap-2">
                  <a-button @click="">详情</a-button>
                </div>
              </div>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <a-modal v-model:visible="showAddCategory" @ok="createCategory()">
        <div class="p-5">
          <div class="text-20px leading-[1.4] font-500 text-[--color-text-1] mb-5">
            添加类别
          </div>
          <div class="flex flex-row gap-5">
            <a-input placeholder="类别名称" v-model="categoryNew.name"></a-input>
            <a-input placeholder="类别描述" v-model="categoryNew.description"></a-input>
          </div>

          <div class="mt-2">
            <a-tabs :editable="true" type="card-gutter" @add="handleAdd" @delete="" show-add-button auto-switch>
              <a-tab-pane v-for="meta in categoryMetadatas" :key="meta.order" :title="(meta.order).toString()">
                <a-form :model="meta">
                  <a-form-item field="order" label="顺序">
                    <a-input-number v-model="meta.order" />
                  </a-form-item>
                  <a-form-item field="type" label="类型">
                    <a-radio-group type="button" v-model="meta.type">
                      <a-radio value="string">字符串</a-radio>
                      <a-radio value="number">数字</a-radio>
                      <a-radio value="switch">开关</a-radio>
                    </a-radio-group>
                  </a-form-item>
                  <a-form-item field="description" label="描述">
                    <a-input v-model="meta.description"/>
                  </a-form-item>
                  <a-form-item field="DefaultValue" label="默认值">
                    <a-input v-model="meta.DefaultValue"/>
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

<style scoped lang="less">
</style>