<template>
  <BModal :title="title">
    <BFormInput v-model="text" />
    <template #footer>
      <BButton @click="emit('close')">Close</BButton>
      <BButton @click="onSave" variant="primary">Save</BButton>
    </template>
  </BModal>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

const text = ref<string | null>(null);

const { title = "Text", presetText } = defineProps<{
  title?: string;
  presetText?: string;
}>();

const onSave = (e: Event) => {
  emit("save", text.value);
};

watch(
  () => presetText,
  (v) => {
    text.value = v || "";
  }
);

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", text: string | null): void;
}>();
</script>

<style scoped></style>
