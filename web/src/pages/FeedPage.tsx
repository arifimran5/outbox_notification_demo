import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../api/apiclient';
import { Input } from '../components/ui/Input';
import { Textarea } from '../components/ui/TextArea';
import { Plus } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import type { CreatePost, Post, Topic } from '../types/type';
import { toast } from "sonner"



const FeedPage = () => {
  const [selectedTopic, setSelectedTopic] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const queryClient = useQueryClient();

  // 1. Fetch Subscriptions
  const { data: subs } = useQuery<Topic[]>({
    queryKey: ['subs'],
    queryFn: async () => (await apiClient.get('/subscriptions')).data
  });

  const activeTopicId = selectedTopic || subs?.[0]?.id;

  // 2. Fetch Posts
  const { data: posts, isLoading } = useQuery<Post[]>({
    queryKey: ['posts', activeTopicId],
    queryFn: async () => {
      if (!activeTopicId) return [];
      return (await apiClient.get(`/topics/${activeTopicId}/posts`)).data;
    },
    enabled: !!activeTopicId
  });

  // 3. Create Post
  const postMutation = useMutation({
    mutationFn: async (data: CreatePost) => apiClient.post(`/topics/${activeTopicId}/posts`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts', activeTopicId] });
      setIsModalOpen(false);
      setTitle('');
      setContent('');
      // alert('Post created! Notification queued.');
		toast.success('Post created!')
		
    }
  });

  const handlePost = (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeTopicId) return;
    postMutation.mutate({ title, content, topic_id: activeTopicId });
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
      {/* Sidebar */}
      <div className="md:col-span-1 space-y-4">
        <h3 className="font-semibold text-gray-900">Subscriptions</h3>
        <div className="flex flex-col space-y-1">
          {subs?.length === 0 && <span className="text-sm text-gray-500">No subscriptions yet.</span>}
          {subs?.map(s => (
            <Button 
              key={s.id} 
              variant={activeTopicId === s.id ? 'secondary' : 'ghost'} 
              className="justify-start font-normal"
              onClick={() => setSelectedTopic(s.id)}
            >
              # {s.name}
            </Button>
          ))}
        </div>
      </div>

      {/* Main Feed */}
      <div className="md:col-span-3 space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">
            {subs?.find(s => s.id === activeTopicId)?.name || 'Feed'}
          </h2>
          {activeTopicId && (
            <div className="relative">
              <Button onClick={() => setIsModalOpen(!isModalOpen)}>
                <Plus className="mr-2 h-4 w-4" /> Create Post
              </Button>
              {isModalOpen && (
                <div className="absolute right-0 top-12 w-[400px] z-50 bg-white border shadow-xl rounded-xl p-6 animate-in fade-in zoom-in-95">
                  <h3 className="font-bold mb-4">New Post</h3>
                  <form onSubmit={handlePost} className="space-y-4">
                    <Input placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} required />
                    <Textarea placeholder="What's happening?" value={content} onChange={e => setContent(e.target.value)} required />
                    <div className="flex justify-end gap-2">
                      <Button type="button" variant="outline" onClick={() => setIsModalOpen(false)}>Cancel</Button>
                      <Button type="submit" disabled={postMutation.isPending}>Publish</Button>
                    </div>
                  </form>
                </div>
              )}
            </div>
          )}
        </div>

        {isLoading ? (
          <div>Loading posts...</div>
        ) : posts?.length === 0 ? (
          <div className="p-10 text-center text-gray-500 bg-gray-50 rounded-xl border border-dashed">
            No posts yet in this topic.
          </div>
        ) : (
          posts?.map(p => (
            <Card key={p.id} className="overflow-hidden">
              <CardHeader className="bg-gray-50/50 pb-4">
                <CardTitle className="text-lg">{p.title}</CardTitle>
                <span className="text-xs text-gray-500">{new Date(p.created_at).toLocaleString()}</span>
              </CardHeader>
              <CardContent className="pt-4">
                <p className="whitespace-pre-wrap text-gray-700">{p.content}</p>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  );
};


export default FeedPage;