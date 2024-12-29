<template>
  <BCard title="Telegram Chat ID Helper">
    <BForm @submit.prevent="startPolling" v-if="!isPolling">
      <BFormGroup label="Bot Token" description="Enter your Telegram bot token">
        <BFormInput
          v-model="botToken"
          placeholder="123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
          required
        />
      </BFormGroup>
      <BButton type="submit" variant="primary" :disabled="!botToken">
        Start Listening
      </BButton>
    </BForm>

    <div v-else>
      <BAlert :model-value="true" show variant="info">
        <div class="d-flex justify-content-between align-items-center">
          <span>Listening for messages... Send a message to your bot.</span>
          <BButton size="sm" variant="outline-danger" @click="stopPolling"
            >Stop</BButton
          >
        </div>
      </BAlert>

      <div v-if="updates.length > 0">
        <h4>Received Messages:</h4>
        <BListGroup>
          <BListGroupItem v-for="update in updates" :key="update.update_id">
            <strong>Chat ID:</strong> {{ update.message.chat.id }}<br />
            <strong>From:</strong> {{ update.message.from.first_name }}<br />
            <strong>Message:</strong> {{ update.message.text }}<br />
            <div class="mt-2">
              <strong>Shoutrrr URL:</strong>
              <div class="d-flex align-items-center gap-2">
                <code class="flex-grow-1">{{
                  constructShoutrrrUrl(update.message.chat.id)
                }}</code>
                <BButton
                  size="sm"
                  variant="outline-primary"
                  @click="copyUrl(update.message.chat.id)"
                >
                  Copy URL
                </BButton>
              </div>
            </div>
          </BListGroupItem>
        </BListGroup>
      </div>

      <BAlert v-if="error" show variant="danger">
        {{ error }}
      </BAlert>
    </div>

    <!-- Zeige Nachrichten auch nach dem Stoppen an -->
    <div v-if="updates.length > 0" class="mt-3">
      <h4>Received Messages:</h4>
      <BListGroup>
        <BListGroupItem v-for="update in updates" :key="update.update_id">
          <strong>Chat ID:</strong> {{ update.message.chat.id }}<br />
          <strong>From:</strong> {{ update.message.from.first_name }}<br />
          <strong>Message:</strong> {{ update.message.text }}<br />
          <div class="mt-2">
            <strong>Shoutrrr URL:</strong>
            <div class="d-flex align-items-center gap-2">
              <code class="flex-grow-1">{{
                constructShoutrrrUrl(update.message.chat.id)
              }}</code>
              <BButton
                size="sm"
                variant="outline-primary"
                @click="copyUrl(update.message.chat.id)"
              >
                Copy URL
              </BButton>
            </div>
          </div>
        </BListGroupItem>
      </BListGroup>
    </div>
  </BCard>
</template>

<script setup lang="ts">
import { onUnmounted, ref } from "vue";

const botToken = ref("");
const isPolling = ref(false);
const updates = ref<any[]>([]);
const error = ref("");
let pollInterval: number | null = null;
const lastUpdateId = ref(0);
const activeController = ref<AbortController | null>(null);

const fetchWithTimeout = async (url: string) => {
  if (activeController.value) {
    activeController.value.abort();
  }
  activeController.value = new AbortController();

  try {
    const response = await fetch(url, {
      signal: activeController.value.signal,
    });
    return await response.json();
  } catch (e) {
    if (e instanceof DOMException && e.name === "AbortError") {
      throw new Error("Request was cancelled");
    }
    throw e;
  }
};

const fetchUpdates = async () => {
  try {
    const data = await fetchWithTimeout(
      `https://api.telegram.org/bot${botToken.value}/getUpdates?offset=${lastUpdateId.value}&timeout=3`
    );

    if (!data.ok) {
      throw new Error(data.description || "Failed to fetch updates");
    }

    if (data.result.length > 0) {
      updates.value.push(...data.result);
      lastUpdateId.value = data.result[data.result.length - 1].update_id + 1;
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : "An error occurred";
    await cancelPolling();
  }
};

const cancelPolling = async () => {
  if (activeController.value) {
    activeController.value.abort();
    activeController.value = null;
  }
  stopPolling();
  try {
    // Final cleanup request
    await fetch(
      `https://api.telegram.org/bot${botToken.value}/getUpdates?offset=-1`
    );
  } catch {
    // Ignore cleanup errors
  }
};

const startPolling = async () => {
  try {
    isPolling.value = true;
    error.value = "";
    updates.value = [];

    const data = await fetchWithTimeout(
      `https://api.telegram.org/bot${botToken.value}/getUpdates?limit=1`
    );

    if (!data.ok) {
      throw new Error(data.description || "Failed to initialize polling");
    }

    if (data.result.length > 0) {
      lastUpdateId.value = data.result[0].update_id + 1;
    }

    pollInterval = window.setInterval(fetchUpdates, 3500);
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Failed to start polling";
    await cancelPolling();
  }
};

const stopPolling = () => {
  if (pollInterval) {
    clearInterval(pollInterval);
    pollInterval = null;
  }
  isPolling.value = false;
  // Entferne das Löschen der updates
  // updates.value = [];
};

const constructShoutrrrUrl = (chatId: string): string => {
  return `telegram://${botToken.value}@telegram?chats=${chatId}`;
};

const copyUrl = async (chatId: string) => {
  try {
    await navigator.clipboard.writeText(constructShoutrrrUrl(chatId));
    // Optional: Zeige Feedback an (wenn Sie einen Toast-Controller haben)
    // toast.show?.({ props: { title: 'Success', body: 'URL copied to clipboard', variant: 'success' }});
  } catch (err) {
    console.error("Failed to copy:", err);
  }
};

// Cleanup on component unmount
onUnmounted(() => {
  cancelPolling();
});
</script>

<style scoped>
.list-group {
  margin-top: 1rem;
}

code {
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
  word-break: break-all;
}
</style>
