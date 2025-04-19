import "bootstrap-vue-next/dist/bootstrap-vue-next.css";
import "bootstrap/dist/css/bootstrap.css";
import { createPinia } from "pinia";
import { createApp } from "vue";
import "./assets/main.css";

import { createBootstrap } from "bootstrap-vue-next";
import App from "./App.vue";
import router from "./router";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(createBootstrap()); // Important
app.mount("#app");

if ("serviceWorker" in navigator) {
  navigator.serviceWorker.register(
    (import.meta as any).env.MODE === "production"
      ? "/service-worker.js"
      : "/dev-sw.js?dev-sw",
    {
      type:
        (import.meta as any).env.MODE === "production" ? "classic" : "module",
    }
  );
}
