import { create } from 'zustand'

interface Notification {
	id: string // generated client side or from event ID
	message: string
	topic_name: string
	post_id: string
	timestamp: string
	read: boolean
}

interface NotificationState {
	notifications: Notification[]
	unreadCount: number
	addNotification: (n: Notification) => void
	markAllRead: () => void
	clear: () => void
}

export const useNotificationStore = create<NotificationState>((set) => ({
	notifications: [],
	unreadCount: 0,
	addNotification: (n) =>
		set((state) => ({
			notifications: [n, ...state.notifications],
			unreadCount: state.unreadCount + 1,
		})),
	markAllRead: () =>
		set((state) => ({
			notifications: state.notifications.map((n) => ({ ...n, read: true })),
			unreadCount: 0,
		})),
	clear: () => set({ notifications: [], unreadCount: 0 }),
}))
