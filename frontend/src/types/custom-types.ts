import type { UsersRecord } from "./pocketbase-types";

// type based on UserRecord but make the id and tokenKey optional
export type UserRequest = Omit<UsersRecord, "id" | "tokenKey"> & {
  passwordConfirm: string;
};

export type AuthStoreUserRecord = Omit<UsersRecord, "password" | "tokenKey">;

export type CreateChannelRequest = {
  name: string;
};

export type CreateSenderRequest = {
  name: string;
  sendurl: string;
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
