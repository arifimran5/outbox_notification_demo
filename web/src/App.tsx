import { Routes, Route, Navigate, BrowserRouter } from 'react-router-dom';
import Layout from './components/Layout';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import TopicsPage from './pages/TopicsPage';
import FeedPage from './pages/FeedPage';
import { useAuthStore } from './store/authStore';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { JSX } from 'react';

const queryClient = new QueryClient();

const ProtectedRoute = ({ children }: { children: JSX.Element }) => {
	const user = useAuthStore(s => s.user);
	if (!user) return <Navigate to="/login" replace />;
	return children;
};

export default function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<BrowserRouter>
				<Routes>
					<Route path="/login" element={<LoginPage />} />
					<Route path="/register" element={<RegisterPage />} />

					<Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
						<Route index element={<Navigate to="/topics" replace />} />
						<Route path="topics" element={<TopicsPage />} />
						<Route path="feed" element={<FeedPage />} />
					</Route>
				</Routes>
			</BrowserRouter>
		</QueryClientProvider>
	);
}