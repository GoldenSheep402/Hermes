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
          field="email"
          :rules="[{ required: true, message: $t('login.form.userName.errMsg') }]"
          :validate-trigger="['change', 'blur']"
          hide-label
      >
        <a-input
            v-model="userInfo.email"
        >
          <template #prefix>
            <icon-user />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item
          field="password"
          :rules="[{ required: true, message: '注册出错，轻刷新重试' }]"
          :validate-trigger="['change', 'blur']"
          hide-label
      >
        <a-input-password
            v-model="userInfo.password"
            allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>
      <a-space :size="16" direction="vertical">
        <div class="login-form-password-actions">
        </div>
        <a-button type="primary" html-type="submit" long class="login-form-register-btn">
          {{ $t('login.form.register') }}
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import useLoading from '@/hooks/loading';
import {AuthService} from "@/services/grpc.ts";
import {RegisterRequest} from "@/lib/proto/gen/proto/auth/v1/auth_service.pb.ts";
import {Message} from "@arco-design/web-vue";
import { useRouter } from 'vue-router';
const errorMessage = ref('');
const { loading, setLoading } = useLoading();

const router = useRouter();


const userInfo = reactive({
  email: '',
  password: '',
});


const reg = () => {
  console.log({userInfo})
  const req = ref<RegisterRequest>({
    email: userInfo.email,
    name: userInfo.email,
    password: userInfo.password,
  });
  AuthService.Register(req.value).then((res) => {
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
