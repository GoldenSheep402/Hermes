<template>
  <div class="w-full flex h-dvh">
    <div class="logo">
      <div class="logo-text">HERMES</div>
    </div>
<!--    <LoginBanner />-->
    <div class="relative flex flex-1 items-center justify-center">
      <div class="content-inner">
        <LoginForm />
      </div>
      <div class="footer">
        <Footer />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import Footer from '@/components/footer/index.vue';
  import LoginForm from './components/login-form.vue';
</script>

<style lang="less" scoped>
  .footer {
    position: absolute;
    right: 0;
    bottom: 0;
    width: 100%;
  }

  .logo {
    position: fixed;
    top: 24px;
    left: 22px;
    z-index: 1;
    display: inline-flex;
    align-items: center;

    &-text {
      margin-right: 4px;
      margin-left: 4px;
      color: var(--color-text-1);
      font-size: 20px;
    }
  }
</style>
