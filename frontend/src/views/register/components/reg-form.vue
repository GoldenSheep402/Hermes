<template>
  <div class="login-form-wrapper">
    <div class="login-form-title">HERMES</div>
    <div class="login-form-sub-title">注册账户</div>
    <div class="login-form-error-msg">{{ errorMessage }}</div>
    <a-form
        ref="regForm"
        :model="userInfo"
        class="reg-form"
        layout="vertical"
        @submit="reg"
    >
      <a-form-item
          :rules="[{ required: true, message: $t('login.form.userName.errMsg') }]"
          :validate-trigger="['change', 'blur']"
          field="email"
          hide-label
      >
        <a-input
            v-model="userInfo.email"
        >
          <template #prefix>
            <icon-user/>
          </template>
        </a-input>
      </a-form-item>
      <a-form-item
          :rules="[{ required: true, message: '注册出错，轻刷新重试' }]"
          :validate-trigger="['change', 'blur']"
          field="password"
          hide-label
      >
        <a-input-password
            v-model="userInfo.password"
            allow-clear
        >
          <template #prefix>
            <icon-lock/>
          </template>
        </a-input-password>
      </a-form-item>

      <a-form-item>
        <div class="flex flex-row justify-between w-full gap-2">
          <a-input v-model="code"></a-input>
          <a-button type="primary" @click="sendEmail">获取验证码</a-button>
        </div>
      </a-form-item>
      <a-space :size="16" direction="vertical">
        <div class="login-form-password-actions">
        </div>
        <a-button class="login-form-register-btn" html-type="submit" long type="primary">
          {{ $t('login.form.register') }}
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import {reactive, ref} from 'vue';
import {AuthService} from "@/services/grpc.ts";
import {Message} from "@arco-design/web-vue";
import {useRouter} from 'vue-router';
import {RegisterWithEmailRequest} from "@/lib/proto/auth/v1/auth.pb.ts";

const errorMessage = ref('');

const router = useRouter();

const code = ref<string>("");

const userInfo = reactive({
  email: '',
  username: '',
  password: '',
});

function sendEmail() {
  AuthService.RegisterSendEmail({email: userInfo.email}).then((res) => {
    Message.success('验证码已发送');
  }).catch((err) => {
    Message.error(err.message);
  });
}


const reg = () => {
  const req = ref<RegisterWithEmailRequest>({
    email: userInfo.email,
    username: userInfo.email,
    password: userInfo.password,
    emailToken: code.value,
  });
  AuthService.RegisterWithEmail(req.value).then((res) => {
    Message.success('注册成功');
  }).catch((err) => {
    Message.error(err.message);
    return
  });
  router.push({
    name: 'login',
  });
};
</script>

<style lang="less" scoped>
.login-form {
  &-wrapper {
    width: 320px;
  }

  &-title {
    color: var(--color-text-1);
    font-weight: 500;
    font-size: 24px;
    line-height: 32px;
  }

  &-sub-title {
    color: var(--color-text-3);
    font-size: 16px;
    line-height: 24px;
  }

  &-error-msg {
    height: 32px;
    color: rgb(var(--red-6));
    line-height: 32px;
  }

  &-password-actions {
    display: flex;
    justify-content: space-between;
  }

  &-register-btn {
    //color: var(--color-text-3) !important;
  }
}
</style>
