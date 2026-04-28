const routes = [
    {
        path: '/system/personal',
        component: () => import('@/views/system/Personal.vue'),
        meta: {sTitle: '个人中心', tTitle: '个人信息'}
    },
    {
        path: '/system/admin',
        component: () => import('@/views/system/Admin.vue'),
        meta: {sTitle: '基础管理', tTitle: '用户信息'}
    },
    {
        path: '/system/role',
        component: () => import('@/views/system/Role.vue'),
        meta: {sTitle: '基础管理', tTitle: '角色信息'}
    },
    {
        path: '/system/menu',
        component: () => import('@/views/system/Menu.vue'),
        meta: {sTitle: '基础管理', tTitle: '菜单信息'}
    },
    {
        path: '/system/dept',
        component: () => import('@/views/system/Dept.vue'),
        meta: {sTitle: '基础管理', tTitle: '部门信息'}
    },
    {
        path: '/system/post',
        component: () => import('@/views/system/Post.vue'),
        meta: {sTitle: '基础管理', tTitle: '岗位信息'}
    },
    {
        path: '/monitor/loginlog',
        component: () => import('@/views/monitor/LoginLog.vue'),
        meta: {sTitle: '日志管理', tTitle: '登录日志'}
    },
    {
        path: '/monitor/operator',
        component: () => import('@/views/monitor/Operator.vue'),
        meta: {sTitle: '日志管理', tTitle: '操作日志'}
    },
    {
        path: '/monitor/dblog',
        component: () => import('@/views/monitor/DBLog.vue'),
        meta: {sTitle: '日志管理', tTitle: '数据日志'}
    }
]

export default routes
