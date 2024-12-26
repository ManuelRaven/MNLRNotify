<template>
  <div>
    <TextResponseModal
      v-model="showTextResponseModal"
      :title="'Edit Reciever'"
      :presetText="currentEditReciever?.name"
      @close="showTextResponseModal = false"
      @save="onNameChange"
    ></TextResponseModal>

    <ChannelSelectionModal
      v-model="showChannelModal"
      :preset-channels="currentEditReciever?.channels"
      @save="updateRecieverChannels"
    />

    <div class="mb-3 d-flex gap-2">
      <BFormInput
        v-model="newReciever.name"
        placeholder="Enter Name"
        class="me-2"
      />
      <BFormSelect
        v-model="newReciever.type"
        :options="typeOptions"
        class="me-2 w-auto"
      />
      <BButton variant="primary" @click="createReciever">Add Receiver</BButton>
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
      :items="recievers"
      :fields="sortFields"
    >
      <template #cell(channels)="row">
        {{ row.item.expand?.channels?.map((c) => c.name).join(", ") }}
        <span v-if="!row.item.expand?.channels?.length" class="text-danger"
          >No channel</span
        >
      </template>
      <template #cell(token)="row">
        <div class="d-flex align-items-center gap-2">
          {{
            showTokens[row.item.id] ? row.item.token : maskToken(row.item.token)
          }}
          <BButton
            size="sm"
            variant="link"
            class="p-0"
            @click="() => toggleToken(row.item.id)"
          >
            {{ showTokens[row.item.id] ? "Hide" : "Show" }}
          </BButton>
        </div>
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
              currentEditReciever = row.item;
              showTextResponseModal = true;
            }
          "
        >
          Edit Name
        </BButton>

        <BButton
          size="sm"
          variant="info"
          @click="
            () => {
              currentEditReciever = row.item;
              showChannelModal = true;
            }
          "
        >
          Manage Channels
        </BButton>
      </template>
    </BTable>
  </div>
</template>

<script setup lang="ts">
import ChannelSelectionModal from "@/components/GenericModals/ChannelSelectionModal.vue";
import TextResponseModal from "@/components/GenericModals/TextResponseModal.vue";
import { usePb } from "@/composeables/usePb";
import type { ExpandChannelsNameId } from "@/types/custom-types";
import {
  RecieverTypeOptions,
  type RecieverRecord,
  type RecieverResponse,
} from "@/types/pocketbase-types";
import {
  useModalController,
  useToastController,
  type TableFieldRaw,
} from "bootstrap-vue-next";
import { onMounted, ref } from "vue";

const showTextResponseModal = ref(false);
const showChannelModal = ref(false);
const currentEditReciever = ref<RecieverResponse | null>(null);

const typeOptions = [
  { value: "gotify", text: "Gotify" },
  { value: "ntfy", text: "Ntfy" },
];

const pb = usePb();
const toast = useToastController();
const { confirm } = useModalController();

const perPage = 10;
const currentPage = ref(1);
const rowsCount = ref(0);

const newReciever = ref<Omit<RecieverRecord, "id">>({
  name: "",
  type: RecieverTypeOptions.gotify,
  owner: "",
});

const recievers = ref<RecieverResponse<ExpandChannelsNameId>[]>([]);

const maskToken = (token?: string): string => {
  if (!token) return "";
  return token.slice(0, 4) + "..." + token.slice(-4);
};

const showTokens = ref<Record<string, boolean>>({});

const toggleToken = (id: string) => {
  showTokens.value[id] = !showTokens.value[id];
};

const onNameChange = async (newName: string | null) => {
  showTextResponseModal.value = false;
  if (currentEditReciever.value && newName) {
    try {
      await pb.collection("reciever").update(currentEditReciever.value.id, {
        name: newName,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Receiver updated successfully",
          variant: "success",
        },
      });

      await fetchRecievers();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update receiver",
          variant: "danger",
        },
      });
    }
  }
};

const updateRecieverChannels = async (channels: string[]) => {
  if (!currentEditReciever.value) return;

  try {
    await pb.collection("reciever").update(currentEditReciever.value.id, {
      channels: channels,
    });

    toast.show?.({
      props: {
        title: "Success",
        body: "Channels updated successfully",
        variant: "success",
      },
    });

    await fetchRecievers();
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to update channels",
        variant: "danger",
      },
    });
  }
};

const onDelete = async (reciever: RecieverResponse, index: number) => {
  try {
    const value = await confirm?.({
      props: {
        title: "Are you sure?",
        body: "Do you want to delete this receiver? This action cannot be undone.",
      },
    });

    if (!value) return;
    await pb.collection("reciever").delete(reciever.id);
    toast.show?.({
      props: { title: "Success", body: "Receiver Deleted", variant: "success" },
    });

    recievers.value.splice(index, 1);
  } catch (error) {
    toast.show?.({
      props: { title: "Error", body: (error as any) || "", variant: "danger" },
    });
  }
};

const createReciever = async () => {
  if (!newReciever.value.name) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Please fill in all fields",
        variant: "danger",
      },
    });
    return;
  }

  try {
    await pb.collection("reciever").create<RecieverResponse>(newReciever.value);
    newReciever.value.name = "";
    await fetchRecievers();
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to create receiver",
        variant: "danger",
      },
    });
  }
};

const sortFields: Exclude<TableFieldRaw<RecieverResponse>, string>[] = [
  { key: "id", label: "ID", sortable: true },
  { key: "name", label: "Name", sortable: true },
  { key: "type", label: "Type", sortable: true },
  { key: "token", label: "Token" },
  { key: "channels", label: "Channels" },
  { key: "actions", label: "Actions" },
];

const fetchRecievers = async () => {
  try {
    let resp = await pb
      .collection("reciever")
      .getList<RecieverResponse<ExpandChannelsNameId>>(
        currentPage.value,
        perPage,
        {
          expand: "channels",
        }
      );

    recievers.value = resp.items;
    rowsCount.value = resp.totalItems;

    // Reset token visibility for new items
    resp.items.forEach((item) => {
      if (!(item.id in showTokens.value)) {
        showTokens.value[item.id] = false;
      }
    });
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch receivers",
        variant: "danger",
      },
    });
  }
};

onMounted(async () => {
  await fetchRecievers();
});
</script>

<style scoped></style>
