<script setup lang="ts">
import { onMounted, ref } from "vue";
import { UserService } from "@/services/grpc";
import USERINFO from "@/router/routes/modules/userInfo";

interface userInfo {
  name: string;
  email: string;
  download: number;
  upload: number;
  torrentPublished: number;
  torrentDownload: number;
  torrentSeeding: number;
  key: string;
}

const mockData = ref<userInfo>({
  name: "admin",
  email: "admin@admin.com",
  download: 0,
  upload: 0,
  torrentPublished: 0,
  torrentDownload: 0,
  torrentSeeding: 0,
  key: "111"
})

function toDescriptionData(_data: userInfo) {
  return [
    {
      label: "用户名",
      value: _data.name,
    },
    {
      label: "邮箱",
      value: _data.email,
    },
    {
      label: "下载量",
      value: _data.download,
    },
    {
      label: "上传量",
      value: _data.upload,
    },
    {
      label: "发布种子",
      value: _data.torrentPublished,
    },
    {
      label: "下载种子",
      value: _data.torrentDownload,
    },
    {
      label: "做种种子",
      value: _data.torrentSeeding,
    },
    {
      label: "密钥",
      value: _data.key,
    }
  ]
}

const userInfo = ref<userInfo>({} as userInfo);
async function fetchUserInfo() { 
  UserService.GetUserInfo({}).then(async (res) => {
    userInfo.value.email = res.email!!;
    userInfo.value.name = res.name!!;
    userInfo.value.download = res.download!!;
    userInfo.value.upload = res.upload!!;
    userInfo.value.torrentPublished = res.torrentPublished!!;
    userInfo.value.torrentDownload = res.torrentDownloaded!!;
    userInfo.value.torrentSeeding = res.torrentSeeding!!;
  }).catch((err)=>{
    console.error('Error fetching user info:', err);
    userInfo.value = mockData.value;
  }).finally(async() => {
    UserService.GetUserPassKey({}).then((res) => {
      userInfo.value.key = res.passKey!!;
    }).catch((err) => {
      console.error('Error fetching user passkey:', err);
    });
  });
}

onMounted(() => {
  fetchUserInfo();
});
</script>

<template>
  <div class="p-5 h-full">
    <div class="p-5 bg-[--color-bg-2] h-full">
      <div class="p-0.5 text-30px leading-[1.4] font-500 text-[--color-text-1] mb-5">
        个人信息
      </div>

      <div>
        <a-descriptions title="用户详细信息" :data="toDescriptionData(userInfo)" :column="1" bordered>
        </a-descriptions>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">

</style>