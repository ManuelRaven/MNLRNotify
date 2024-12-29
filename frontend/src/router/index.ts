import { useAuth } from "@/composeables/useAuth";
import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/system/HomeView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: () => import("../layouts/MainLayout.vue"),
      children: [
        {
          path: "",
          name: "home",
          component: HomeView,
        },
        {
          path: "/about",
          name: "about",
          // route level code-splitting
          // this generates a separate chunk (About.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () => import("@/views/system/AboutView.vue"),
        },
        {
          path: "/testing",
          name: "testing",
          component: () => import("@/views/system/TestingView.vue"),
        },
        {
          path: "/flows",
          name: "flows",
          component: () => import("@/views/FlowView.vue"),
        },
        {
          path: "/channels",
          name: "channels",
          component: () => import("@/views/ChannelView.vue"),
        },
        {
          path: "/messagesView",
          name: "messages",
          component: () => import("@/views/MessageView.vue"),
        },
        {
          path: "/senders",
          name: "senders",
          component: () => import("@/views/SenderView.vue"),
        },
        {
          path: "/recievers",
          name: "recievers",
          component: () => import("@/views/RecieverView.vue"),
        },
        {
          path: "/guides",
          name: "guides",
          component: () => import("@/views/guides/GuidesView.vue"),
        },
        {
          path: "/guides/proxmox",
          name: "guides-proxmox",
          component: () => import("@/views/guides/ProxmoxGuide.vue"),
        },
        {
          path: "/guides/telegram",
          name: "guides-proxmox",
          component: () => import("@/views/guides/TelegramGuide.vue"),
        },
        {
          path: "/chat",
          name: "chat",
          component: () => import("../views/MessageChatView.vue"),
        },
      ],
    },
    {
      path: "/",
      component: () => import("../layouts/AnonymousLayout.vue"),
      children: [
        {
          path: "/login",
          name: "login",
          // route level code-splitting
          // this generates a separate chunk (About.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () => import("@/views/system/LoginView.vue"),
        },
        {
          path: "/register",
          name: "register",
          component: () => import("@/views/system/RegisterView.vue"),
        },
      ],
    },
  ],
});

//RouteGuard to check if user is authenticated and redirect to login page if not

const unauthorizedRoutes = ["login", "register"];

router.beforeEach(async (to, from, next) => {
  const isAuthenticated = await useAuth().getAuthState();
  if (!isAuthenticated && !unauthorizedRoutes.includes(to.name as string)) {
    next({ name: "login" });
  } else next();
});

export default router;
