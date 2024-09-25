<template>
  <a-grid-item
    :class="isLast ? 'panel-col-no-border' : 'panel-col'"
    :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }"
  >
    <a-space>
      <a-avatar :size="54" class="col-avatar">
        <div style="margin: 120px"><slot /></div>
      </a-avatar>
      <a-statistic
        :title="title"
        :value="value"
        :precision="0"
        :value-from="0"
        animation
        show-group-separator
      >
      </a-statistic>
    </a-space>
  </a-grid-item>
</template>

<script lang="ts">
  export default {
    props: {
      title: {
        type: String,
        default: '',
      },
      value: {
        type: Number,
        default: 0,
      },
      isLast: {
        type: Boolean,
        default: false,
      },
    },
  };
</script>

<style lang="less" scoped>
  .panel-col {
    padding-left: 43px;
    border-right: 1px solid rgb(var(--gray-2));
  }

  .panel-col-no-border {
    padding-left: 43px;
  }

  .col-avatar {
    margin-right: 12px;
    background-color: var(--color-fill-2);
  }

  .up-icon {
    color: rgb(var(--red-6));
  }

  .unit {
    margin-left: 8px;
    color: rgb(var(--gray-8));
    font-size: 12px;
  }

  :deep(.panel-border) {
    margin: 4px 0 0 0;
  }
</style>
