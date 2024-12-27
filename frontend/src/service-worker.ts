declare let self: ServiceWorkerGlobalScope;

//precacheAndRoute(self.__WB_MANIFEST);

import { clientsClaim } from "workbox-core";

self.skipWaiting();
clientsClaim();

self.addEventListener("fetch", (event) => {
  console.log("fetch event", event.request.url);
  event.respondWith(fetch(event.request));
});
