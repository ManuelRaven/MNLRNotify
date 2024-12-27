declare let self: ServiceWorkerGlobalScope;

//precacheAndRoute(self.__WB_MANIFEST);

import { clientsClaim } from "workbox-core";

self.skipWaiting();
clientsClaim();
