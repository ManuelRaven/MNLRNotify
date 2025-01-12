<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type {
  ChannelResponse,
  RecieverResponse,
  SenderResponse,
  WebpushDevicesResponse,
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
    const [receivers, channels, senders, webpushDevices] = await Promise.all([
      pb.collection("reciever").getList<RecieverResponse>(1, 100),
      pb.collection("channel").getList<ChannelResponse>(1, 100),
      pb.collection("sender").getList<SenderResponse>(1, 100),
      pb.collection("webpushDevices").getList<WebpushDevicesResponse>(1, 100),
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

    // Calculate positions for bottom row (senders and webpush)
    const bottomY = 700;
    const spacing = 180; // Fixed smaller spacing instead of window-based calculation

    const webpushNodes = webpushDevices.items.map((device, index) => ({
      id: `webpush-${device.id}`,
      position: { x: 100 + index * spacing, y: bottomY },
      data: { label: "📱 " + device.Name },
      type: "output",
      class: "webpush-node",
    }));

    const senderNodes = senders.items.map((sender, index) => ({
      id: `sender-${sender.id}`,
      position: {
        x: 100 + (index + webpushDevices.items.length) * spacing,
        y: bottomY,
      },
      data: { label: sender.name },
      type: "output",
      class: "sender-node",
    }));

    addNodes([
      ...receiverNodes,
      ...channelNodes,
      ...webpushNodes,
      ...senderNodes,
    ]);

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

    const webpushEdges = webpushDevices.items.flatMap((device) =>
      (device.channel || []).map((channelId) => ({
        id: `e-${channelId}-${device.id}`,
        source: `channel-${channelId}`,
        target: `webpush-${device.id}`,
        animated: true,
        class: "webpush-edge",
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

    addEdges([...receiverEdges, ...webpushEdges, ...senderEdges]);

    // Center view
    setViewport({ x: 0, y: 0, zoom: 0.7 }); // Adjusted zoom to show all nodes
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

.webpush-node {
  background-color: var(--bs-info);
  color: white;
}

/* Edge styling */
.receiver-edge {
  stroke: var(--bs-primary);
}

.sender-edge {
  stroke: var(--bs-success);
}

.webpush-edge {
  stroke: var(--bs-info);
}

.vue-flow__node {
  padding: 10px;
  border-radius: 5px;
  font-weight: bold;
  text-align: center;
  width: 150px;
}
</style>
