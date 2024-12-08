<script lang="ts" setup>
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()

const gotoPath = route.params.path as string

console.log({ gotoPath })
router.replace({ path: gotoPath })
</script>

<template>
  <div />
</template>

<style scoped lang="less"></style>
