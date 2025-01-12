import vue from "@vitejs/plugin-vue";
import { BootstrapVueNextResolver } from "bootstrap-vue-next";
import path from "path";
import IconsResolve from "unplugin-icons/resolver";
import Icons from "unplugin-icons/vite";
import Components from "unplugin-vue-components/vite";
import { defineConfig, loadEnv } from "vite";
import { VitePWA } from "vite-plugin-pwa";
import vueDevTools from "vite-plugin-vue-devtools";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const backendUrl = env.VITE_BACKEND_URL || "http://localhost:8090";
  return {
    plugins: [
      vue(),
      Components({
        dts: true,
        resolvers: [BootstrapVueNextResolver(), IconsResolve()],
      }),
      Icons({
        compiler: "vue3",
        autoInstall: true,
      }),
      vueDevTools(),
      VitePWA({
        srcDir: "frontend/src",
        filename: "service-worker.ts",
        strategies: "injectManifest",
        registerType: "autoUpdate",
        devOptions: {
          type: "module",
          enabled: true,
        },
        injectRegister: "auto",
        manifest: {
          name: "MNLRNotify",
          short_name: "MNLRNotify",
          description: "MNLRNotify",
          theme_color: "#ffffff",
          icons: [
            {
              src: "pwa-64x64.png",
              sizes: "64x64",
              type: "image/png",
            },
            {
              src: "pwa-192x192.png",
              sizes: "192x192",
              type: "image/png",
            },
            {
              src: "pwa-512x512.png",
              sizes: "512x512",
              type: "image/png",
            },
          ],
        },
      }),
    ],
    build: {
      outDir: "./backend/dist",
      emptyOutDir: true,
    },
    publicDir: "./frontend/public",
    resolve: { alias: { "@": path.resolve(__dirname, "./frontend/src") } },
    server: {
      proxy: {
        "/api": {
          target: backendUrl,
          changeOrigin: true,
        },
        "/eapi": {
          target: backendUrl,
          changeOrigin: true,
        },
      },
    },
  };
});
