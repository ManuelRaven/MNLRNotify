declare let self: ServiceWorkerGlobalScope;

import { cleanupOutdatedCaches, precacheAndRoute } from "workbox-precaching";

cleanupOutdatedCaches();

precacheAndRoute(self.__WB_MANIFEST);

interface PushMessageFormat {
  channelName: string;
  message: string;
}

import { clientsClaim } from "workbox-core";

self.skipWaiting();
clientsClaim();

self.addEventListener("push", (event) => {
  const data = event.data?.json() as PushMessageFormat;
  const options: NotificationOptions = {
    body: data.message,
  };

  event.waitUntil(
    self.registration.showNotification(
      `MNLRNotify: ${data.channelName}`,
      options
    )
  );
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  event.waitUntil(self.clients.openWindow("/messagesView"));
});
