import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { Button } from '../components/ui/Button';
import { Loader2 } from 'lucide-react';
import apiClient from '../api/apiclient';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import { Input } from '../components/ui/Input';
import type { AxiosError } from 'axios';

const LoginPage = () => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const setAuth = useAuthStore(s => s.setAuth);
	const navigate = useNavigate();

	const mutation = useMutation({
		mutationFn: async (data: { username: string; password: string }) => (await apiClient.post('/login', data)).data,
		onSuccess: (data) => {
			setAuth({ id: data.user_id, username }, data.token);
			navigate('/topics');
		},
		onError: (err: AxiosError) => alert(err.response?.data || 'Login failed')
	});

	return (
		<div className="flex h-screen items-center justify-center bg-gray-50">
			<Card className="w-[350px]">
				<CardHeader>
					<CardTitle>Welcome Back</CardTitle>
				</CardHeader>
				<CardContent>
					<form onSubmit={(e) => { e.preventDefault(); mutation.mutate({ username, password }); }} className="space-y-4">
						<Input placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
						<Input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
						<Button className="w-full" disabled={mutation.isPending}>
							{mutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Login
						</Button>
					</form>
					<div className="mt-4 text-center text-sm text-gray-500">
						<Link to="/register" className="underline">Create an account</Link>
					</div>
				</CardContent>
			</Card>
		</div>
	);
};
export default LoginPage;