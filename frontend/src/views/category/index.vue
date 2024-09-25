<script setup lang="ts">
import {CategoryService} from "@/services/grpc.ts";
import {onMounted, ref} from "vue";

interface Category {
  id: string;
  name: string;
  description: string;
}

const categoryList = ref<Category[]>([]);

function fetchCategoryList() {
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

onMounted(() => {
  fetchCategoryList()
});
</script>

<template>
  <div class="p-5">
    <div class="p-5 bg-[--color-bg-2]">
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
    </div>
  </div>
</template>

<style scoped lang="less">
</style>