<script setup lang="ts">
import { ref, computed } from 'vue'
import { Message } from '@arco-design/web-vue'

interface ShopItem {
  id: number
  name: string
  description: string
  price: number
  type: 'upload' | 'download' | 'ticket' | 'title' | 'vip' | 'gift' | 'ad' | 'charity' | 'hr' | 'rainbow'
  amount?: number
  unit?: string
}

const currentMagicPoints = ref(293712.8)

const shopItems = ref<ShopItem[]>([
  {
    id: 1,
    name: '1.0 GB上传量',
    description: '如果有足够的魔力值，你可以用它来换取上传量。交易完成后，你的魔力值会减少，上传量则会增加。',
    price: 1200,
    type: 'upload',
    amount: 1,
    unit: 'GB',
  },
  {
    id: 2,
    name: '5.0 GB上传量',
    description: '如果有足够的魔力值，你可以用它来换取上传量。交易完成后，你的魔力值会减少，上传量则会增加。',
    price: 3200,
    type: 'upload',
    amount: 5,
    unit: 'GB',
  },
  {
    id: 3,
    name: '10.0 GB上传量',
    description: '如果有足够的魔力值，你可以用它来换取上传量。交易完成后，你的魔力值会减少，上传量则会增加。',
    price: 4800,
    type: 'upload',
    amount: 10,
    unit: 'GB',
  },
  {
    id: 4,
    name: '1 个限时邀请',
    description: '购买后你会获得一个有效期为7天的邀请名额，如果过期未使用将自动消失，请合理购买。',
    price: 588888,
    type: 'ticket',
  },
  {
    id: 5,
    name: '自定义头衔',
    description: '如果有足够的魔力值，你就可以给自己定一个头衔了。注意：禁止使用脏话、攻击性的词汇或用户等级作为头衔。',
    price: 50000,
    type: 'title',
  },
  {
    id: 6,
    name: '贵宾待遇',
    description: '如果有足够的魔力值，你可以用它来换取一个月的贵宾待遇。交易完成后，你的魔力值会减少，同时你的等级将变为贵宾。',
    price: 200000000,
    type: 'vip',
  },
  {
    id: 11,
    name: '100.0 GB上传量',
    description: '如果有足够的魔力值，你可以用它来换取上传量。交易完成后，你的魔力值会减少，上传量则会增加。',
    price: 40000,
    type: 'upload',
    amount: 100,
    unit: 'GB',
  },
  {
    id: 12,
    name: '增加100.0 GB下载量',
    description: '如果有足够的魔力值，你可以用它来换取下载量。交易完成后，你的魔力值会减少，下载量则会增加。需上传量大于1TB。',
    price: 40000,
    type: 'download',
    amount: 100,
    unit: 'GB',
  },
  {
    id: 13,
    name: '减少100.0 GB下载量',
    description: '如果有足够的魔力值，你可以用它来减少下载量。交易完成后，你的魔力值会减少，下载量也会减少。需上传量大于1TB。',
    price: 80000,
    type: 'download',
    amount: -100,
    unit: 'GB',
  },
  {
    id: 14,
    name: '彩虹ID(月)',
    description: '如果有足够的魔力值，你可以用它来换取一个月的彩虹ID显示。用户名留空是买给自己，填写可赠送给其他会员，需要加收20%手续费。',
    price: 60000,
    type: 'rainbow',
  },
  {
    id: 15,
    name: '彩虹ID(年)',
    description: '如果有足够的魔力值，你可以用它来换取一年的彩虹ID显示。用户名为空是买给自己，填写可赠送给其他会员，需要加收20%手续费。',
    price: 720000,
    type: 'rainbow',  
  },
])

const canAfford = (price: number) => {
  return currentMagicPoints.value >= price
}

const handlePurchase = (item: ShopItem) => {
  if (!canAfford(item.price)) {
    Message.error('魔力值不足')
    return
  }
  // TODO: 实现购买逻辑
  Message.success('购买成功')
  currentMagicPoints.value -= item.price
}
</script>

<template>
  <div class="magic-shop-container">
    <div class="h-full p-5">
      <div class="h-full bg-[--color-bg-2] p-5">
        <div class="mb-5 flex items-center justify-between">
          <div class="p-0.5 text-30px text-[--color-text-1] font-500 leading-[1.4]">
            魔力值商店
          </div>
          <div class="flex items-center gap-2">
            <icon-gift class="text-24px text-[--color-primary-6]" />
            <span class="text-20px font-500 text-[--color-text-1]">
              当前魔力值：{{ currentMagicPoints.toLocaleString() }}
            </span>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4">
          <a-card v-for="item in shopItems" :key="item.id" class="info-card">
            <template #title>
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <icon-gift class="mr-2" />
                  {{ item.name }}
                </div>
                <div class="flex items-center gap-4">
                  <span class="text-[--color-text-2]">
                    {{ item.price.toLocaleString() }} 魔力值
                  </span>
                  <a-button
                    type="primary"
                    :disabled="!canAfford(item.price)"
                    @click="handlePurchase(item)"
                  >
                    购买
                  </a-button>
                </div>
              </div>
            </template>
            <div class="text-[--color-text-2]">
              {{ item.description }}
            </div>
          </a-card>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.info-card {
  margin-bottom: 16px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }
}
</style> 