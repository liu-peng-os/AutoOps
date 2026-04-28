const routes = [
    {
        path: '/ops/tools',
        component: () => import('@/views/Tools/Tools'),
        meta: {sTitle: '运维工具', tTitle: '工具列表'}
    },
    {
        path: '/ops/agent',
        component: () => import('@/views/Tools/Agent'),
        meta: {sTitle: '运维工具', tTitle: 'agent列表'}
    },
    {
        path: '/ops/knowledge',
        component: () => import('@/views/Tools/Knowledge'),
        meta: {sTitle: '运维工具', tTitle: '运维知识库'}
    },
]

export default routes
