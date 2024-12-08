<script setup lang="ts">
import { UserService } from '@/services/grpc'
import { onMounted, ref } from 'vue'

interface UserInfo {
  name: string
  email: string
  download: number
  upload: number
  torrentPublished: number
  torrentDownload: number
  torrentSeeding: number
  key: string
}

const mockData = ref<UserInfo>({
  name: 'admin',
  email: 'admin@admin.com',
  download: 0,
  upload: 0,
  torrentPublished: 0,
  torrentDownload: 0,
  torrentSeeding: 0,
  key: '111',
})

function toDescriptionData(_data: UserInfo) {
  return [
    {
      label: '用户名',
      value: _data.name,
    },
    {
      label: '邮箱',
      value: _data.email,
    },
    {
      label: '下载量',
      value: _data.download,
    },
    {
      label: '上传量',
      value: _data.upload,
    },
    {
      label: '发布种子',
      value: _data.torrentPublished,
    },
    {
      label: '下载种子',
      value: _data.torrentDownload,
    },
    {
      label: '做种种子',
      value: _data.torrentSeeding,
    },
    {
      label: '密钥',
      value: _data.key,
    },
  ]
}

const userInfo = ref<UserInfo>({} as UserInfo)
async function fetchUserInfo() {
  UserService.GetUserInfo({}).then(async (res) => {
    userInfo.value.email = res.email!
    userInfo.value.name = res.name!
    userInfo.value.download = res.download!
    userInfo.value.upload = res.upload!
    userInfo.value.torrentPublished = res.torrentPublished!
    userInfo.value.torrentDownload = res.torrentDownloaded!
    userInfo.value.torrentSeeding = res.torrentSeeding!
  }).catch((err) => {
    console.error('Error fetching user info:', err)
    userInfo.value = mockData.value
  }).finally(async () => {
    UserService.GetUserPassKey({}).then((res) => {
      userInfo.value.key = res.passKey!
    }).catch((err) => {
      console.error('Error fetching user passkey:', err)
    })
  })
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<template>
  <div class="h-full p-5">
    <div class="h-full bg-[--color-bg-2] p-5">
      <div class="mb-5 p-0.5 text-30px text-[--color-text-1] font-500 leading-[1.4]">
        个人信息
      </div>

      <div>
        <a-descriptions title="用户详细信息" :data="toDescriptionData(userInfo) as any" :column="1" bordered />
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">

</style>
