<script setup lang="ts">
import { onMounted } from "vue";
import { RouterView } from "vue-router";
import { useAuth } from "./composeables/useAuth";
import { usePb } from "./composeables/usePb";

const pb = usePb();
const auth = useAuth();

onMounted(async () => {
  let health = await pb.health.check();
  console.log(health);
});
</script>

<template>
  <BModalOrchestrator />
  <BToastOrchestrator />
  <RouterView />
</template>

<style scoped></style>
