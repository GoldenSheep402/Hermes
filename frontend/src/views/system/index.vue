<script lang="ts" setup>
import {onMounted, ref, UnwrapRef} from "vue";
import {SystemService} from "@/services/grpc.ts";
import {InnetTracker, Settings} from "@/lib/proto/system/v1/system.pb.ts";
import {Notification} from "@arco-design/web-vue";

const setting = ref<Settings>({
  peerExpireTime: 0,
  allowedNets: [],
  smtpEnable: false,
  smtpHost: '',
  smtpPort: 0,
  smtpUser: '',
  smtpPassword: '',
  registerEnable: false,
  loginEnable: false,
  publishEnable: false,
});

const trackers = ref<InnetTracker[]>([]);

function fetchSystemSettings() {
  setting.value = {}
  trackers.value = []
  SystemService.GetSettings({}).then((res) => {
    if (res.settings) {
      setting.value = res.settings;
      trackers.value = res.settings.innetTracker!!;
    }
  }).catch((err) => {
    console.error(err);
  });
}

function handleSubmit() {
  SystemService.SetSettings({settings: setting.value}).then((res) => {
    console.log(res);
  }).catch((err) => {
    handleNotification('error', '修改失败', '系统设置修改失败');
    console.error(err);
  }).finally(() => {
    handleNotification('success', '修改成功', '系统设置已修改');
    fetchSystemSettings();
  });
}

function deleteTracker(index: number) {
  trackers.value[index].toBeDeleted = !trackers.value[index].toBeDeleted;
}

function setButtonStatus(isDanger: UnwrapRef<InnetTracker["toBeDeleted"]> | undefined) {
  return !isDanger ? 'primary' : 'secondary'
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


onMounted(() => {
  fetchSystemSettings();
});
</script>

<template>
  <div class="p-5 h-full">
    <div class="p-5 bg-[--color-bg-2] h-full">
      <div class="p-0.5 text-20px leading-[1.4] font-500 text-[--color-text-1] mb-5">
        系统设置
      </div>
      <a-form :model="setting" @submit="handleSubmit">
        <!-- peerExpireTime -->
        <a-form-item field="peerExpireTime" label="Peer过期时间">
          <a-input-number v-model="setting.peerExpireTime"/>
        </a-form-item>

        <!-- allowedNets -->
        <!--        <subnets v-model:subnets="setting.allowedNets!!"/>-->

        <!-- smtpEnable -->
        <a-form-item field="smtpEnable" label="SMTP启用">
          <a-switch v-model="setting.smtpEnable"/>
        </a-form-item>

        <!-- smtpHost -->
        <a-form-item field="smtpHost" label="SMTP主机">
          <a-input v-model="setting.smtpHost"/>
        </a-form-item>

        <!-- smtpPort -->
        <a-form-item field="smtpPort" label="SMTP端口">
          <a-input-number v-model="setting.smtpPort"/>
        </a-form-item>

        <!-- smtpUser -->
        <a-form-item field="smtpUser" label="SMTP用户名">
          <a-input v-model="setting.smtpUser"/>
        </a-form-item>

        <!-- smtpPassword -->
        <a-form-item field="smtpPassword" label="SMTP密码">
          <a-input-password v-model="setting.smtpPassword"/>
        </a-form-item>

        <!-- registerEnable -->
        <a-form-item field="registerEnable" label="注册启用">
          <a-switch v-model="setting.registerEnable"/>
        </a-form-item>

        <!-- loginEnable -->
        <a-form-item field="loginEnable" label="登录启用">
          <a-switch v-model="setting.loginEnable"/>
        </a-form-item>

        <!-- publishEnable -->
        <a-form-item field="publishEnable" label="发布启用">
          <a-switch v-model="setting.publishEnable"/>
        </a-form-item>

        <!-- trackers -->

        <a-form-item field="trackers" label="内网跟踪器">
          <div class="flex flex-col">
            <a-button type="primary" @click="trackers.push({addr: '', enable: false})">添加</a-button>
            <div v-for="(tracker, index) in trackers" :key="index">
              <div class="flex flex-row items-center gap-2 m-t2">
                <a-input v-model="trackers[index].addr"/>
                <a-button :type="setButtonStatus(trackers[index].toBeDeleted)" @click="deleteTracker(index)">删除
                </a-button>
                <a-switch v-model="trackers[index].enable" class="justify-center flex"/>
              </div>
            </div>
          </div>
        </a-form-item>

        <a-form-item>
          <a-button html-type="submit">修改</a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<style lang="less" scoped>

</style>
