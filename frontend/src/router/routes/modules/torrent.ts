import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const TORRENT: AppRouteRecordRaw = {
	path: '/torrent',
	name: 'Torrent',
	component: DEFAULT_LAYOUT,
	meta: {
		label: '种子',
		requiresAuth: true,
		icon: 'icon-relation',
		order: 4,
	},
	children: [
		{
			path: 'list',
			name: 'TorrentList',
			component: () => import('@/views/torrent/index.vue'),
			meta: {
				label: '种子列表',
				requiresAuth: true,
				roles: ['*'],
			},
		},
		{
			path: 'create',
			name: 'TorrentCreate',
			component: () => import('@/views/torrent/create/index.vue'),
			meta: {
				label: '创建种子',
				requiresAuth: true,
				roles: ['*'],
			},
		},
	],
};

export default TORRENT;
