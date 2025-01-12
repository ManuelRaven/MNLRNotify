import type {
  ChannelMutationsRecord,
  UsersRecord,
  WebpushDevicesRecord,
} from "./pocketbase-types";

// type based on UserRecord but make the id and tokenKey optional
export type UserRequest = Omit<UsersRecord, "id" | "tokenKey"> & {
  passwordConfirm: string;
};

export type WebpushDeviceRequest = Omit<WebpushDevicesRecord, "id">;

export type AuthStoreUserRecord = Omit<UsersRecord, "password" | "tokenKey">;

export type ChannelMutationsRequest = Omit<ChannelMutationsRecord, "id">;

export type CreateSenderRequest = {
  name: string;
  sendurl: string;
};

export type CreateChannelRequest = {
  name: string;
};

export type ExpandChannelNameId = {
  channel: {
    id: string;
    name: string;
  }[];
};

export type ExpandChannelsNameId = {
  channels: {
    id: string;
    name: string;
  }[];
};

export type ExpandChannelNameIdSingle = {
  channel: {
    id: string;
    name: string;
  };
};

export type ExpandMutations = {
  mutations: {
    id: string;
    name: string;
  }[];
};
