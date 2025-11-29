export interface Topic {
	id: string
	name: string
	description: string
}

export interface Post {
	id: string
	title: string
	content: string
	created_at: string
}

export type CreatePost = {
	title: string
	content: string
	topic_id: string
}
