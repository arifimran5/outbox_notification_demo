import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { useNavigate, Link } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';
import { Loader2 } from 'lucide-react';
import apiClient from '../api/apiclient';
import { toast } from "sonner"


const RegisterPage = () => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const navigate = useNavigate();

	const mutation = useMutation({
		mutationFn: async (data: { username: string; password: string }) => (await apiClient.post('/register', data)).data,
		onSuccess: () => {
			toast.success('Registered')

			setTimeout(() => {
				toast.info('Navigating to login page..')
				setTimeout(() => {
					navigate('/login');
				}, 1000)

			}, 1000)

		},
		onError: () => toast.error('Failed to register')
	});

	return (
		<div className="flex h-screen items-center justify-center bg-gray-50">
			<Card className="w-[350px]">
				<CardHeader>
					<CardTitle>Create Account</CardTitle>
				</CardHeader>
				<CardContent>
					<form onSubmit={(e) => { e.preventDefault(); mutation.mutate({ username, password }); }} className="space-y-4">
						<Input placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
						<Input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
						<Button className="w-full" disabled={mutation.isPending}>
							{mutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />} Register
						</Button>
					</form>
					<div className="mt-4 text-center text-sm text-gray-500">
						<Link to="/login" className="underline">Back to Login</Link>
					</div>
				</CardContent>
			</Card>
		</div>
	);
};

export default RegisterPage;