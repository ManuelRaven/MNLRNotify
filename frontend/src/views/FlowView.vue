<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type {
  ChannelResponse,
  RecieverResponse,
  SenderResponse,
} from "@/types/pocketbase-types";
import { Background } from "@vue-flow/background";
import { useVueFlow, VueFlow } from "@vue-flow/core";
import { useToastController } from "bootstrap-vue-next";
import { onMounted } from "vue";

const pb = usePb();
const toast = useToastController();
const { nodes, edges, addNodes, addEdges, setViewport } = useVueFlow();

// Fetch data and create flow
const createFlow = async () => {
  try {
    const [receivers, channels, senders] = await Promise.all([
      pb.collection("reciever").getList<RecieverResponse>(1, 100),
      pb.collection("channel").getList<ChannelResponse>(1, 100),
      pb.collection("sender").getList<SenderResponse>(1, 100),
    ]);

    // Create nodes
    const receiverNodes = receivers.items.map((receiver, index) => ({
      id: `receiver-${receiver.id}`,
      position: { x: 100 + index * 200, y: 100 }, // Changed y to fixed, x varies with index
      data: { label: receiver.name },
      type: "input",
      class: "receiver-node",
    }));

    const channelNodes = channels.items.map((channel, index) => ({
      id: `channel-${channel.id}`,
      position: { x: 100 + index * 200, y: 400 }, // Changed y to middle position
      data: { label: channel.name },
      class: "channel-node",
    }));

    const senderNodes = senders.items.map((sender, index) => ({
      id: `sender-${sender.id}`,
      position: { x: 100 + index * 200, y: 700 }, // Changed y to bottom position
      data: { label: sender.name },
      type: "output",
      class: "sender-node",
    }));

    addNodes([...receiverNodes, ...channelNodes, ...senderNodes]);

    // Create edges
    const receiverEdges = receivers.items.flatMap((receiver) =>
      (receiver.channels || []).map((channelId) => ({
        id: `e-${receiver.id}-${channelId}`,
        source: `receiver-${receiver.id}`,
        target: `channel-${channelId}`,
        animated: true,
        class: "receiver-edge",
      }))
    );

    const senderEdges = senders.items.flatMap((sender) =>
      (sender.channel || []).map((channelId) => ({
        id: `e-${channelId}-${sender.id}`,
        source: `channel-${channelId}`,
        target: `sender-${sender.id}`,
        animated: true,
        class: "sender-edge",
      }))
    );

    addEdges([...receiverEdges, ...senderEdges]);

    // Center view
    setViewport({ x: 0, y: 0, zoom: 0.8 });
  } catch (error) {
    toast.show?.({
      props: {
        title: "Error",
        body: "Failed to load flow data",
        variant: "danger",
      },
    });
  }
};

onMounted(() => {
  createFlow();
});
</script>

<template>
  <div class="vueflow-page-container">
    <VueFlow
      :default-viewport="{ zoom: 0.8 }"
      :min-zoom="0.2"
      :max-zoom="4"
      fit-view-on-init
    >
      <Background pattern-color="#aaa" :gap="16" />
    </VueFlow>
  </div>
</template>

<style>
@import "@vue-flow/core/dist/style.css";
@import "@vue-flow/core/dist/theme-default.css";

@media (max-width: 768px) {
  .vueflow-page-container {
    height: calc(100vh - 76px);
  }
}

@media (min-width: 768px) {
  .vueflow-page-container {
    height: calc(100vh - 56px);
  }
}

.vueflow-page-container {
  display: flex;
  flex-direction: column;
}

/* Node styling */
.receiver-node {
  background-color: var(--bs-primary);
  color: white;
}

.channel-node {
  background-color: var(--bs-secondary);
  color: white;
}

.sender-node {
  background-color: var(--bs-success);
  color: white;
}

/* Edge styling */
.receiver-edge {
  stroke: var(--bs-primary);
}

.sender-edge {
  stroke: var(--bs-success);
}

.vue-flow__node {
  padding: 10px;
  border-radius: 5px;
  font-weight: bold;
  text-align: center;
  width: 150px;
}
</style>
