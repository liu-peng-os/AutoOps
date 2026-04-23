import Host from '@/views/cmdb/cmdbHost.vue'
import Group from '@/views/cmdb/cmdbGroup.vue'
import Db from '@/views/cmdb/cmdbDB.vue'
import SSH from '@/views/cmdb/Host/SSH.vue'
import DBdetails from '@/views/cmdb/DBdetails.vue'
import ExternalModuleLanding from '@/views/integration/ExternalModuleLanding.vue'


const routes = [
    {
        path: '/cmdb/ecs',
        component: Host,
        meta: {sTitle: '资产管理', tTitle: '主机管理'}
    },
    {
        path: '/cmdb/group',
        component: Group,
        meta: {sTitle: '资产管理', tTitle: '业务分组'}
    },
    {
        path: '/cmdb/db',
        component: Db,
        meta: {sTitle: '资产管理', tTitle: '数据管理'}
    }, 
    {
        path: '/cmdb/ssh',
        component: SSH,
        meta: {sTitle: '资产管理', tTitle: '终端登录'}
    },
    {
        path: '/cmdb/dbdetails',
        component: DBdetails,
        meta: {sTitle: '数据管理', tTitle: '数据库操作'}
    },
    {
        path: '/cmdb/site',
        component: ExternalModuleLanding,
        props: {
            pageKey: 'site-management',
            title: '站点管理',
            subtitle: '站点数据将继续以外部表格为主数据源，平台负责统一入口和后续同步回写。'
        },
        meta: {sTitle: '资产管理', tTitle: '站点管理'}
    }
]

export default routes
