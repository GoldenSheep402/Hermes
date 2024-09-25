import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const CATEGORY: AppRouteRecordRaw = {
	path: '/category',
	name: 'Category',
	component: DEFAULT_LAYOUT,
	meta: {
		label: '类别',
		requiresAuth: true,
		icon: 'icon-apps',
		order: 1,
	},
	children: [
		{
			path: 'list',
			name: 'CategoryList',
			component: () => import('@/views/category/index.vue'),
			meta: {
				label: '列表',
				requiresAuth: true,
				roles: ['*'],
			},
		},
	],
};

export default CATEGORY;
