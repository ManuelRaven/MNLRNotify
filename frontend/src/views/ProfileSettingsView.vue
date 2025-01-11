<template>
  <BCard title="Push Notifications">
    <BAlert
      :model-value="true"
      variant="warning"
      :dismissible="false"
      class="mb-3"
    >
      <p>
        This is only available if the WebApp is served over HTTPS and the
        browser supports Push Notifications. On IOS devices, you need to use
        this app as a PWA.
      </p>
    </BAlert>
    <BFormCheckbox v-model="pushNotificationsEnabled">
      Enable Push Notifications
    </BFormCheckbox>
  </BCard>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

const pushNotificationsEnabled = ref(false);

watch(
  () => pushNotificationsEnabled.value,
  (enabled) => {
    if (enabled) {
      enablePushNotifications();
    } else {
      disablePushNotifications();
    }
  }
);

const disablePushNotifications = async () => {
  if (!("serviceWorker" in navigator)) {
    console.error("Service worker not supported");
    return;
  }

  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();

  if (!subscription) {
    console.log("Not subscribed");
    return;
  }

  await subscription.unsubscribe();
  console.log("Unsubscribed");
};

// Enable PWA Push Notifications, service worker are already registered in main.js

const registerPushNotifications = async () => {
  // Use serviceWorker.ready to ensure that you can subscribe for push
  navigator.serviceWorker.ready.then((serviceWorkerRegistration) => {
    const options = {
      userVisibleOnly: true,
      applicationServerKey:
        "BPast68hR3ihXdYIAtbkM1Rv-uPJd3PTXljuUYUq-4GDXC5sZqLPIAVO20cHA05pZEc22B1iKyQaHuTEmZZC2vI",
    };
    serviceWorkerRegistration.pushManager.subscribe(options).then(
      (pushSubscription) => {
        console.log(pushSubscription.endpoint);
        // The push subscription details needed by the application
        // server are now available, and can be sent to it using,
        // for example, the fetch() API.
      },
      (error) => {
        // During development it often helps to log errors to the
        // console. In a production environment it might make sense to
        // also report information about errors back to the
        // application server.
        console.error(error);
      }
    );
  });
};

// Enable PWA Push Notifications
const enablePushNotifications = async () => {
  if (!("serviceWorker" in navigator)) {
    console.error("Service worker not supported");
    return;
  }

  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();

  console.log("Subscription", subscription);

  if (subscription) {
    console.log("Already subscribed");
    return;
  }

  await registerPushNotifications();
};
</script>

<style scoped></style>
