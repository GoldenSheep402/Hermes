import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const USERINFO: AppRouteRecordRaw = {
	path: '/user',
	name: 'user',
	component: DEFAULT_LAYOUT,
	meta: {
		label: '用户',
		requiresAuth: true,
		icon: 'icon-command',
		order: 10,
		hideInMenu: true,
	},
	children: [
		{
			path: 'info',
			name: 'UserInfo',
			component: () => import('@/views/userInfo/index.vue'),
			meta: {
				label: '用户信息',
				requiresAuth: true,
				roles: ['*'],
				hideInMenu: true,
			},
		},
	],
};

export default USERINFO;
