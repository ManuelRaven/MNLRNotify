/**
* This file was @generated using pocketbase-typegen
*/

import type PocketBase from 'pocketbase'
import type { RecordService } from 'pocketbase'

export enum Collections {
	Authorigins = "_authOrigins",
	Externalauths = "_externalAuths",
	Mfas = "_mfas",
	Otps = "_otps",
	Superusers = "_superusers",
	Channel = "channel",
	Message = "message",
	Reciever = "reciever",
	Sender = "sender",
	Users = "users",
}

// Alias types for improved usability
export type IsoDateString = string
export type RecordIdString = string
export type HTMLString = string

// System fields
export type BaseSystemFields<T = never> = {
	id: RecordIdString
	collectionId: string
	collectionName: Collections
	expand?: T
}

export type AuthSystemFields<T = never> = {
	email: string
	emailVisibility: boolean
	username: string
	verified: boolean
} & BaseSystemFields<T>

// Record types for each collection

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
	name: string
	owner: RecordIdString
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

// Response types include system fields and match responses from the PocketBase API
export type AuthoriginsResponse<Texpand = unknown> = Required<AuthoriginsRecord> & BaseSystemFields<Texpand>
export type ExternalauthsResponse<Texpand = unknown> = Required<ExternalauthsRecord> & BaseSystemFields<Texpand>
export type MfasResponse<Texpand = unknown> = Required<MfasRecord> & BaseSystemFields<Texpand>
export type OtpsResponse<Texpand = unknown> = Required<OtpsRecord> & BaseSystemFields<Texpand>
export type SuperusersResponse<Texpand = unknown> = Required<SuperusersRecord> & AuthSystemFields<Texpand>
export type ChannelResponse<Texpand = unknown> = Required<ChannelRecord> & BaseSystemFields<Texpand>
export type MessageResponse<Texpand = unknown> = Required<MessageRecord> & BaseSystemFields<Texpand>
export type RecieverResponse<Texpand = unknown> = Required<RecieverRecord> & BaseSystemFields<Texpand>
export type SenderResponse<Texpand = unknown> = Required<SenderRecord> & BaseSystemFields<Texpand>
export type UsersResponse<Texpand = unknown> = Required<UsersRecord> & AuthSystemFields<Texpand>

// Types containing all Records and Responses, useful for creating typing helper functions

export type CollectionRecords = {
	_authOrigins: AuthoriginsRecord
	_externalAuths: ExternalauthsRecord
	_mfas: MfasRecord
	_otps: OtpsRecord
	_superusers: SuperusersRecord
	channel: ChannelRecord
	message: MessageRecord
	reciever: RecieverRecord
	sender: SenderRecord
	users: UsersRecord
}

export type CollectionResponses = {
	_authOrigins: AuthoriginsResponse
	_externalAuths: ExternalauthsResponse
	_mfas: MfasResponse
	_otps: OtpsResponse
	_superusers: SuperusersResponse
	channel: ChannelResponse
	message: MessageResponse
	reciever: RecieverResponse
	sender: SenderResponse
	users: UsersResponse
}

// Type for usage with type asserted PocketBase instance
// https://github.com/pocketbase/js-sdk#specify-typescript-definitions

export type TypedPocketBase = PocketBase & {
	collection(idOrName: '_authOrigins'): RecordService<AuthoriginsResponse>
	collection(idOrName: '_externalAuths'): RecordService<ExternalauthsResponse>
	collection(idOrName: '_mfas'): RecordService<MfasResponse>
	collection(idOrName: '_otps'): RecordService<OtpsResponse>
	collection(idOrName: '_superusers'): RecordService<SuperusersResponse>
	collection(idOrName: 'channel'): RecordService<ChannelResponse>
	collection(idOrName: 'message'): RecordService<MessageResponse>
	collection(idOrName: 'reciever'): RecordService<RecieverResponse>
	collection(idOrName: 'sender'): RecordService<SenderResponse>
	collection(idOrName: 'users'): RecordService<UsersResponse>
}
