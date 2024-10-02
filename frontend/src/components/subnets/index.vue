<script setup lang="ts">
import { defineModel, ref, watch } from 'vue';

const subnets = defineModel<string[]>(
    'subnets',
    { required: true, default: () => [] }
);

const cidrString = ref(subnets.value.join(', '));

watch(subnets, (newSubnets) => {
  const newString = newSubnets.join(', ');
  if (newString !== cidrString.value) {
    cidrString.value = newString;
  }
});

watch(cidrString, (newString) => {
  const newArray = newString.split(',')
      .map(item => item.trim())
      .filter(item => item !== '');

  if (newArray.join(', ') !== subnets.value.join(', ')) {
    subnets.value = newArray;
  }
});
</script>

<template>
  <a-form-item label="允许的网络">
    <a-input v-model="cidrString" placeholder="输入CIDR，用逗号分隔" />
  </a-form-item>
</template>

<style scoped lang="less">
</style>
