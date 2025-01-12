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
    <BFormGroup label="Device Name" label-for="deviceName" class="mb-3">
      <BFormInput
        id="deviceName"
        v-model="deviceName"
        placeholder="Enter device name"
        :state="deviceNameValid"
        :disabled="isSubscribed"
      />
      <BFormText v-if="!deviceNameValid && !isSubscribed">
        {{ deviceNameErrorMessage }}
      </BFormText>
    </BFormGroup>
    <div class="d-grid gap-2 mb-4">
      <BAlert
        v-if="isRegistering"
        :model-value="true"
        variant="info"
        :dismissible="false"
        class="mb-3"
      >
        <p>
          Your browser may ask for permission to enable notifications. Please
          check your browser's address bar or notification toggle prompt.
        </p>
      </BAlert>
      <BButton
        v-if="!isSubscribed"
        variant="primary"
        :disabled="!deviceNameValid || isRegistering"
        @click="registerDevice"
      >
        <BSpinner v-if="isRegistering" small class="me-2"></BSpinner>
        {{ isRegistering ? "Registering..." : "Register Device" }}
      </BButton>
      <BButton v-else variant="danger" @click="unregisterDevice">
        Unregister This Device
      </BButton>
    </div>

    <!-- Device List -->
    <div v-if="devices.length > 0" class="mt-4">
      <h5>Registered Devices</h5>
      <BListGroup>
        <BListGroupItem
          v-for="device in devices"
          :key="device.id"
          class="d-flex justify-content-between align-items-center"
        >
          <div>
            {{ device.Name }}
            <BBadge v-if="device.id === currentDevice?.id" class="ms-2">
              This Device
            </BBadge>
            <div class="text-muted small">
              {{
                device.expand?.channel?.map((c) => c.name).join(", ") ||
                "No channels"
              }}
            </div>
          </div>
          <div>
            <BButton
              size="sm"
              variant="info"
              class="me-2"
              @click="openChannelModal(device)"
            >
              Channels
            </BButton>
            <BButton
              size="sm"
              variant="danger"
              @click="removeDevice(device.id)"
              :disabled="device.id === currentDevice?.id"
            >
              Remove
            </BButton>
          </div>
        </BListGroupItem>
      </BListGroup>
    </div>

    <BAlert
      v-if="errorMessage"
      :model-value="true"
      variant="danger"
      dismissible
      @dismissed="errorMessage = ''"
    >
      {{ errorMessage }}
    </BAlert>
  </BCard>
  <ChannelSelectionModal
    v-model="showChannelModal"
    :preset-channels="currentManageDevice?.channel"
    @save="updateDeviceChannels"
  />
</template>

<script setup lang="ts">
import ChannelSelectionModal from "@/components/GenericModals/ChannelSelectionModal.vue";
import { usePb } from "@/composeables/usePb";
import type {
  ExpandChannelNameId,
  WebpushDeviceRequest,
} from "@/types/custom-types";
import type { WebpushDevicesResponse } from "@/types/pocketbase-types";
import { useToastController } from "bootstrap-vue-next";
import { computed, onMounted, ref } from "vue";

const pb = usePb();
const vapidPublicKey = ref<string | null>(null);
const deviceName = ref("");
const errorMessage = ref("");
const devices = ref<
  WebpushDevicesResponse<PushSubscription, ExpandChannelNameId>[]
>([]);
const isSubscribed = ref(false);
const showChannelModal = ref(false);
const currentManageDevice = ref<WebpushDevicesResponse | null>(null);
const isRegistering = ref(false);

const currentDevice = computed(() => {
  if (!currentEndpoint.value) return null;
  return devices.value.find(
    (device) => device.subscription?.endpoint === currentEndpoint.value
  );
});
const toast = useToastController();
const currentEndpoint = ref<string | null>(null);

onMounted(async () => {
  await initializeDevices();
  let resp = await pb.send<string>("/eapi/vapid", { method: "GET" });
  vapidPublicKey.value = resp;
});

const findNextAvailableName = (baseName: string): string => {
  if (
    !devices.value.some(
      (device) => device.Name.toLowerCase() === baseName.toLowerCase()
    )
  ) {
    return baseName;
  }

  let counter = 2;
  while (
    devices.value.some(
      (device) =>
        device.Name.toLowerCase() === `${baseName} ${counter}`.toLowerCase()
    )
  ) {
    counter++;
  }
  return `${baseName} ${counter}`;
};

const generateDeviceName = (): string => {
  const ua = navigator.userAgent;
  let deviceName = "Unknown Device";

  // Detect browser
  const browserMatches = ua.match(
    /(chrome|safari|firefox|edge|opera(?=\/))\/?\s*(\d+)/i
  );
  const browser = browserMatches ? browserMatches[1] : "";

  // Detect OS
  const osMatches = ua.match(/windows nt|mac os x|linux|android|ios/i);
  const os = osMatches ? osMatches[0] : "";

  if (browser && os) {
    deviceName = `${browser.charAt(0).toUpperCase() + browser.slice(1)} on ${
      os.charAt(0).toUpperCase() + os.slice(1)
    }`;
  } else if (browser) {
    deviceName = browser.charAt(0).toUpperCase() + browser.slice(1);
  } else if (os) {
    deviceName = os.charAt(0).toUpperCase() + os.slice(1);
  }

  return findNextAvailableName(deviceName);
};

const initializeDevices = async () => {
  await updateCurrentEndpoint();
  await loadDevices();
  isSubscribed.value = !!currentDevice.value;
  if (currentDevice.value) {
    deviceName.value = currentDevice.value.Name;
  } else {
    deviceName.value = generateDeviceName();
  }
};

const updateCurrentEndpoint = async () => {
  if (!("serviceWorker" in navigator)) return;
  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();
  currentEndpoint.value = subscription?.endpoint || null;
};

const loadDevices = async () => {
  try {
    const result = await pb
      .collection("webpushDevices")
      .getFullList<
        WebpushDevicesResponse<PushSubscription, ExpandChannelNameId>
      >({
        filter: pb.authStore.record?.id
          ? `owner = "${pb.authStore.record?.id}"`
          : "",
        expand: "channel",
      });
    devices.value = result;
  } catch (error) {
    console.error("Failed to load devices:", error);
  }
};

const registerDeviceInDatabase = async (subscription: PushSubscription) => {
  try {
    const data: WebpushDeviceRequest = {
      Name: deviceName.value,
      subscription: subscription.toJSON(),
      owner: "",
    };

    await pb.collection("webpushDevices").create(data);
    await loadDevices();
    await updateCurrentEndpoint();
  } catch (error) {
    console.error("Failed to register device:", error);
    throw new Error("Failed to register device in database");
  }
};

const removeDevice = async (deviceId: string) => {
  try {
    await pb.collection("webpushDevices").delete(deviceId);
    await loadDevices();
  } catch (error: any) {
    errorMessage.value = "Failed to remove device: " + error.message;
  }
};

// Modify unregisterDevice to use currentDevice
const unregisterDevice = async () => {
  try {
    await disablePushNotifications();
    isSubscribed.value = false;
    deviceName.value = "";
    await loadDevices();
  } catch (error: any) {
    errorMessage.value = error.message || "Failed to unregister device";
  }
};

const disablePushNotifications = async () => {
  if (!("serviceWorker" in navigator)) {
    throw new Error("Service worker not supported");
  }

  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();

  if (subscription && currentDevice.value) {
    await subscription.unsubscribe();
    await pb.collection("webpushDevices").delete(currentDevice.value.id);
  }
};

const registerDevice = async () => {
  if (isRegistering.value) return;

  try {
    isRegistering.value = true;
    await registerPushNotifications();
    isSubscribed.value = true;
  } catch (error: any) {
    errorMessage.value = error.message || "Failed to register device";
  } finally {
    isRegistering.value = false;
  }
};

const registerPushNotifications = async () => {
  if (!("serviceWorker" in navigator)) {
    throw new Error("Service worker not supported");
  }

  if (vapidPublicKey.value === null) {
    throw new Error("Vapid public key not available");
  }

  if (!deviceNameValid.value) {
    throw new Error("Please enter a valid device name (min. 3 characters)");
  }

  const registration = await navigator.serviceWorker.ready;
  const options = {
    userVisibleOnly: true,
    applicationServerKey: vapidPublicKey.value,
  };

  try {
    const subscription = await registration.pushManager.subscribe(options);
    await registerDeviceInDatabase(subscription);
    console.log("Push notification subscription successful");
  } catch (error) {
    console.error("Push subscription failed:", error);
    throw error;
  }
};

const isNameUnique = computed(() => {
  return !devices.value.some(
    (device) =>
      device.Name.toLowerCase() === deviceName.value.toLowerCase() &&
      device.id !== currentDevice.value?.id
  );
});

const deviceNameErrorMessage = computed(() => {
  if (deviceName.value.length < 3) {
    return "Device name must be at least 3 characters long";
  }
  if (!isNameUnique.value) {
    return "This device name is already in use";
  }
  return "";
});

const deviceNameValid = computed(
  () => deviceName.value.length >= 3 && isNameUnique.value
);

const openChannelModal = (device: WebpushDevicesResponse) => {
  currentManageDevice.value = device;
  showChannelModal.value = true;
};

const updateDeviceChannels = async (channels: string[]) => {
  if (!currentManageDevice.value) return;

  try {
    await pb.collection("webpushDevices").update(currentManageDevice.value.id, {
      channel: channels,
    });

    toast.show?.({
      props: {
        title: "Success",
        body: "Channels updated successfully",
        variant: "success",
      },
    });

    await loadDevices();
  } catch (error) {
    errorMessage.value = "Failed to update channels";
  }
};
</script>

<style scoped>
.list-group-item {
  word-break: break-all;
}
</style>
