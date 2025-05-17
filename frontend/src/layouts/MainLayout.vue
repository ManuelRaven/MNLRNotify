<script setup lang="ts">
import TheFrame from "@/components/nav/TheFrame.vue";
</script>

<template>
  <TheFrame
    ><router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component"></component>
      </keep-alive>
    </router-view>
  </TheFrame>
</template>

<style scoped></style>
