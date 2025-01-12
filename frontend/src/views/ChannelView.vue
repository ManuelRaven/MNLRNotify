<template>
  <div>
    <TextResponseModal
      v-model="showTextResponseModal"
      :title="'Edit Channel'"
      :presetText="currentEditChannel?.name"
      @close="showTextResponseModal = false"
      @save="onNameChange"
    ></TextResponseModal>

    <NumberResponseModal
      v-model="showLifetimeModal"
      title="Edit Lifetime in Seconds"
      :preset-value="currentEditChannel?.lifetime_seconds"
      :min="0"
      placeholder="Enter lifetime in seconds"
      @save="onLifetimeChange"
    />

    <ChannelMutationsModal
      v-model="showMutationsModal"
      :preset-mutations="currentEditChannel?.mutations"
      @save="onMutationsChange"
    />

    <div class="mb-3 d-flex gap-2">
      <BFormInput
        v-model="newChannel.name"
        placeholder="Enter Name"
        @keyup.enter="createChannel"
      />

      <BButton variant="primary" @click="createChannel">Add Channel</BButton>
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
      :items="channels"
      :fields="sortFields"
    >
      <template #cell(lifetime_seconds)="row">
        {{ formatLifetime(row.item.lifetime_seconds) }}
      </template>
      <template #cell(mutations)="row">
        {{
          row.item.expand?.mutations?.map((m) => m.name).join(" → ") ||
          "No mutations"
        }}
      </template>
      <template #cell(actions)="row">
        <BButton
          size="sm"
          variant="danger"
          class="me-1"
          @click="onDelete(row.item, row.index)"
        >
          Delete
        </BButton>

        <BButton
          size="sm"
          variant="primary"
          class="me-1"
          @click="
            () => {
              currentEditChannel = row.item;
              showTextResponseModal = true;
            }
          "
        >
          Edit Name
        </BButton>

        <BButton
          size="sm"
          variant="primary"
          @click="
            () => {
              currentEditChannel = row.item;
              showLifetimeModal = true;
            }
          "
        >
          Edit Lifetime
        </BButton>

        <BButton
          size="sm"
          variant="primary"
          @click="
            () => {
              currentEditChannel = row.item;
              showMutationsModal = true;
            }
          "
        >
          Manage Mutations
        </BButton>
      </template>
    </BTable>
  </div>
</template>

<script setup lang="ts">
import ChannelMutationsModal from "@/components/GenericModals/ChannelMutationsModal.vue";
import NumberResponseModal from "@/components/GenericModals/NumberResponseModal.vue";
import TextResponseModal from "@/components/GenericModals/TextResponseModal.vue";
import { usePb } from "@/composeables/usePb";
import type {
  CreateChannelRequest,
  ExpandMutations,
} from "@/types/custom-types";
import type { ChannelResponse } from "@/types/pocketbase-types";
import {
  useModalController,
  useToastController,
  type TableFieldRaw,
} from "bootstrap-vue-next";
import { onMounted, ref } from "vue";

// Modals
const showTextResponseModal = ref(false);
const showLifetimeModal = ref(false);
const showMutationsModal = ref(false);
const currentEditChannel = ref<ChannelResponse<ExpandMutations> | null>(null);

const onNameChange = async (newName: string | null) => {
  console.log("newName", newName);
  showTextResponseModal.value = false;
  if (currentEditChannel.value && newName) {
    try {
      await pb.collection("channel").update(currentEditChannel.value.id, {
        name: newName,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Data updated successfully",
          variant: "success",
        },
      });

      await fetchChannels();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update data",
          variant: "danger",
        },
      });
    }
  }
};

const onLifetimeChange = async (lifetime: number | null) => {
  if (currentEditChannel.value) {
    try {
      await pb.collection("channel").update(currentEditChannel.value.id, {
        lifetime_seconds: lifetime,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Lifetime updated successfully",
          variant: "success",
        },
      });

      await fetchChannels();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update lifetime",
          variant: "danger",
        },
      });
    }
  }
};

const onMutationsChange = async (mutations: string[]) => {
  showMutationsModal.value = false;
  if (currentEditChannel.value) {
    try {
      await pb.collection("channel").update(currentEditChannel.value.id, {
        mutations: mutations,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Mutations updated successfully",
          variant: "success",
        },
      });

      await fetchChannels();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update mutations",
          variant: "danger",
        },
      });
    }
  }
};

const pb = usePb();
const toast = useToastController();
const { confirm } = useModalController();

const perPage = 10;
const currentPage = ref(1);
const rowsCount = ref(0);

const newChannel = ref<CreateChannelRequest>({
  name: "",
});

const channels = ref<ChannelResponse<ExpandMutations>[]>([]);

const onDelete = async (channel: ChannelResponse, index: number) => {
  try {
    const value = await confirm?.({
      props: {
        title: "Are you sure?",
        body: "Do you want to delete this channel? This action cannot be undone.",
      },
    });

    if (!value) return;
    await pb.collection("channel").delete(channel.id);
    toast.show?.({
      props: { title: "Success", body: "Bookmark Deleted", variant: "success" },
    });

    channels.value.splice(index, 1);
  } catch (error) {
    toast.show?.({
      props: { title: "Error", body: (error as any) || "", variant: "danger" },
    });
  }
};

const createChannel = async () => {
  await pb.collection("channel").create<ChannelResponse>({
    name: newChannel.value.name,
  });

  newChannel.value.name = "";

  await fetchChannels();
};

const sortFields: Exclude<TableFieldRaw<ChannelResponse>, string>[] = [
  { key: "name", label: "Name", sortable: true },
  { key: "lifetime_seconds", label: "Lifetime of Messages", sortable: true },
  { key: "mutations", label: "Mutations" },
  { key: "actions", label: "Actions" },
];

const fetchChannels = async () => {
  try {
    let resp = await pb
      .collection("channel")
      .getList<ChannelResponse<ExpandMutations>>(currentPage.value, perPage, {
        expand: "mutations",
      });

    channels.value = resp.items;
    rowsCount.value = resp.totalItems;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch data",
        variant: "danger",
      },
    });
  }
};

const formatLifetime = (seconds?: number): string => {
  if (!seconds) return "Never expires";
  if (seconds < 60) return `${seconds} seconds`;
  if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours`;
  return `${Math.floor(seconds / 86400)} days`;
};

onMounted(() => {
  fetchChannels();
});
</script>

<style scoped></style>
