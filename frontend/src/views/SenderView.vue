<template>
  <div>
    <TextResponseModal
      v-model="showTextResponseModal"
      :title="'Edit Sender'"
      :presetText="currentEditSender?.name"
      @close="showTextResponseModal = false"
      @save="onNameChange"
    ></TextResponseModal>

    <TextResponseModal
      v-model="showUrlModal"
      :title="'Edit Send URL'"
      :presetText="currentEditSender?.sendurl"
      @close="showUrlModal = false"
      @save="onUrlChange"
    ></TextResponseModal>

    <NumberResponseModal
      v-model="showSplitLimitModal"
      :title="'Edit Split Limit'"
      :preset-value="currentEditSender?.splitLimit"
      :min="0"
      :placeholder="'Enter split limit (0 for no splitting)'"
      :help-text="'Maximum number of characters per message. Messages longer than this will be split. Set to 0 to disable splitting. Example: Use 1000 for Telegram\'s character limit.'"
      @save="onSplitLimitChange"
    />

    <ChannelSelectionModal
      v-model="showChannelModal"
      :preset-channels="currentManageSender?.channel"
      @save="updateSenderChannels"
    />
    <BAccordion>
      <BAccordionItem title="Add New" class="mb-3">
        <div>
          <BCard class="mb-3">
            <template #header>
              <h5 class="mb-0">Shoutrrr URL Format Help</h5>
            </template>
            <p class="mb-2">
              URLs must follow the shoutrrr format:
              <code>service://parameters</code>
            </p>
            <div class="mb-2">
              <strong>Examples:</strong>
              <ul class="mb-0">
                <li><code>discord://token@id</code> - Discord webhook</li>
                <li>
                  <code>telegram://token@telegram?chats=@channel</code> -
                  Telegram channel
                </li>
                <li>
                  <code>matrix://user:password@host/channel</code> - Matrix room
                </li>
              </ul>
            </div>
            <BButton
              variant="link"
              class="p-0"
              href="https://containrrr.dev/shoutrrr/v0.8/services/overview/"
              target="_blank"
            >
              View full documentation →
            </BButton>
          </BCard>

          <div class="d-flex gap-2">
            <BFormInput
              v-model="newSender.name"
              placeholder="Enter Name"
              class="me-2"
            />
            <BFormInput
              v-model="newSender.sendurl"
              placeholder="Enter Send URL"
              class="me-2"
            />
            <BButton variant="primary" @click="createSender"
              >Add Sender</BButton
            >
          </div>
        </div>
      </BAccordionItem>
    </BAccordion>

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
      :items="senders"
      :fields="sortFields"
    >
      <template #cell(channels)="row">
        {{ row.item.expand?.channel?.map((c) => c.name).join(", ") }}
        <span v-if="!row.item.expand?.channel?.length" class="text-danger"
          >No channel</span
        >
      </template>
      <template #cell(sendurl)="row">
        {{ maskSendUrl(row.item.sendurl) }}
      </template>
      <template #cell(splitLimit)="row">
        <div class="d-flex align-items-center">
          {{ row.item.splitLimit || "No splitting" }}
          <BButton
            size="sm"
            variant="link"
            class="p-0 ms-1"
            v-b-tooltip.hover
            title="Maximum characters per message. Messages longer than this will be split into multiple messages. Set to 0 to disable splitting."
          >
            <IBiInfoCircle />
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
              currentEditSender = row.item;
              showTextResponseModal = true;
            }
          "
        >
          Edit Name
        </BButton>

        <BButton
          size="sm"
          variant="primary"
          class="me-1"
          @click="
            () => {
              currentEditSender = row.item;
              showUrlModal = true;
            }
          "
        >
          Edit URL
        </BButton>

        <BButton
          size="sm"
          variant="primary"
          class="me-1"
          @click="
            () => {
              currentEditSender = row.item;
              showSplitLimitModal = true;
            }
          "
        >
          Edit Split Limit
        </BButton>

        <BButton size="sm" variant="info" @click="openChannelModal(row.item)">
          Manage Channels
        </BButton>
      </template>
    </BTable>
  </div>
</template>

<script setup lang="ts">
import ChannelSelectionModal from "@/components/GenericModals/ChannelSelectionModal.vue";
import NumberResponseModal from "@/components/GenericModals/NumberResponseModal.vue";
import TextResponseModal from "@/components/GenericModals/TextResponseModal.vue";
import { usePb } from "@/composeables/usePb";
import type {
  CreateSenderRequest,
  ExpandChannelNameId,
} from "@/types/custom-types";
import type { SenderResponse } from "@/types/pocketbase-types";
import {
  useModalController,
  useToastController,
  type TableFieldRaw,
} from "bootstrap-vue-next";
import { onMounted, ref } from "vue";

const showTextResponseModal = ref(false);
const showUrlModal = ref(false);
const showSplitLimitModal = ref(false);
const currentEditSender = ref<SenderResponse | null>(null);

const onNameChange = async (newName: string | null) => {
  showTextResponseModal.value = false;
  if (currentEditSender.value && newName) {
    try {
      await pb.collection("sender").update(currentEditSender.value.id, {
        name: newName,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Sender updated successfully",
          variant: "success",
        },
      });

      await fetchSenders();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update sender",
          variant: "danger",
        },
      });
    }
  }
};

const validateShoutrrrUrl = (url: string): boolean => {
  const pattern = /^[a-z]+:\/\/.+/i;
  return pattern.test(url);
};

const onUrlChange = async (newUrl: string | null) => {
  showUrlModal.value = false;
  if (currentEditSender.value && newUrl) {
    if (!validateShoutrrrUrl(newUrl)) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Invalid Shoutrrr URL format. Please check the format help above.",
          variant: "danger",
        },
      });
      return;
    }

    try {
      await pb.collection("sender").update(currentEditSender.value.id, {
        sendurl: newUrl,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Send URL updated successfully",
          variant: "success",
        },
      });

      await fetchSenders();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update Send URL",
          variant: "danger",
        },
      });
    }
  }
};

const onSplitLimitChange = async (newLimit: number | null) => {
  if (currentEditSender.value) {
    try {
      await pb.collection("sender").update(currentEditSender.value.id, {
        splitLimit: newLimit,
      });

      toast.show?.({
        props: {
          title: "Success",
          body: "Split limit updated successfully",
          variant: "success",
        },
      });

      await fetchSenders();
    } catch (error) {
      toast.show?.({
        props: {
          title: "Error",
          body: "Failed to update split limit",
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

const newSender = ref<CreateSenderRequest>({
  name: "",
  sendurl: "",
});

const senders = ref<SenderResponse<ExpandChannelNameId>[]>([]);

const onDelete = async (sender: SenderResponse, index: number) => {
  try {
    const value = await confirm?.({
      props: {
        title: "Are you sure?",
        body: "Do you want to delete this sender? This action cannot be undone.",
      },
    });

    if (!value) return;
    await pb.collection("sender").delete(sender.id);
    toast.show?.({
      props: { title: "Success", body: "Sender Deleted", variant: "success" },
    });

    senders.value.splice(index, 1);
  } catch (error) {
    toast.show?.({
      props: { title: "Error", body: (error as any) || "", variant: "danger" },
    });
  }
};

const createSender = async () => {
  if (!newSender.value.name || !newSender.value.sendurl) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Please fill in all fields",
        variant: "danger",
      },
    });
    return;
  }

  if (!validateShoutrrrUrl(newSender.value.sendurl)) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Invalid Shoutrrr URL format. Please check the format help above.",
        variant: "danger",
      },
    });
    return;
  }

  try {
    await pb.collection("sender").create<SenderResponse>({
      name: newSender.value.name,
      sendurl: newSender.value.sendurl,
    });

    newSender.value.name = "";
    newSender.value.sendurl = "";

    await fetchSenders();
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to create sender",
        variant: "danger",
      },
    });
  }
};

const sortFields: Exclude<TableFieldRaw<SenderResponse>, string>[] = [
  { key: "name", label: "Name", sortable: true },
  { key: "sendurl", label: "Send URL", sortable: true },
  { key: "splitLimit", label: "Split Limit", sortable: true },
  { key: "channels", label: "Channels" },
  { key: "actions", label: "Actions" },
];

const fetchSenders = async () => {
  try {
    let resp = await pb
      .collection("sender")
      .getList<SenderResponse<ExpandChannelNameId>>(
        currentPage.value,
        perPage,
        {
          expand: "channel",
        }
      );

    senders.value = resp.items;
    rowsCount.value = resp.totalItems;
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to fetch senders",
        variant: "danger",
      },
    });
  }
};

const showChannelModal = ref(false);
const currentManageSender = ref<SenderResponse | null>(null);

const openChannelModal = async (sender: SenderResponse) => {
  currentManageSender.value = sender;
  showChannelModal.value = true;
};

const updateSenderChannels = async (channels: string[]) => {
  if (!currentManageSender.value) return;

  try {
    await pb.collection("sender").update(currentManageSender.value.id, {
      channel: channels,
    });

    toast.show?.({
      props: {
        title: "Success",
        body: "Channels updated successfully",
        variant: "success",
      },
    });

    await fetchSenders();
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

const maskSendUrl = (url: string): string => {
  const parts = url.split("://");
  if (parts.length !== 2) {
    return url.slice(0, 4) + "..." + url.slice(-4);
  }

  const [protocol, rest] = parts;
  return `${protocol}://${rest.slice(0, 4)}...${rest.slice(-4)}`;
};

onMounted(async () => {
  await fetchSenders();
});
</script>

<style scoped>
code {
  background: var(--bs-dark-bg-subtle);
  border-radius: 4px;
}
.bi-info-circle {
  font-size: 0.875rem;
  color: #6c757d;
}
</style>
