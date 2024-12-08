<script setup lang="ts">
import type { CategoryMetaData as CategoryMetadataBase } from '@/lib/proto/category/v1/category.pb.ts'
import { CategoryService } from '@/services/grpc.ts'
import { onMounted, ref } from 'vue'

const id = defineModel<string>('id', {
  required: true,
})

const categoryMetadata = ref<CategoryMetadataBase[]>([])
function fetchCategory(id: string) {
  CategoryService.GetCategory({ id }).then((res) => {
    for (let i = 0; i < res.category!.metaData!.length; i++) {
      categoryMetadata.value.push({
        id: res.category!.metaData![i].id!,
        order: res.category!.metaData![i].order!,
        categoryId: res.category!.metaData![i].categoryId!,
        description: res.category!.metaData![i].description!,
        key: res.category!.metaData![i].key!,
        defaultValue: res.category!.metaData![i].defaultValue!,
        value: res.category!.metaData![i].value!,
      })
    }
  })
}

onMounted(() => {
  fetchCategory(id.value)
})
</script>

<template>
  <div>
    {{ categoryMetadata }}
  </div>
</template>

<style scoped lang="less">

</style>
