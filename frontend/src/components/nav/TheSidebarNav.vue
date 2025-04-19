<template>
  <BNav vertical pills>
    <BNavItem>
      <BFormSelect v-model="mode">
        <option value="light">Light</option>
        <option value="dark">Dark</option>
      </BFormSelect>
    </BNavItem>
    <BNavItem>
      <BFormSelect v-model="selectedTheme" @change="changeTheme">
        <option v-for="theme in availableThemes" :key="theme" :value="theme">
          {{ theme.charAt(0).toUpperCase() + theme.slice(1) }}
        </option>
      </BFormSelect>
    </BNavItem>
    <BNavItem exact exact-active-class="active" :to="{ name: 'home' }"
      >Home</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'channels' }"
      >Channels</BNavItem
    >
    <BNavItemDropdown
      exact
      exact-active-class="active"
      :to="{ name: 'channel-mutations' }"
    >
      <template #button-content> Channel Mutations </template>
      <BDropdownItem
        exact-active-class="active"
        :to="{ name: 'channel-mutations' }"
      >
        List
      </BDropdownItem>
      <BDropdownItem
        exact-active-class="active"
        :to="{ name: 'channel-mutations-editor' }"
      >
        Editor
      </BDropdownItem>
    </BNavItemDropdown>
    <BNavItem exact exact-active-class="active" :to="{ name: 'senders' }"
      >Senders</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'webpush' }"
      >WebPush</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'messages' }"
      >Messages</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'chat' }"
      >Messages (Chat)</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'receivers' }"
      >Receivers</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'guides' }"
      >Integration Guides</BNavItem
    >
    <BNavItem exact exact-active-class="active" :to="{ name: 'flows' }"
      >Flow Diagram</BNavItem
    >

    <BNavItemDropdown right>
      <!-- Using 'button-content' slot -->
      <template #button-content>
        <em>{{ auth.authModel.value?.email }}</em>
      </template>
      <!-- <BDropdownItem exact-active-class="active" :to="{ name: 'settings' }"
        >Settings</BDropdownItem
      > -->
      <BDropdownItem @click="logout">Sign Out</BDropdownItem>
    </BNavItemDropdown>
    <BNavItem disabled>v{{ version }}</BNavItem>
  </BNav>
</template>

<script setup lang="ts">
import { useAuth } from "@/composeables/useAuth";
import { BFormSelect, useColorMode } from "bootstrap-vue-next";
import { ref, watchEffect } from "vue";
import { useRouter } from "vue-router";
import pkg from "../../../../package.json";

const version = pkg.version;

const auth = useAuth();
const router = useRouter();

const darkEnabled = ref(false);

const mode = useColorMode({
  persist: true,
});

import { useThemeStore, availableThemes } from "@/stores/useThemeStore";

const themeStore = useThemeStore();
const selectedTheme = ref(themeStore.theme);

watchEffect(() => {
  darkEnabled.value = mode.value === "dark";
});

const changeColor = () => {
  mode.value = mode.value === "dark" ? "light" : "dark";
};

const changeTheme = () => {
  console.log("Changing theme to:", selectedTheme.value);

  themeStore.changeTheme(selectedTheme.value);
};

const logout = async () => {
  await auth.logout();
  await router.replace("/login");
};
</script>

<style scoped></style>
