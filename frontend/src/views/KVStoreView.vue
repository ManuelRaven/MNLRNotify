<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>KV Store Management</h2>
      <BButton variant="primary" @click="showAddModal"> Add New Entry </BButton>
    </div>
    <BAlert
      :model-value="true"
      variant="warning"
      :dismissible="false"
      class="mb-3"
    >
      <p>
        <strong>Note:</strong> The KV Store is a simple key-value store. Keys
        must be unique on the whole Instance, and the value can be any
        JSON-serializable object.
        <br />

        Keys are only restricted at a user level - all Mutations can access,
        modify, or delete these entries, so use caution when storing sensitive
        information.
      </p>
    </BAlert>
    <div class="mb-3">
      <BInputGroup>
        <BFormInput
          v-model="searchQuery"
          placeholder="Search for keys..."
          @input="onSearch"
        />
        <BInputGroupAppend>
          <BButton @click="onSearch">Search</BButton>
        </BInputGroupAppend>
      </BInputGroup>
    </div>

    <b-pagination
      v-model="currentPage"
      :total-rows="rowsCount"
      :per-page="perPage"
    />

    <BTable
      stacked="md"
      striped
      label-stacked
      hover
      :items="kvItems"
      :fields="sortFields"
      :busy="isLoading"
    >
      <template #table-busy>
        <div class="text-center my-2">
          <BSpinner class="align-middle"></BSpinner>
          <strong> Loading data...</strong>
        </div>
      </template>

      <template #cell(id)="row">
        <small>{{ row.item.id }}</small>
      </template>

      <template #cell(key)="row">
        <span class="text-break">{{ row.item.id }}</span>
      </template>

      <template #cell(value)="row">
        <span class="text-break">{{ formatValue(row.item.value) }}</span>
      </template>

      <template #cell(created)="row">
        {{ new Date(row.item.created).toLocaleString() }}
      </template>

      <template #cell(updated)="row">
        {{ new Date(row.item.updated).toLocaleString() }}
      </template>

      <template #cell(actions)="row">
        <div class="d-flex gap-2">
          <BButton size="sm" variant="primary" @click="showEditModal(row.item)">
            Edit
          </BButton>
          <BButton
            size="sm"
            variant="danger"
            @click="onDelete(row.item, row.index)"
          >
            Delete
          </BButton>
        </div>
      </template>
    </BTable>

    <!-- Modal for adding or editing KV entries -->
    <BModal
      v-model="showModal"
      :title="isEditMode ? 'Edit KV Entry' : 'Add New KV Entry'"
      @ok="handleSave"
      @hidden="resetModal"
      ok-title="Save"
      cancel-title="Cancel"
    >
      <BForm>
        <BFormGroup label="Key:" label-for="key-input">
          <BFormInput
            id="key-input"
            v-model="editItem.key"
            :disabled="isEditMode"
            required
            placeholder="Enter key"
          ></BFormInput>
        </BFormGroup>

        <BFormGroup label="Value:" label-for="value-input">
          <BFormTextarea
            id="value-input"
            v-model="editItem.value"
            required
            placeholder="Enter value"
            rows="5"
          ></BFormTextarea>
          <small class="form-text text-muted">
            JSON objects can be entered as formatted values.
          </small>
        </BFormGroup>
      </BForm>
    </BModal>
  </div>
</template>

<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type { KVStoreResponse } from "@/types/pocketbase-types";
import {
  BvTriggerableEvent,
  useModalController,
  useToastController,
  type TableFieldRaw,
} from "bootstrap-vue-next";
import { onMounted, ref, watch } from "vue";

const pb = usePb();
const toast = useToastController();
const { confirm } = useModalController();

const isLoading = ref(false);
const perPage = 10;
const currentPage = ref(1);
const rowsCount = ref(0);
const searchQuery = ref("");

const kvItems = ref<KVStoreResponse[]>([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editItem = ref({
  id: "",
  key: "",
  value: "",
});

// Format value for display
const formatValue = (value: any): string => {
  if (value === null) return "null";
  if (typeof value === "object") {
    try {
      return JSON.stringify(value, null, 2);
    } catch (e) {
      return String(value);
    }
  }
  return String(value);
};

// Parse input value to actual value
const parseValue = (value: string): any => {
  if (!value || value === "null") return null;

  try {
    return JSON.parse(value);
  } catch (e) {
    // If not valid JSON, return it as string
    return value;
  }
};

const onDelete = async (item: KVStoreResponse, index: number) => {
  try {
    const value = await confirm?.({
      props: {
        title: "Are you sure?",
        body: "Do you really want to delete this KV entry? This action cannot be undone.",
      },
    });

    if (!value) return;

    await pb.collection("KVStore").delete(item.id);

    toast.show?.({
      props: {
        title: "Success",
        body: "KV entry has been deleted",
        variant: "success",
      },
    });

    kvItems.value.splice(index, 1);
    rowsCount.value--;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: (error as any)?.message || "Error deleting entry",
        variant: "danger",
      },
    });
  }
};

const showAddModal = () => {
  isEditMode.value = false;
  editItem.value = {
    id: "",
    key: "",
    value: "",
  };
  showModal.value = true;
};

const showEditModal = (item: KVStoreResponse) => {
  isEditMode.value = true;
  editItem.value = {
    id: item.id,
    key: item.id,
    value: formatValue(item.value),
  };
  showModal.value = true;
};

const resetModal = () => {
  editItem.value = {
    id: "",
    key: "",
    value: "",
  };
};

const handleSave = async (bvModalEvent: BvTriggerableEvent) => {
  bvModalEvent.preventDefault();

  try {
    const data = {
      value: parseValue(editItem.value.value),
    };

    if (isEditMode.value) {
      // Update existing item
      await pb.collection("KVStore").update(editItem.value.id, data);

      toast.show?.({
        props: {
          title: "Success",
          body: "KV entry has been updated",
          variant: "success",
        },
      });
    } else {
      // Create new item
      await pb.collection("KVStore").create({
        id: editItem.value.key,
        value: parseValue(editItem.value.value),
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "New KV entry has been added",
          variant: "success",
        },
      });
    }

    showModal.value = false;
    fetchKVItems();
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: (error as any)?.message || "Error saving entry",
        variant: "danger",
      },
    });
  }
};

const sortFields: Exclude<TableFieldRaw<KVStoreResponse>, string>[] = [
  { key: "id", label: "Key", sortable: true },
  { key: "value", label: "Value" },
  { key: "created", label: "Created on", sortable: true },
  { key: "updated", label: "Updated on", sortable: true },
  { key: "actions", label: "Actions" },
];

const fetchKVItems = async () => {
  isLoading.value = true;
  try {
    const filter = searchQuery.value ? `id ~ "${searchQuery.value}"` : "";

    const resp = await pb
      .collection("KVStore")
      .getList(currentPage.value, perPage, {
        sort: "-updated",
        filter: filter,
      });

    kvItems.value = resp.items;
    rowsCount.value = resp.totalItems;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Error loading KV entries",
        variant: "danger",
      },
    });
  } finally {
    isLoading.value = false;
  }
};

const onSearch = () => {
  currentPage.value = 1;
  fetchKVItems();
};

onMounted(async () => {
  await fetchKVItems();
});

// Load new data when changing pages
watch(currentPage, () => {
  fetchKVItems();
});
</script>

<style scoped>
.text-break {
  word-break: break-word;
  white-space: pre-wrap;
}
</style>
