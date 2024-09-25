<script setup lang="ts">
import {CategoryService} from "@/services/grpc.ts";
import {onMounted, ref} from "vue";

interface Category {
  id: number;
  name: string;
  description: string;
}

const categoryList = ref<Category[]>([]);

function fetchCategoryList() {
  CategoryService.GetCategoryList({}).then((res) => {
    console.log(res);
  });
}

onMounted(() => {
  fetchCategoryList()
});
</script>

<template>
  <div class="p-5">
    <div class="p-5 bg-[--color-bg-2]">
    </div>
  </div>
</template>

<style scoped lang="less">
</style>