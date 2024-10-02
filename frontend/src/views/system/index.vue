<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {SystemService} from "@/services/grpc.ts";
import Subnets from "@/components/subnets/index.vue";
import {Settings} from "@/lib/proto/system/v1/system.pb.ts";

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

function fetchSystemSettings() {
  SystemService.GetSettings({}).then((res) => {
    if (res.settings) {
      setting.value = res.settings;
      console.log(setting.value);
    }
  }).catch((err) => {
    console.error(err);
  });
}

function handleSubmit() {
  SystemService.SetSettings({settings: setting.value}).then((res) => {
    console.log(res);
  }).catch((err) => {
    console.error(err);
  });
}

onMounted(() => {
  fetchSystemSettings();
});
</script>

<template>
  <div class="p-5 h-full">
    <div class="p-5 bg-[--color-bg-2] h-full">
      <a-form :model="setting" @submit="handleSubmit">
        <!-- peerExpireTime -->
        <a-form-item field="peerExpireTime" label="Peer过期时间">
          <a-input-number v-model="setting.peerExpireTime"/>
        </a-form-item>

        <!-- allowedNets -->
        <subnets v-model:subnets="setting.allowedNets!!"/>

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

        <a-form-item>
          <a-button html-type="submit">修改</a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<style lang="less" scoped>

</style>
