export default [
	{
		path: 'reports',
		components: {
			default: () => import('./Reports.vue')
		},
		props: {
			default: true,
			sidebar: false
		},
		meta: {
			authentication: {
				rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
			}
		}
	}
];