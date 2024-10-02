import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const SYSTEM: AppRouteRecordRaw = {
	path: '/system',
	name: 'System',
	component: DEFAULT_LAYOUT,
	meta: {
		label: '系统',
		requiresAuth: true,
		icon: 'icon-settings',
		order: 4,
		roles: ['admin'],
	},
	children: [
		{
			path: 'setting',
			name: 'SystemSetting',
			component: () => import('@/views/system/index.vue'),
			meta: {
				label: '系统设置',
				requiresAuth: true,
				roles: ['admin'],
			},
		},
	],
};

export default SYSTEM;
