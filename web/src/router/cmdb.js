const routes = [
    {
        path: '/cmdb/ecs',
        component: () => import('@/views/cmdb/cmdbHost.vue'),
        meta: {sTitle: '资产管理', tTitle: '主机管理'}
    },
    {
        path: '/cmdb/group',
        component: () => import('@/views/cmdb/cmdbGroup.vue'),
        meta: {sTitle: '资产管理', tTitle: '业务分组'}
    },
    {
        path: '/cmdb/db',
        component: () => import('@/views/cmdb/cmdbDB.vue'),
        meta: {sTitle: '资产管理', tTitle: '数据管理'}
    },
    {
        path: '/cmdb/ssh',
        component: () => import('@/views/cmdb/Host/SSH.vue'),
        meta: {sTitle: '资产管理', tTitle: '终端登录'}
    },
    {
        path: '/cmdb/dbdetails',
        component: () => import('@/views/cmdb/DBdetails.vue'),
        meta: {sTitle: '数据管理', tTitle: '数据库操作'}
    },
    {
        path: '/cmdb/site',
        component: () => import('@/views/integration/ExternalModuleLanding.vue'),
        props: {
            pageKey: 'site-management',
            title: '站点管理',
            subtitle: '站点数据将继续以外部表格为主数据源，平台负责统一入口和后续同步回写。'
        },
        meta: {sTitle: '资产管理', tTitle: '站点管理'}
    }
]

export default routes
