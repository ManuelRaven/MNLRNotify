<template>
  <BModal
    v-bind:model-value="modelValue"
    v-on:update:model-value="(data) => $emit('update:modelValue', data)"
    :title="
      currentStep === 1 ? 'Select Channel Mutations' : 'Order Channel Mutations'
    "
    size="lg"
  >
    <div v-if="currentStep === 1" class="mutations-list">
      <div
        v-for="mutation in mutations"
        :key="mutation.id"
        class="d-flex align-items-center mb-2"
      >
        <BFormCheckbox
          v-model="selectedMutations"
          :value="mutation.id"
          class="flex-grow-1"
        >
          {{ mutation.name }}
        </BFormCheckbox>
      </div>
    </div>

    <div v-else class="mutations-list">
      <div
        v-for="mutationId in selectedMutations"
        :key="mutationId"
        class="d-flex align-items-center mb-2"
      >
        <div class="flex-grow-1">
          {{ getMutationName(mutationId) }}
        </div>
        <div class="d-flex flex-column ms-2">
          <BButton
            size="sm"
            variant="outline-secondary"
            class="mb-1"
            :disabled="isFirst(mutationId)"
            @click="moveUp(mutationId)"
          >
            ↑
          </BButton>
          <BButton
            size="sm"
            variant="outline-secondary"
            :disabled="isLast(mutationId)"
            @click="moveDown(mutationId)"
          >
            ↓
          </BButton>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="w-100 d-flex justify-content-between">
        <BButton
          variant="secondary"
          @click="
            currentStep === 1
              ? $emit('update:modelValue', false)
              : currentStep--
          "
        >
          {{ currentStep === 1 ? "Cancel" : "Back" }}
        </BButton>
        <BButton
          variant="primary"
          @click="currentStep === 1 ? nextStep() : onSave()"
          :disabled="currentStep === 1 && selectedMutations.length === 0"
        >
          {{ currentStep === 1 ? "Next" : "Save Changes" }}
        </BButton>
      </div>
    </template>
  </BModal>
</template>

<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type { ChannelMutationsResponse } from "@/types/pocketbase-types";
import { useToastController } from "bootstrap-vue-next";
import { onMounted, ref, watch } from "vue";

const props = defineProps<{
  modelValue: boolean;
  presetMutations?: string[];
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  save: [mutations: string[]];
}>();

const pb = usePb();
const toast = useToastController();
const mutations = ref<ChannelMutationsResponse[]>([]);
const selectedMutations = ref<string[]>([]);
const currentStep = ref(1);

const fetchMutations = async () => {
  try {
    const resp = await pb
      .collection("channel_mutations")
      .getFullList<ChannelMutationsResponse>();
    mutations.value = resp;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch mutations",
        variant: "danger",
      },
    });
  }
};

watch(
  () => props.presetMutations,
  (newVal) => {
    selectedMutations.value = newVal || [];
  },
  { immediate: true }
);

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      currentStep.value = 1;
      fetchMutations();
    }
  }
);

const onSave = () => {
  emit("save", selectedMutations.value);
  emit("update:modelValue", false);
};

const moveUp = (id: string) => {
  const index = selectedMutations.value.indexOf(id);
  if (index > 0) {
    const newArray = [...selectedMutations.value];
    [newArray[index - 1], newArray[index]] = [
      newArray[index],
      newArray[index - 1],
    ];
    selectedMutations.value = newArray;
  }
};

const moveDown = (id: string) => {
  const index = selectedMutations.value.indexOf(id);
  if (index < selectedMutations.value.length - 1) {
    const newArray = [...selectedMutations.value];
    [newArray[index], newArray[index + 1]] = [
      newArray[index + 1],
      newArray[index],
    ];
    selectedMutations.value = newArray;
  }
};

const isFirst = (id: string) => {
  return selectedMutations.value.indexOf(id) === 0;
};

const isLast = (id: string) => {
  return (
    selectedMutations.value.indexOf(id) === selectedMutations.value.length - 1
  );
};

const getMutationName = (id: string) => {
  return mutations.value.find((m) => m.id === id)?.name || "";
};

const nextStep = () => {
  if (selectedMutations.value.length > 0) {
    currentStep.value = 2;
  }
};

onMounted(() => {
  if (props.modelValue) {
    fetchMutations();
  }
});
</script>

<style scoped>
.mutations-list {
  max-height: 300px;
  overflow-y: auto;
}

.handle {
  cursor: move;
  color: #666;
}
</style>
