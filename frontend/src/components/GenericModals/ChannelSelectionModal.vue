<template>
  <BModal
    v-bind:model-value="modelValue"
    v-on:update:model-value="(data) => $emit('update:modelValue', data)"
    title="Manage Channels"
    size="lg"
  >
    <div class="channel-selection">
      <BFormCheckboxGroup
        v-model="selectedChannels"
        :options="channels"
        value-field="id"
        text-field="name"
        stacked
      ></BFormCheckboxGroup>
    </div>
    <template #footer>
      <BButton variant="secondary" @click="$emit('update:modelValue', false)">
        Cancel
      </BButton>
      <BButton variant="primary" @click="onSave">Save Changes</BButton>
    </template>
  </BModal>
</template>

<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type { ChannelResponse } from "@/types/pocketbase-types";
import { useToastController } from "bootstrap-vue-next";
import { onMounted, ref, watch } from "vue";

const props = defineProps<{
  modelValue: boolean;
  presetChannels?: string[];
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  save: [channels: string[]];
}>();

const pb = usePb();
const toast = useToastController();
const channels = ref<ChannelResponse[]>([]);
const selectedChannels = ref<string[]>([]);

const fetchChannels = async () => {
  try {
    const resp = await pb.collection("channel").getFullList<ChannelResponse>();
    channels.value = resp;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch channels",
        variant: "danger",
      },
    });
  }
};

watch(
  () => props.presetChannels,
  (newVal) => {
    selectedChannels.value = newVal || [];
  },
  { immediate: true }
);

watch(
  () => props.modelValue,
  async (newVal) => {
    if (newVal) {
      await fetchChannels();
    }
  }
);

const onSave = () => {
  emit("save", selectedChannels.value);
  emit("update:modelValue", false);
};

onMounted(() => {
  if (props.modelValue) {
    fetchChannels();
  }
});
</script>

<style scoped>
.channel-selection {
  max-height: 300px;
  overflow-y: auto;
}
</style>
