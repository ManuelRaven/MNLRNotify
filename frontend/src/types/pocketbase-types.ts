/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	ApplicationKV = "ApplicationKV",
	KVStore = "KVStore",
	Authorigins = "_authOrigins",
	Externalauths = "_externalAuths",
	Mfas = "_mfas",
	Otps = "_otps",
	Superusers = "_superusers",
	Channel = "channel",
	ChannelMutations = "channel_mutations",
	Message = "message",
	Reciever = "reciever",
	Sender = "sender",
	Users = "users",
	WebpushDevices = "webpushDevices",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

type ExpandType<T> = unknown extends T
	? T extends unknown
		? { expand?: unknown }
		: { expand: T }
	: { expand: T }

// System fields
export type BaseSystemFields<T = unknown> = {
	id: RecordIdString
	collectionId: string
	collectionName: Collections
} & ExpandType<T>

export type AuthSystemFields<T = unknown> = {
	email: string
	emailVisibility: boolean
	username: string
	verified: boolean
} & BaseSystemFields<T>

// Record types for each collection

export type ApplicationKVRecord<Tvalue = unknown> = {
	created?: IsoDateString
	id: string
	key: string
	updated?: IsoDateString
	value: null | Tvalue
}

export type KVStoreRecord<Tvalue = unknown> = {
	created?: IsoDateString
	id: string
	owner?: RecordIdString
	updated?: IsoDateString
	value?: null | Tvalue
}

export type AuthoriginsRecord = {
	collectionRef: string
	created?: IsoDateString
	fingerprint: string
	id: string
	recordRef: string
	updated?: IsoDateString
}

export type ExternalauthsRecord = {
	collectionRef: string
	created?: IsoDateString
	id: string
	provider: string
	providerId: string
	recordRef: string
	updated?: IsoDateString
}

export type MfasRecord = {
	collectionRef: string
	created?: IsoDateString
	id: string
	method: string
	recordRef: string
	updated?: IsoDateString
}

export type OtpsRecord = {
	collectionRef: string
	created?: IsoDateString
	id: string
	password: string
	recordRef: string
	sentTo?: string
	updated?: IsoDateString
}

export type SuperusersRecord = {
	created?: IsoDateString
	email: string
	emailVisibility?: boolean
	id: string
	password: string
	tokenKey: string
	updated?: IsoDateString
	verified?: boolean
}

export type ChannelRecord = {
	created?: IsoDateString
	id: string
	lifetime_seconds?: number
	mutations?: RecordIdString[]
	name: string
	owner: RecordIdString
	updated?: IsoDateString
}

export enum ChannelMutationsExecutorOptions {
	"goja" = "goja",
	"cmd" = "cmd",
}
export type ChannelMutationsRecord = {
	created?: IsoDateString
	executor: ChannelMutationsExecutorOptions
	id: string
	name: string
	owner?: RecordIdString
	script: string
	updated?: IsoDateString
}

export enum MessageDeliveryStateOptions {
	"pending" = "pending",
	"failure" = "failure",
	"success" = "success",
}
export type MessageRecord = {
	channel: RecordIdString
	created?: IsoDateString
	deliveryMessage?: string
	deliveryState?: MessageDeliveryStateOptions
	id: string
	text: string
	updated?: IsoDateString
}

export enum RecieverTypeOptions {
	"gotify" = "gotify",
	"ntfy" = "ntfy",
}
export type RecieverRecord = {
	channels?: RecordIdString[]
	created?: IsoDateString
	id: string
	name: string
	owner: RecordIdString
	token?: string
	type: RecieverTypeOptions
	updated?: IsoDateString
}

export type SenderRecord = {
	channel?: RecordIdString[]
	created?: IsoDateString
	id: string
	name: string
	owner: RecordIdString
	sendurl: string
	splitLimit?: number
	updated?: IsoDateString
}

export type UsersRecord = {
	avatar?: string
	created?: IsoDateString
	email: string
	emailVisibility?: boolean
	id: string
	name?: string
	password: string
	tokenKey: string
	updated?: IsoDateString
	verified?: boolean
}

export type WebpushDevicesRecord<Tsubscription = unknown> = {
	Name: string
	channel?: RecordIdString[]
	created?: IsoDateString
	id: string
	owner: RecordIdString
	subscription: null | Tsubscription
	updated?: IsoDateString
}

// Response types include system fields and match responses from the PocketBase API
export type ApplicationKVResponse<Tvalue = unknown, Texpand = unknown> = Required<ApplicationKVRecord<Tvalue>> & BaseSystemFields<Texpand>
export type KVStoreResponse<Tvalue = unknown, Texpand = unknown> = Required<KVStoreRecord<Tvalue>> & BaseSystemFields<Texpand>
export type AuthoriginsResponse<Texpand = unknown> = Required<AuthoriginsRecord> & BaseSystemFields<Texpand>
export type ExternalauthsResponse<Texpand = unknown> = Required<ExternalauthsRecord> & BaseSystemFields<Texpand>
export type MfasResponse<Texpand = unknown> = Required<MfasRecord> & BaseSystemFields<Texpand>
export type OtpsResponse<Texpand = unknown> = Required<OtpsRecord> & BaseSystemFields<Texpand>
export type SuperusersResponse<Texpand = unknown> = Required<SuperusersRecord> & AuthSystemFields<Texpand>
export type ChannelResponse<Texpand = unknown> = Required<ChannelRecord> & BaseSystemFields<Texpand>
export type ChannelMutationsResponse<Texpand = unknown> = Required<ChannelMutationsRecord> & BaseSystemFields<Texpand>
export type MessageResponse<Texpand = unknown> = Required<MessageRecord> & BaseSystemFields<Texpand>
export type RecieverResponse<Texpand = unknown> = Required<RecieverRecord> & BaseSystemFields<Texpand>
export type SenderResponse<Texpand = unknown> = Required<SenderRecord> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>
export type WebpushDevicesResponse<Tsubscription = unknown, Texpand = unknown> = Required<WebpushDevicesRecord<Tsubscription>> & BaseSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	ApplicationKV: ApplicationKVRecord
	KVStore: KVStoreRecord
	_authOrigins: AuthoriginsRecord
	_externalAuths: ExternalauthsRecord
	_mfas: MfasRecord
	_otps: OtpsRecord
	_superusers: SuperusersRecord
	channel: ChannelRecord
	channel_mutations: ChannelMutationsRecord
	message: MessageRecord
	reciever: RecieverRecord
	sender: SenderRecord
	users: UsersRecord
	webpushDevices: WebpushDevicesRecord
}

export type CollectionResponses = {
	ApplicationKV: ApplicationKVResponse
	KVStore: KVStoreResponse
	_authOrigins: AuthoriginsResponse
	_externalAuths: ExternalauthsResponse
	_mfas: MfasResponse
	_otps: OtpsResponse
	_superusers: SuperusersResponse
	channel: ChannelResponse
	channel_mutations: ChannelMutationsResponse
	message: MessageResponse
	reciever: RecieverResponse
	sender: SenderResponse
	users: UsersResponse
	webpushDevices: WebpushDevicesResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: 'ApplicationKV'): RecordService<ApplicationKVResponse>
	collection(idOrName: 'KVStore'): RecordService<KVStoreResponse>
	collection(idOrName: '_authOrigins'): RecordService<AuthoriginsResponse>
	collection(idOrName: '_externalAuths'): RecordService<ExternalauthsResponse>
	collection(idOrName: '_mfas'): RecordService<MfasResponse>
	collection(idOrName: '_otps'): RecordService<OtpsResponse>
	collection(idOrName: '_superusers'): RecordService<SuperusersResponse>
	collection(idOrName: 'channel'): RecordService<ChannelResponse>
	collection(idOrName: 'channel_mutations'): RecordService<ChannelMutationsResponse>
	collection(idOrName: 'message'): RecordService<MessageResponse>
	collection(idOrName: 'reciever'): RecordService<RecieverResponse>
	collection(idOrName: 'sender'): RecordService<SenderResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
	collection(idOrName: 'webpushDevices'): RecordService<WebpushDevicesResponse>
}
