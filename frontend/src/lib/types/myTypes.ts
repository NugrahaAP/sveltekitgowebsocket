export enum responseToastEnum {
	primary = 'primary',
	secondary = 'secondary',
	tertiary = 'tertiary',
	warning = 'warning',
	error = 'error'
}

export enum messageTypeEnum {
	message = 'message',
	info = 'info'
}

export enum chatRoomTypeEnum {
	personal = 'personal',
	group = 'group'
}

export type responseToast = {
	error: boolean;
	message: string[];
	type: responseToastEnum;
};

export type ServerResponse = {
	isError: boolean;
	messages: string[];
	data: any;
};

export type User = {
	id: string;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
	name: string;
	email: string;
};

export type Message = {
	id: string;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
	chatRoomId: string;
	messageBody: string;
	messageType: messageTypeEnum;
	messageLink: string;
	sender: User;
	userId: string;
};

export type ChatRoom = {
	id: string;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
	participant: User[];
	messages: Message[];
	chatRoomType: chatRoomTypeEnum;
};

export type GroupChatRoom = {
	id: string;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
	chatRoom: ChatRoom;
	roomName: string;
	roleAdmin: User[];
	chatRoomId: string;
};

export enum actionEnum {
	show = 'show',
	delete = 'delete'
}

export type WsResponse = {
	action: string;
	message: string;
	senderId: string;
	createdAt: string;
	name: string;
};

export type WsUser = Omit<User, 'createdAt' | 'updatedAt' | 'deletedAt' | 'email'>;
export type WsChat = Omit<Message, 'updatedAt' | 'deletedAt' | 'sender'>;
