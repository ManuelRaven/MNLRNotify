<template>
  <div class="guides-container">
    <h1>Integration Guides</h1>
    <p class="guides-intro">
      Learn how to integrate MNLRNotify with popular applications and services.
    </p>

    <div class="guides-grid">
      <BCard
        v-for="guide in guides"
        :key="guide.path"
        :title="guide.title"
        class="guide-card"
        @click="router.push(guide.path)"
      >
        <BCardText>{{ guide.description }}</BCardText>
        <template #footer>
          <small class="text-muted">Difficulty: {{ guide.difficulty }}</small>
        </template>
      </BCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();

const guides = [
  {
    title: "Proxmox Integration",
    description: "Configure Proxmox to send notifications through MNLRNotify",
    difficulty: "Intermediate",
    path: "/guides/proxmox",
  },
  {
    title: "Telegram Bot Setup",
    description: "Learn how to create a Telegram bot and get your bot token",
    difficulty: "Beginner",
    path: "/guides/telegram",
  },
  // More guides can be added here
];
</script>

<style scoped>
.guides-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.guides-intro {
  margin-bottom: 2rem;
  color: var(--bs-secondary-color);
}

.guides-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.guide-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.guide-card:hover {
  transform: translateY(-5px);
}
</style>
