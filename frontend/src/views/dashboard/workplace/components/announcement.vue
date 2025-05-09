<script lang="ts" setup>
import { ref } from 'vue'

const announcements = ref([
  {
    title: '技术资源奖励',
    content: '发布下列技术相关的模型，整合包，教程，代码等资源，每个种子可获得5000魔力奖励（重复的不计算）：\n- Stable Diffusion\n- Stable Video Diffusion\n- Stable Audio\n- Llama 2',
  },
])

const rules = ref([
  {
    title: '总则',
    subtitle: '不遵守这些将导致帐号被封！',
    items: [
      '请不要做管理员明文禁止的事情。',
      '不允许发送垃圾信息。',
      '账号保留规则：',
      [
        'Veteran User及以上等级用户会永远保留；',
        'Elite User及以上等级用户封存账号（在控制面板）后不会被删除帐号；',
        '封存账号的用户连续400天不登录将被删除帐号；',
        '未封存账号的用户连续 365 天不登录将被删除帐号；',
        '没有流量的用户（即上传/下载数据都为0）连续100天不登录将被删除帐号。',
      ],
      '一切作弊的帐号会被封，请勿心存侥幸。',
      '注册多个Hermes账号的用户将被禁止。',
      '不要把本站的种子文件上传到其他Tracker！(具体请看常见问题)',
      '第一次在论坛或服务器中的捣乱行为会受到警告，第二次您将永远无缘Hermes。',
    ],
  },
  {
    title: '下载规则',
    subtitle: '违规将会失去下载权限！',
    items: [
      '过低的分享率会导致严重的后果，包括禁止帐号。详见常见问题。',
      '种子促销规则：',
      [
        '随机促销（种子上传后系统自动随机设为促销）：',
        '10%的概率成为"50%下载"；',
        '5%的概率成为"免费"；',
        '5%的概率成为"2x上传"；',
        '3%的概率成为"50%下载&2x上传"；',
        '1%的概率成为"免费&2x上传"。',
        '文件总体积大于20GB的种子将自动成为"免费"。',
        'Blu-ray Disk, HD DVD原盘将成为"免费"。',
        '电视剧等每季的第一集将成为"免费"。',
        '关注度高的种子将被设置为促销（由管理员定夺）。',
      ],
      '促销时限：',
      [
        '除了"2x上传"以外，其余类型的促销限时7天（自种子发布时起）；',
        '"2x上传"无时限。',
        '所有种子在发布1个月后将自动永久成为"2x上传"。',
        '我们也将不定期开启全站免费，届时尽情狂欢吧~',
      ],
      '你只能使用允许的BT客户端下载本站资源。详见常见问题。',
    ],
  },
  {
    title: '上传规则',
    subtitle: '谨记: 违规的种子将不经提醒而直接删除',
    items: [
      '上传总则：',
      [
        '上传者必须对上传的文件拥有合法的传播权。',
        '上传者必须保证上传速度与做种时间。如果在其他人完成前撤种或做种时间不足24小时，或者故意低速上传，上传者将会被警告甚至取消上传权限。',
        '对于自己发布的种子，发布者将获得双倍的上传量。',
        '如果你有一些违规但却有价值的资源，请将详细情况告知管理组，我们可能破例允许其发布。',
      ],
      '上传者资格：',
      [
        '任何人都能发布资源。',
        '不过，有些用户需要先在候选区提交候选。详见常见问题中的相关说明。',
        '对于游戏类资源，只有上传员及以上等级的用户，或者是管理组特别指定的用户，才能自由上传。其他用户必须先在候选区提交候选。',
      ],
      '允许的资源和文件：',
      [
        '高清（HD）视频，包括',
        '完整高清媒介，如蓝光（Blu-ray）原碟、HD DVD原碟等，或remux，',
        'HDTV流媒体，',
        '来源于上述媒介的高清重编码（至少为720p标准），',
        '其他高清视频，如高清DV；',
        '标清（SD）视频，只能是',
        '来源于高清媒介的标清重编码（至少为480p标准）；',
        'DVDR/DVDISO，',
        'DVDRip、CNDVDRip；',
        '无损音轨（及相应cue表单），如FLAC、Monkey\'s Audio等；',
        '5.1声道或以上标准的电影音轨、音乐音轨（DTS、DTSCD镜像等），评论音轨；',
        'PC游戏（必须为原版光盘镜像）；',
        '7日内发布的高清预告片；',
        '与高清相关的软件和文档。',
      ],
      '不允许的资源和文件：',
      [
        '总体积小于100MB的资源；',
        '标清视频upscale或部分upscale而成的视频文件；',
        '属于标清级别但质量较差的视频文件，包括CAM、TC、TS、SCR、DVDSCR、R5、R5.Line、HalfCD等；',
        'RealVideo编码的视频（通常封装于RMVB或RM）、flv文件；',
        '单独的样片（样片请和正片一起上传）；',
        '未达到5.1声道标准的有损音频文件，如常见的有损MP3、有损WMA等；',
        '无正确cue表单的多轨音频文件；',
        '硬盘版、高压版的游戏资源，非官方制作的游戏镜像，第三方mod，小游戏合集，单独的游戏破解或补丁；',
        'RAR等压缩文件；',
        '重复（dupe）的资源（判定规则见下文）；',
        '涉及禁忌或敏感内容（如色情、敏感政治话题等）的资源；',
        '损坏的文件，指在读取或回放过程中出现错误的文件；',
        '垃圾文件，如病毒、木马、网站链接、广告文档、种子中包含的种子文件等，或无关文件。',
      ],
    ],
  },
])
</script>

<template>
  <div class="announcement-container">
    <!-- 站点公告 -->
    <a-card class="mb-4">
      <template #title>
        站点公告
      </template>
      <a-collapse>
        <a-collapse-item v-for="(announcement, index) in announcements" :key="index" :header="announcement.title">
          <div class="whitespace-pre-line">
            {{ announcement.content }}
          </div>
        </a-collapse-item>
      </a-collapse>
    </a-card>

    <!-- 站点规则 -->
    <a-card>
      <template #title>
        站点规则
      </template>
      <a-collapse>
        <a-collapse-item v-for="(rule, index) in rules" :key="index" :header="rule.title">
          <template #extra>
            {{ rule.subtitle }}
          </template>
          <div class="rule-content">
            <template v-for="(item, itemIndex) in rule.items" :key="itemIndex">
              <div v-if="Array.isArray(item)" class="ml-4">
                <div v-for="(subItem, subIndex) in item" :key="subIndex" class="rule-item">
                  {{ subItem }}
                </div>
              </div>
              <div v-else class="rule-item">
                {{ item }}
              </div>
            </template>
          </div>
        </a-collapse-item>
      </a-collapse>
    </a-card>
  </div>
</template>

<style scoped lang="less">
.announcement-container {
  width: 100%;
}

.rule-content {
  .rule-item {
    margin-bottom: 8px;
    line-height: 1.5;
  }
}
</style>
