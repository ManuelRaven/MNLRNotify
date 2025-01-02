<template>
  <div class="chat-container">
    <!-- Channel Sidebar -->
    <div
      class="channel-sidebar p-2 border-end overflow-y-auto"
      :class="{ 'sidebar-hidden': !showSidebar }"
    >
      <h5 class="mb-3">Channels</h5>
      <div class="channel-list">
        <div
          v-for="channel in channels"
          :key="channel.id"
          class="channel-item p-2 mb-1"
          :class="{ active: selectedChannel?.id === channel.id }"
          @click="selectChannel(channel)"
        >
          {{ channel.name }}
        </div>
      </div>
    </div>

    <!-- Message Area -->
    <div class="message-area">
      <div class="chat-header p-3 border-bottom d-flex align-items-center">
        <BButton
          class="me-2 d-md-none"
          @click="toggleSidebar"
          variant="outline-primary"
          size="sm"
        >
          <IBiList />
        </BButton>
        <h4 class="mb-0">{{ selectedChannel?.name || "Select a channel" }}</h4>
      </div>
      <div class="message-list p-3" ref="messageContainer">
        <div v-if="!selectedChannel" class="text-center text-muted mt-5">
          Select a channel to view messages
        </div>
        <template v-else>
          <div class="messages-wrapper">
            <div
              v-for="message in [...messages].reverse()"
              :key="message.id"
              class="message-item mb-3 p-2"
              :class="getMessageStateClass(message.deliveryState)"
            >
              <div class="message-time text-muted small">
                {{ new Date(message.created).toLocaleString() }}
              </div>
              <div class="message-text">{{ message.text }}</div>
              <BAccordion v-if="message.deliveryMessage" class="mt-2">
                <BAccordionItem title="Delivery Details">
                  <div class="delivery-message text-break pre-wrap">
                    {{ message.deliveryMessage }}
                  </div>
                </BAccordionItem>
              </BAccordion>
            </div>
          </div>
        </template>
      </div>

      <!-- Message Input Area -->
      <div class="message-input-area p-3 border-top" v-if="selectedChannel">
        <form @submit.prevent="sendMessage" class="d-flex gap-2">
          <BFormTextarea
            v-model="newMessage"
            placeholder="Type a message..."
            rows="1"
            max-rows="4"
            class="flex-grow-1"
            :disabled="!selectedChannel"
            @keydown.enter.exact.prevent="sendMessage"
          />
          <BButton
            type="submit"
            variant="primary"
            :disabled="!newMessage.trim()"
          >
            Send
          </BButton>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type {
  ChannelResponse,
  MessageDeliveryStateOptions,
  MessageResponse,
} from "@/types/pocketbase-types";
import { useToastController } from "bootstrap-vue-next";
import { nextTick, onMounted, ref } from "vue";

const pb = usePb();
const toast = useToastController();

const channels = ref<ChannelResponse[]>([]);
const messages = ref<MessageResponse[]>([]);
const selectedChannel = ref<ChannelResponse | null>(null);
const newMessage = ref("");
const messageContainer = ref<HTMLElement | null>(null);
const showSidebar = ref(true);

const fetchChannels = async () => {
  try {
    const response = await pb.collection("channel").getList(1, 50);
    channels.value = response.items;
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

const getMessageStateClass = (state?: MessageDeliveryStateOptions) => {
  switch (state) {
    case "success":
      return "message-success";
    case "failure":
      return "message-failure";
    default:
      return "message-pending";
  }
};

const scrollToBottom = () => {
  if (messageContainer.value) {
    messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
  }
};

const fetchMessages = async (channelId: string) => {
  try {
    const response = await pb.collection("message").getList(1, 100, {
      filter: `channel = "${channelId}"`,
      sort: "-created",
    });
    messages.value = response.items;
    // Scroll to bottom after messages load
    nextTick(() => {
      scrollToBottom();
    });
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

const selectChannel = (channel: ChannelResponse) => {
  selectedChannel.value = channel;
  fetchMessages(channel.id);
};

const sendMessage = async () => {
  if (!selectedChannel.value || !newMessage.value.trim()) return;

  try {
    const message = await pb.collection("message").create({
      text: newMessage.value,
      channel: selectedChannel.value.id,
    });

    messages.value.unshift(message);
    newMessage.value = "";
    // Scroll to bottom after sending
    nextTick(() => {
      scrollToBottom();
    });
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to send message",
        variant: "danger",
      },
    });
  }
};

const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value;
};

onMounted(async () => {
  await fetchChannels();
});
</script>

<style scoped>
.chat-container {
  display: flex;
  height: calc(100vh - 76px); /* Adjust based on your navbar height */
  overflow: hidden;
}

.channel-sidebar {
  width: 250px;
  background-color: var(--bs-dark-bg-subtle);
  height: 100%;
  transition: transform 0.3s ease;
}

@media (max-width: 767.98px) {
  .channel-sidebar {
    position: fixed;
    left: 0;
    top: 64px; /* Adjust based on your navbar height */
    z-index: 1030;
    height: calc(100vh - 56px);
  }

  .sidebar-hidden {
    transform: translateX(-100%);
  }
}

.message-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.chat-header {
  background-color: var(--bs-body-bg);
  z-index: 1;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  background-color: var(--bs-body-bg);
  display: flex;
  flex-direction: column;
}

.channel-item {
  cursor: pointer;
  border-radius: 4px;
  color: var(--bs-body-color);
}

.channel-item:hover {
  background-color: var(--bs-tertiary-bg);
}

.channel-item.active {
  background-color: var(--bs-primary);
  color: var (--bs-light);
}

.message-container {
  overflow-y: auto;
  background-color: var(--bs-body-bg);
}

.message-item {
  background-color: var(--bs-dark-bg-subtle);
  border-radius: 4px;
  max-width: 80%;
  color: var(--bs-body-color);
  margin-left: auto;
  border-width: 2px;
  border-style: solid;
}

.message-success {
  border-color: var(--bs-success);
}

.message-failure {
  border-color: var(--bs-danger);
}

.message-pending {
  border-color: var(--bs-warning);
}

.delivery-message {
  font-family: monospace;
  font-size: 0.9em;
  padding: 0.5rem;
  background-color: var(--bs-tertiary-bg);
  border-radius: 4px;
  white-space: pre-wrap;
}

.message-text {
  word-break: break-word;
  white-space: pre-wrap;
}

.message-time {
  color: var(--bs-secondary-color) !important;
}

/* Add border color for the header */
.border-bottom {
  border-bottom-color: var(--bs-border-color) !important;
}

h4,
h5 {
  color: var(--bs-heading-color);
}

.message-input-area {
  background-color: var(--bs-body-bg);
  border-top-color: var(--bs-border-color) !important;
}

.message-input-area textarea {
  resize: none;
  background-color: var(--bs-dark-bg-subtle);
  border-color: var(--bs-border-color);
  color: var(--bs-body-color);
}

.message-input-area textarea:focus {
  background-color: var(--bs-tertiary-bg);
  border-color: var(--bs-primary);
  box-shadow: 0 0 0 0.25rem rgb(var(--bs-primary-rgb), 0.25);
}

.messages-wrapper {
  display: flex;
  flex-direction: column;
  min-height: min-content;
  justify-content: flex-end;
  margin-top: auto;
}
</style>
