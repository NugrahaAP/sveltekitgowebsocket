import { writable } from 'svelte/store';
import { responseToastEnum } from '$lib/types/myTypes';

export type toastType = {
	message: string[] | undefined;
	type: responseToastEnum | undefined;
};

let obj: toastType = {
	message: undefined,
	type: responseToastEnum.warning
};

export const toastData = writable(obj);
export const currentRoomId = writable({ roomId: '' });
export const groupName = writable('');
