<script setup lang="ts">
import {ref} from "vue";

interface userInfo {
  name: string;
  email: string;
  download: number;
  upload: number;
  torrentPubic: number;
  torrentDownload: number;
  torrentSeeding: number;
  key: string;
}

const mockData = ref<userInfo>({
})

function toDescriptionData(_data: userInfo) {
  return [
    {
      label: "用户名",
      content: _data.name
    },
    {
      label: "邮箱",
      content: _data.email
    },
    {
      label: "下载量",
      content: _data.download
    },
    {
      label: "上传量",
      content: _data.upload
    },
    {
      label: "发布种子",
      content: _data.torrentPubic
    },
    {
      label: "下载种子",
      content: _data.torrentDownload
    },
    {
      label: "做种种子",
      content: _data.torrentSeeding
    },
    {
      label: "密钥",
      content: _data.key
    }
  ]
}

function fetchUserInfo()
</script>

<template>
  <div class="p-5 h-full">
    <div class="p-5 bg-[--color-bg-2] h-full">
      <div class="p-0.5 text-30px leading-[1.4] font-500 text-[--color-text-1] mb-5">
        个人信息
      </div>

      <div>
        <a-descriptions title="User Info" :data="toDescriptionData(mockData)" :column="{xs:1, md:3, lg:4}" bordered>

        </a-descriptions>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">

</style>