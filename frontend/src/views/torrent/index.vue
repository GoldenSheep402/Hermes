<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {TorrentService, TrackerService, UserService} from "@/services/grpc.ts";
import {Notification} from "@arco-design/web-vue";

interface TorrentMessage {
  id: string;
  name: string;
  description: string;
  categoryId: string;
  categoryName: string;
  downloading: number;
  seeding: number;
  finished: number;
}

const torrentList = ref<TorrentMessage[]>([]);

async function fetchTorrentData() {
  TorrentService.GetTorrentV1List({}).then(async (res) => {
    for (let i = 0; i < res.torrents!!.length; i++) {
      const downloadingCount = ref<number>(0);
      const seedingCount = ref<number>(0);
      const finishedCount = ref<number>(0);
      await TrackerService.GetTorrentDownloadingStatus({torrentId: res.torrents!![i].id})
          .then((statusRes) => {
            downloadingCount.value = statusRes.downloading!!;
            seedingCount.value = statusRes.seeding!!;
            finishedCount.value = statusRes.finished!!;
          })
          .catch((err) => {
            console.error('Failed to get torrent status', err);
          });

      torrentList.value.push({
        id: res.torrents!![i].id!!,
        name: res.torrents!![i].name!!,
        description: res.torrents!![i].description!!,
        categoryId: res.torrents!![i].categoryId!!,
        categoryName: res.torrents!![i].categoryName!!,
        downloading: downloadingCount.value,
        seeding: seedingCount.value,
        finished: finishedCount.value,
      });
    }
  }).catch((err) => {
    console.error('Failed to fetch torrent list', err);
  });
}

function base64ToUint8Array(base64: string): Uint8Array {
  const binaryString = window.atob(base64);
  const len = binaryString.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes;
}

function downloadTorrent(id: string, name: string) {
  TorrentService.DownloadTorrentV1({id: id}).then((res) => {
    if (res.data) {
      const uint8Array = base64ToUint8Array(res.data);
      const blob = new Blob([uint8Array], {type: "application/octet-stream"});
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${name}.torrent`;
      document.body.appendChild(link);
      link.click();

      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    }
    handleNotification("success", "下载成功", "种子下载成功");
  }).catch((err) => {
    handleNotification("error", "下载失败", "种子下载失败");
    console.error('Error downloading torrent:', err);
  });
}


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


const passkey = ref<string>("");

function getPasskey() {
  UserService.GetUserPassKey({}).then((res) => {
    passkey.value = res.passKey!!;
  }).catch(() => {
    console.log("Get Passkey Fail");
  });
}

function genUrl(id: string) {
  const baseUrl = import.meta.env.VITE_GAPI_URL;
  return `${baseUrl}/api/torrent/download/${passkey.value}?id=${id}`;
}

onMounted(() => {
  getPasskey();
  fetchTorrentData();
  console.log(torrentList.value);
});
</script>

<template>
  <div class="p-5">
    <div class="p-5 bg-[--color-bg-2]">
      <div class="p-0.5 text-20px leading-[1.4] font-500 text-[--color-text-1] mb-5">
        种子列表
      </div>

      <a-table :data="torrentList">
        <template #columns>
          <a-table-column key="name" dataIndex="name" title="名称"></a-table-column>
          <a-table-column key="categoryName" dataIndex="categoryName" title="类别名称"></a-table-column>
          <!--          <a-table-column key="finished" dataIndex="finished" title="下载"></a-table-column>-->
          <a-table-column key="status" title="状态">
            <template #cell="{record}">
              <div class="flex flex-row gap-2">
                <a-statistic :value="record.seeding">
                  <template #suffix>
                    <icon-arrow-up/>
                  </template>
                </a-statistic>

                <a-statistic :value="record.finished">
                  <template #suffix>
                    <icon-check/>
                  </template>
                </a-statistic>
              </div>

            </template>
          </a-table-column>
          <a-table-column key="action" title="操作">
            <template #cell="{record}">
              <a-button type="primary"
                        :href="genUrl(record.id)"
                        @click.prevent="downloadTorrent(record.id,record.name)">
                下载
              </a-button>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>
  </div>
</template>

<style lang="less" scoped>

</style>