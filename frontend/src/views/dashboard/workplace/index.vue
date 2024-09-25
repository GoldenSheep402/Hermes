<template>
  <div class="flex p-5">
    <div class="w-full flex flex-col gap-5">
      <div class="w-full p-5 bg-[--color-bg-2]">
        <Banner/>
<!--        <div class="p-5 grid cols-[repeat(auto-fill,minmax(140px,1fr))]">-->

<!--        </div>-->
      </div>
      <div class="w-full p-5 bg-[--color-bg-2]">
        <a-space size="large">
<!--          <a-statistic title="产品数" :value="data.productNum" show-group-separator />-->
<!--          <a-statistic title="现有项目数" :value="data.projectNum" />-->
<!--          <div :hidden="isAdmin">-->
<!--            <a-statistic title="项目数限制" :value="data.projectLimit"/>-->
<!--          </div>-->
        </a-space>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Banner from './components/banner.vue';
import {onMounted, ref} from "vue";
import {useUserStore} from "@/store";

const UserStore = useUserStore();

interface data {
  productNum: number,
  projectNum: number,
  projectLimit: number,
}

const data = ref<data>({
  productNum: 0,
  projectNum: 0,
  projectLimit: <number>UserStore.projectLimit,
});

const isAdmin = ref<boolean>(false);
isAdmin.value = UserStore.role === 'admin';


// if allow normal user it will be diasbled
const refreshData = () => {
}

onMounted(() => {
  console.log(isAdmin.value)
  refreshData();
})
</script>

<script lang="ts">
export default {
  name: 'Dashboard', // If you want the include property of keep-alive to take effect, you must name the component
};
</script>

<style lang="less" scoped>

</style>

<style>
.amap-logo,.amap-copyright{
  @apply absolute -top-100%;
}
</style>
