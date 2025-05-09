<script setup lang="ts">
import { Tag } from '@arco-design/web-vue'
import { computed } from 'vue'

const props = defineProps<{
  ip: string
  type: 'v4' | 'v6'
}>()

function isPrivateIP(ip: string, type: 'v4' | 'v6'): boolean {
  if (!ip || ip.trim() === '')
    return false

  if (type === 'v4') {
    // IPv4私有地址范围
    const privateRanges = [
      { start: '10.0.0.0', end: '10.255.255.255' },
      { start: '172.16.0.0', end: '172.31.255.255' },
      { start: '192.168.0.0', end: '192.168.255.255' },
    ]

    const ipNum = ip.split('.').reduce((acc, octet) => (acc << 8) + Number.parseInt(octet), 0)

    return privateRanges.some((range) => {
      const start = range.start.split('.').reduce((acc, octet) => (acc << 8) + Number.parseInt(octet), 0)
      const end = range.end.split('.').reduce((acc, octet) => (acc << 8) + Number.parseInt(octet), 0)
      return ipNum >= start && ipNum <= end
    })
  }
  else {
    // IPv6私有地址范围
    return ip.startsWith('fe80:') || ip.startsWith('fc00:') || ip.startsWith('fd00:')
  }
}

const isPrivate = computed(() => isPrivateIP(props.ip, props.type))
</script>

<template>
  <div class="ip-tag">
    <span v-if="ip && ip.trim() !== ''" class="ip-address">{{ ip }}</span>
    <span v-else class="ip-address empty">-</span>
    <Tag v-if="ip && ip.trim() !== ''" :color="isPrivate ? 'arcoblue' : 'orangered'" size="small">
      {{ isPrivate ? '内网' : '公网' }}
    </Tag>
  </div>
</template>

<style scoped lang="less">
.ip-tag {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ip-address {
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: var(--color-text-1);
  letter-spacing: 0.5px;

  &.empty {
    color: var(--color-text-3);
  }
}
</style>
