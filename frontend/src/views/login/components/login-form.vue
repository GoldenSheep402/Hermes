<template>
  <div class="login-form-wrapper">
    <div class="login-form-title">{{ $t('site.maintitle') }}</div>
    <div class="login-form-sub-title">{{ $t('login.subtitile') }}</div>
    <div class="login-form-error-msg">{{ errorMessage }}</div>
    <a-form :model="loginForm" class="login-form" layout="vertical" @submit="handleSubmit">
      <a-form-item field="email" :rules="[
        { required: true, message: $t('login.form.email.required') },
        {
          validator: (value, cb) => {
            if (EMAIL_REGEX.test(value)) {
              cb();
            } else {
              cb($t('login.form.email.invalid'));
            }
          }
        }
      ]" :validate-trigger="['change', 'blur']" hide-label feedback>
        <a-input v-model="loginForm.email" :placeholder="$t('login.form.email.placeholder')">
          <template #prefix>
            <icon-email />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item field="password" :rules="[{ required: true, message: $t('login.form.password.required') }]"
        :validate-trigger="['change', 'blur']" hide-label>
        <a-input-password v-model="loginForm.password" :placeholder="$t('login.form.password.placeholder')" allow-clear>
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>
      <a-space :size="16" direction="vertical">
        <div class="login-form-password-actions">
          <a-checkbox checked="rememberPassword" :model-value="loginConfig.rememberPassword"
            @change="setRememberPassword as any">
            {{ $t('login.form.rememberPassword') }}
          </a-checkbox>
          <a-link>{{ $t('login.form.forgetPassword') }}</a-link>
        </div>
        <a-button type="primary" html-type="submit" long :loading="loading">
          {{ $t('login.form.submit') }}
        </a-button>
        <a-button type="text" long class="login-form-register-btn" @click="handleReg">
          {{ $t('login.form.register') }}
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { ValidatedError } from '@arco-design/web-vue/es/form/interface';
import { useI18n } from 'vue-i18n';
import { useStorage } from '@vueuse/core';
import { useUserStore } from '@/store';
import useLoading from '@/hooks/loading';
import type { LoginData } from '@/api/user';
import { EMAIL_REGEX } from '@/utils/constants';

const router = useRouter();
const { t } = useI18n();
const errorMessage = ref('');
const { loading, setLoading } = useLoading();
const userStore = useUserStore();

const loginConfig = useStorage('login-config', {
  rememberPassword: true,
  email: '',
  password: '',
});

const loginForm = reactive({
  email: loginConfig.value.email,
  password: loginConfig.value.password,
});

const handleSubmit = async ({ errors, values, }: {
  errors: Record<string, ValidatedError> | undefined;
  values: Record<string, any>;
}
) => {
  if (loading.value) return;
  if (!errors) {
    setLoading(true);
    try {
      await userStore.login(values as LoginData);
      const { redirect, ...othersQuery } = router.currentRoute.value.query;
      // const validRedirects = ["ProjectDetail", "ProductDetail", "UserDetail"];
      // if (validRedirects.includes(redirect?.toString() || '')) {
      await router.push({ name: 'Workplace' });
      // } else {
      //   router.push({ name: (redirect as string) || 'Workplace' });
      // }
      Message.success(t('login.form.login.success'));

    } catch (err) {
      if (err instanceof Error) {
        errorMessage.value = err.message;
      } else {
        errorMessage.value = t('login.form.login.error');
      }
    } finally {
      setLoading(false);
    }
  }
};

const setRememberPassword = (value: boolean) => {
  loginConfig.value.rememberPassword = value;
};

const handleReg = () => {
  router.push({
    name: 'register',
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
    color: var(--color-text-3) !important;
  }
}
</style>
