<template>
  <div class="register-form-wrapper">
    <div class="register-form-title">{{ $t('site.maintitle') }}</div>
    <div class="register-form-sub-title">{{ $t('register.subtitle') }}</div>
    <div class="register-form-error-msg">{{ errorMessage }}</div>
    <a-form :model="registerForm" class="reg-form" layout="vertical" @submit="handleSubmit">
      <a-form-item :rules="[
        { required: true, message: $t('register.form.email.required') },
        {
          validator: (value, cb) => {
            if (EMAIL_REGEX.test(value)) {
              cb();
            } else {
              cb($t('register.form.email.invalid'));
            }
          }
        }
      ]" :validate-trigger="['change', 'blur']" field="email" hide-label feedback>
        <a-input v-model="registerForm.email" :placeholder="$t('register.form.email.placeholder')">
          <template #prefix>
            <icon-email />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item :rules="[{ required: true, message: $t('register.form.password.required') }, {
        minLength: PASSWORD_MIN, message: $t('register.form.password.min', { minLength: PASSWORD_MIN })
      }]" :validate-trigger="['change', 'blur']" field="password" hide-label>
        <a-input-password v-model="registerForm.password" allow-clear
          :placeholder="$t('register.form.password.placeholder', { minLength: PASSWORD_MIN })">
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <a-form-item>
        <div class="flex flex-row justify-between w-full gap-2">
          <a-input v-model="registerForm.emailToken" :placeholder="$t('register.form.emailToken.placeholder')">
            <template #prefix>
              <icon-safe />
            </template>
          </a-input>
          <a-button type="primary" @click="sendEmail">{{ $t('register.form.emailToken.send') }}</a-button>
        </div>
      </a-form-item>
      <a-space :size="16" direction="vertical">
        <div class="register-form-password-actions">
        </div>
        <a-button class="register-form-register-btn" html-type="submit" long type="primary" :loading="loading">
          {{ $t('register.form.sumbit') }}
        </a-button>
        <a-button type="text" long class="register-form-login-btn" @click="handleLogin">
          {{ $t('register.form.login') }}
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { AuthService } from "@/services/grpc.ts";
import { Message, ValidatedError } from "@arco-design/web-vue";
import { useRouter } from 'vue-router';
import { RegisterWithEmailRequest } from "@/lib/proto/auth/v1/auth.pb.ts";
import { EMAIL_REGEX, PASSWORD_MIN } from '@/utils/constants';
import useLoading from '@/hooks/loading';
import { set } from 'nprogress';

const errorMessage = ref('');

const router = useRouter();
const { loading, setLoading } = useLoading();

const registerForm = reactive({
  email: '',
  username: '',
  password: '',
  emailToken: '',
});

function sendEmail() {
  if (!EMAIL_REGEX.test(registerForm.email)) {
    Message.error('邮箱格式不正确');
    return;
  }
  AuthService.RegisterSendEmail({ email: registerForm.email }).then((res) => {
    Message.success('验证码已发送');
  }).catch((err) => {
    Message.error(err.message);
  });
}


const handleSubmit = ({ values, errors }: { values: Record<string, any>; errors: Record<string, ValidatedError> | undefined }, ev: Event) => {
  if (loading.value) return;

  if (errors !== undefined) {
    return;
  }

  setLoading(true);

  const req = ref<RegisterWithEmailRequest>({
    email: registerForm.email,
    username: registerForm.email,
    password: registerForm.password,
    emailToken: registerForm.emailToken,
  });

  AuthService.RegisterWithEmail(req.value).then((res) => {
    setLoading(false);
    Message.success('注册成功');
    router.push({
      name: 'login',
    });
  }).catch((err) => {
    setLoading(false);
    Message.error(err.message);
    return
  });
};

const handleLogin = () => {
  router.push({
    name: 'login',
  });
};

</script>

<style lang="less" scoped>
.register-form {
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

  &-login-btn {
    color: var(--color-text-3) !important;
  }
}
</style>
