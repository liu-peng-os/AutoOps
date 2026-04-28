const routes = [
    {
        path: '/integration/domain',
        component: () => import('@/views/integration/ExternalModuleLanding.vue'),
        props: {
            pageKey: 'domain-management',
            title: '域名管理',
            subtitle: '统一接入外部 DNS 管理系统，当前阶段先提供稳定入口和后续增强说明。'
        },
        meta: { sTitle: '域名管理', tTitle: '统一入口' }
    },
    {
        path: '/integration/workorder',
        component: () => import('@/views/integration/ExternalModuleLanding.vue'),
        props: {
            pageKey: 'operations-workorder',
            title: '运营工单',
            subtitle: '以外部表格为主数据源，后续在这里承接同步、展示、编辑和回写能力。'
        },
        meta: { sTitle: '运营工单', tTitle: '工单中心' }
    }
]

export default routes
