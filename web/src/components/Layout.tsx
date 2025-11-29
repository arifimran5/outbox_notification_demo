import { useEffect, useState } from 'react';
import { Outlet, Link, useNavigate } from 'react-router-dom';
import { Bell, LogOut, Activity } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useNotificationStore } from '../store/notificationStore';
import { Button } from './ui/Button';
import { Toaster } from './ui/Sonner';

import { fetchEventSource } from '@microsoft/fetch-event-source';

const Layout = () => {
	const { user, token, logout } = useAuthStore();
	const { notifications, unreadCount, addNotification, markAllRead } = useNotificationStore();
	const [showDropdown, setShowDropdown] = useState(false);
	const navigate = useNavigate();

	useEffect(() => {
		if (!token) return;

		const controller = new AbortController();

		const API_URL = `${import.meta.env.VITE_API_URL}/api`

		fetchEventSource(`${API_URL}/events`, {
			method: 'GET',
			headers: {
				'Authorization': `Bearer ${token}`,
				'Content-Type': 'text/event-stream',
			},
			signal: controller.signal,
			openWhenHidden: true,  // Keep connection alive even when tab hidden

			// onopen(response) {
			// 	console.log('SSE Connected');
			// 	if (response.headers.get('content-type') !== 'text/event-stream') {
			// 		throw new Error('Expected SSE stream');
			// 	}
			// 	return new Promise((res) => res());
			// },

			onmessage(ev) {
				if (ev.data === 'connected') return;

				try {
					const payload = JSON.parse(ev.data);
					addNotification({
						...payload,
						id: crypto.randomUUID(),
						timestamp: new Date().toLocaleTimeString(),
						read: false
					});
				} catch (e) {
					console.error('Failed to parse SSE', e);
				}
			},

			onclose() {
				console.log('SSE Connection closed');
			},

			onerror(err) {
				console.error('SSE Error:', err);
				// Auto-reconnects automatically
			}
		});

		return () => {
			controller.abort();
		};
	}, [token, addNotification]);

	const handleLogout = () => {
		logout();
		navigate('/login');
	};

	const handleNotificationClick = () => {
		setShowDropdown(!showDropdown);
		navigate('/feed')
	}

	return (
		<div className="min-h-screen bg-gray-50">
			<Toaster />
			<header className="sticky top-0 z-50 w-full border-b bg-white">
				<div className="container mx-auto flex h-14 items-center px-4">
					<div className="flex items-center gap-2 font-bold text-xl mr-8">
						<div className="bg-blue-600 text-white p-1 rounded">
							<Activity className="h-5 w-5" />
						</div>
						OutboxDemo
					</div>

					<nav className="flex items-center space-x-6 text-sm font-medium text-gray-600">
						<Link to="/topics" className="hover:text-blue-600">Topics</Link>
						<Link to="/feed" className="hover:text-blue-600">My Feed</Link>
					</nav>

					<div className="ml-auto flex items-center gap-4">
						<div className="relative">
							<Button variant="ghost" className="relative" onClick={() => { setShowDropdown(!showDropdown); if (!showDropdown) markAllRead(); }}>
								<Bell className="h-5 w-5" />
								{unreadCount > 0 && (
									<span className="absolute top-0 right-0 h-4 w-4 bg-red-500 rounded-full text-[10px] text-white flex items-center justify-center">
										{unreadCount}
									</span>
								)}
							</Button>
							{showDropdown && (
								<div className="absolute right-0 mt-2 w-80 rounded-md border bg-white shadow-lg z-50">
									<div className="p-3 border-b text-sm font-semibold">Notifications</div>
									<div className="max-h-64 overflow-y-auto">
										{notifications.length === 0 ? (
											<div className="p-4 text-center text-sm text-gray-500">No new notifications</div>
										) : (
											notifications.map(n => (
												<div key={n.id} className="p-3 border-b hover:bg-gray-50 text-sm cursor-pointer" onClick={handleNotificationClick}>
													<p className="font-medium">{n.message}</p>
													<p className="text-xs text-gray-500 mt-1">{n.timestamp}</p>
												</div>
											))
										)}
									</div>
								</div>
							)}
						</div>

						<div className="flex items-center gap-2 border-l pl-4">
							<span className="text-sm font-medium">{user?.username}</span>
							<Button variant="ghost" onClick={handleLogout}>
								<LogOut className="h-4 w-4" />
							</Button>
						</div>
					</div>
				</div>
			</header>
			<main className="container mx-auto py-8 px-4">
				<Outlet />
			</main>
		</div>
	);
};


export default Layout;