<template>
  <BModal
    v-bind:model-value="modelValue"
    v-on:update:model-value="(data) => $emit('update:modelValue', data)"
    :title="title"
  >
    <div>
      <BFormInput
        v-model="numberText"
        type="number"
        :min="min"
        :placeholder="placeholder"
        :state="isValid"
        @keyup.enter="onSave"
        @input="validateInput"
      />
      <BFormInvalidFeedback :state="isValid">
        Value must be {{ min !== undefined ? `at least ${min}` : "valid" }}
      </BFormInvalidFeedback>
    </div>
    <template #footer>
      <BButton @click="$emit('update:modelValue', false)">Close</BButton>
      <BButton @click="onSave" variant="primary" :disabled="!isValid">
        Save
      </BButton>
    </template>
  </BModal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";

const props = defineProps<{
  modelValue: boolean;
  title?: string;
  presetValue?: number;
  min?: number;
  placeholder?: string;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  save: [value: number | null];
}>();

const numberText = ref<string>("");
const isValid = computed(() => {
  if (numberText.value === "") return true;
  const num = Number(numberText.value);
  if (isNaN(num)) return false;
  if (props.min !== undefined && num < props.min) return false;
  return true;
});

watch(
  () => props.presetValue,
  (v) => {
    numberText.value = v?.toString() || "";
  },
  { immediate: true }
);

const onSave = () => {
  const num = numberText.value === "" ? null : Number(numberText.value);
  if (num !== null && props.min !== undefined && num < props.min) {
    return;
  }
  emit("save", num);
  emit("update:modelValue", false);
};

const validateInput = () => {
  isValid.value;
};
</script>
