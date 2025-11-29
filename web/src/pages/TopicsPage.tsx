import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../api/apiclient';
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { Badge } from '../components/ui/Badge';
import { Check } from 'lucide-react';
import type { Topic } from '../types/type';

const TopicsPage = () => {
	const queryClient = useQueryClient();

	const { data: topics, isLoading } = useQuery<Topic[]>({
		queryKey: ['topics'],
		queryFn: async () => (await apiClient.get('/topics')).data
	});

	const { data: subs } = useQuery<Topic[]>({
		queryKey: ['subs'],
		queryFn: async () => (await apiClient.get('/subscriptions')).data
	});

	const subMutation = useMutation({
		mutationFn: async (id: string) => apiClient.post(`/topics/${id}/subscribe`),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ['subs'] })
	});

	const unsubMutation = useMutation({
		mutationFn: async (id: string) => apiClient.post(`/topics/${id}/unsubscribe`),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ['subs'] })
	});

	const isSubbed = (id: string) => subs?.some(s => s.id === id);

	if (isLoading) return <div className="p-8 text-center">Loading...</div>;

	return (
		<div className="space-y-6">
			<h2 className="text-2xl font-bold">Available Topics</h2>
			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{topics?.map(t => (
					<Card key={t.id} className="hover:shadow-md transition-shadow">
						<CardHeader>
							<div className="flex justify-between items-start">
								<CardTitle>{t.name}</CardTitle>
								{isSubbed(t.id) && <Badge variant="secondary" className="gap-1"><Check className="w-3 h-3" /> Subscribed</Badge>}
							</div>
							<p className="text-sm text-gray-500 mt-2">{t.description}</p>
						</CardHeader>
						<CardContent>

							<>
								{isSubbed(t.id) && <Button
									variant="outline"
									className="w-full"
									onClick={() => unsubMutation.mutate(t.id)}
									disabled={unsubMutation.isPending}
								>
									Unsubscribe
								</Button>}

								{!isSubbed(t.id) && <Button
									variant="default"
									className="w-full"
									onClick={() => subMutation.mutate(t.id)}
									disabled={subMutation.isPending}
								>
									Subscribe
								</Button>}
							</>
						</CardContent>
					</Card>
				))}
			</div>
		</div>
	);
};

export default TopicsPage;