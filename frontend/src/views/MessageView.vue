<template>
  <div>
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
      :items="messages"
      :fields="sortFields"
    >
      <template #cell(channel)="row">
        {{ row.item.expand?.channel?.name }}
      </template>
      <template #cell(text)="row">
        <span class="text-break">{{ row.item.text }}</span>
      </template>
      <template #cell(deliveryState)="row">
        <BBadge :variant="getDeliveryStateVariant(row.item.deliveryState)">
          {{ row.item.deliveryState || "pending" }}
        </BBadge>
      </template>
      <template #cell(deliveryMessage)="row">
        <span class="text-break">{{ row.item.deliveryMessage }}</span>
      </template>
      <template #cell(actions)="row">
        <BButton
          size="sm"
          variant="danger"
          @click="onDelete(row.item, row.index)"
        >
          Delete
        </BButton>
      </template>
    </BTable>
  </div>
</template>

<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type { ExpandChannelNameIdSingle } from "@/types/custom-types";
import type {
  MessageDeliveryStateOptions,
  MessageResponse,
} from "@/types/pocketbase-types";
import {
  useModalController,
  useToastController,
  type TableFieldRaw,
} from "bootstrap-vue-next";
import { onMounted, ref, watch } from "vue";

const pb = usePb();
const toast = useToastController();
const { confirm } = useModalController();

const perPage = 10;
const currentPage = ref(1);
const rowsCount = ref(0);

const messages = ref<MessageResponse<ExpandChannelNameIdSingle>[]>([]);

const onDelete = async (message: MessageResponse, index: number) => {
  try {
    const value = await confirm?.({
      props: {
        title: "Are you sure?",
        body: "Do you want to delete this message? This action cannot be undone.",
      },
    });

    if (!value) return;
    await pb.collection("message").delete(message.id);
    toast.show?.({
      props: { title: "Success", body: "Message Deleted", variant: "success" },
    });

    messages.value.splice(index, 1);
  } catch (error) {
    toast.show?.({
      props: { title: "Error", body: (error as any) || "", variant: "danger" },
    });
  }
};

const getDeliveryStateVariant = (state?: MessageDeliveryStateOptions) => {
  switch (state) {
    case "success":
      return "success";
    case "failure":
      return "danger";
    default:
      return "warning";
  }
};

const sortFields: Exclude<TableFieldRaw<MessageResponse>, string>[] = [
  { key: "created", label: "Created", sortable: true },
  { key: "channel", label: "Channel", sortable: true },
  { key: "text", label: "Message" },
  { key: "deliveryState", label: "Status", sortable: true },
  { key: "deliveryMessage", label: "Delivery Message" },
  { key: "actions", label: "Actions" },
];

const fetchMessages = async () => {
  try {
    let resp = await pb
      .collection("message")
      .getList<MessageResponse<ExpandChannelNameIdSingle>>(
        currentPage.value,
        perPage,
        {
          sort: "-created",
          expand: "channel",
          fields: "*,expand.channel.name,expand.channel.id",
        }
      );

    messages.value = resp.items;
    rowsCount.value = resp.totalItems;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch messages",
        variant: "danger",
      },
    });
  }
};

onMounted(async () => {
  await Promise.all([fetchMessages()]);
});

// Watch page changes
watch(currentPage, () => {
  fetchMessages();
});
</script>

<style scoped>
.text-break {
  word-break: break-word;
  white-space: pre-wrap;
}
</style>
