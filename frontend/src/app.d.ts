// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
import type { PrismaClient } from '@prisma/client';
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			userId: string;
			name: string;
			email: string;
		}
		// interface PageData {}
		// interface Platform {}
	}
}

export {};
