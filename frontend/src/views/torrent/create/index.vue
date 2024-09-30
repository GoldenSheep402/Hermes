<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {CategoryService, TorrentService} from "@/services/grpc.ts";
import {Category as CategoryBase} from "@/lib/proto/category/v1/category.pb.ts"
import {FileItem} from "@arco-design/web-vue";
import {CreateTorrentV1Request} from "@/lib/proto/torrent/v1/torrent.pb.ts";


interface Category {
  id: string;
  name: string;
  description: string;
}

const categoryList = ref<Category[]>([] as Category[]);

function fetchCategoryList() {
  CategoryService.GetCategoryList({}).then((res) => {
    for (let i = 0; i < res.category!!.length; i++) {
      categoryList.value.push({
        id: res.category!![i].id!!,
        name: res.category!![i].name!!,
        description: res.category!![i].description!!,
      });

      CategoryService.GetCategory({id: res.category!![i].id!!}).then((_res) => {
        for (let j = 0; j < _res.category!!.metaData!!.length; j++) {
          categoryFullInfo.value.push({
            id: _res.category!!.id!!,
            name: _res.category!!.name!!,
            description: _res.category!!.description!!,
            metaData: [
              {
                type: _res.category!!.metaData!![j].type!!,
                id: _res.category!!.metaData!![j].id!!,
                order: _res.category!!.metaData!![j].order!!,
                categoryId: _res.category!!.metaData!![j].categoryId!!,
                description: _res.category!!.metaData!![j].description!!,
                key: _res.category!!.metaData!![j].key!!,
                defaultValue: _res.category!!.metaData!![j].defaultValue!!,
                value: _res.category!!.metaData!![j].value!!,
              }
            ]
          });
        }
      }).catch((err) => {
        console.error(err);
      });
    }
  }).catch((err) => {
    console.error(err);
  });
}

const checkedID = ref('');
const categoryFullInfo = ref<CategoryBase[]>([]);


const torrentComment = ref("");

const sendFile = () => {
  const req = ref<CreateTorrentV1Request>({});
  req.value.comment = torrentComment.value;
  req.value.categoryId = checkedID.value;

  if (!req.value.metadata) {
    req.value.metadata = [];
  }

  for (let i = 0; i < categoryFullInfo.value.length; i++) {
    if (categoryFullInfo.value[i].id === checkedID.value) {
      if (categoryFullInfo.value[i].metaData) {
        for (let j = 0; j < categoryFullInfo.value[i].metaData!!.length; j++) {
          req.value.metadata!!.push({
            id: categoryFullInfo.value[i].metaData!![j].id,
            categoryId: checkedID.value,
            value: categoryFullInfo.value[i].metaData!![j].value?.toString(),
          });
        }
      }
    }
  }

  if (uint8Array.value) {
    req.value.torrent = {
      data: uint8Array.value,
    };
  }

  TorrentService.CreateTorrentV1(req.value).then((res) => {
    console.log('CreateTorrentV1:', res);
  }).catch((err) => {
    console.error(err);
  });
};


const fileList = ref<FileItem[]>([]);

const uint8Array = ref<Uint8Array | null>(null);

function handleFileChange(fileList: FileItem[], file: FileItem) {
  const fileObj = file.file as File;
  if (fileObj) {
    const reader = new FileReader();
    reader.onload = (event) => {
      const content = event.target?.result as ArrayBuffer;
      uint8Array.value = new Uint8Array(content);
      console.log('File Content as Uint8Array:', uint8Array.value);
    };
    reader.readAsArrayBuffer(fileObj);
  }
}

onMounted(() => {
  console.log('onMounted');
  fetchCategoryList()
});
</script>

<template>
  <div class="p-5">
    <div class="p-5 bg-[--color-bg-2]">
      <div class="p-0.5 text-20px leading-[1.4] font-500 text-[--color-text-1] mb-5">
        发布种子
      </div>

      <div>
        <div class="p-0.5 text-20px leading-[1.4] font-500 text-[--color-text-1] mb-5">
          选择类别
        </div>

        <div>
          <a-radio-group v-model="checkedID" class="mb-5" type="button">
            <div v-for="category in categoryList" :key="category.id">
              <a-radio :value="category.id">
                {{ category.name }}
              </a-radio>
            </div>
          </a-radio-group>

          <div v-if="checkedID !== ''">
            <div v-for="metas in categoryFullInfo" :key="metas.id">
              <div v-if="metas.id === checkedID">
                <div v-for="meta in metas.metaData!!" :key="meta.id" class="flex items-center mb-4">
                  <div class="flex items-center mr-4">
                    <div class="flex items-center justify-center text-gray-700 font-semibold">
                      {{ meta.key }}:
                    </div>
                  </div>
                  <div class="flex-grow">
                    <a-input
                        v-if="meta.type === 'number'"
                        v-model="meta.value"
                        :max="100"
                        :min="0"
                        :step="1"
                        class="w-full"
                    ></a-input>
                    <a-switch
                        v-else-if="meta.type === 'switch'"
                        v-model="meta.value"
                        class="mr-2"
                    ></a-switch>
                    <a-input
                        v-else-if="meta.type === 'string'"
                        v-model="meta.value"
                        class="w-full"
                    ></a-input>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>


      <div class="flex flex-col space-y-4">
        <div class="flex items-center">
          <div class="w-25 font-semibold text-[--color-text-1]">
            Comment:
          </div>
          <a-input v-model="torrentComment" class="flex-grow" placeholder="请输入Comment"></a-input>
          <a-button type="primary" @click="sendFile" class="ml-5">发布</a-button>
        </div>

        <a-upload
            :auto-upload="false"
            :file-list="fileList"
            @change="handleFileChange"
        >
        </a-upload>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>

</style>