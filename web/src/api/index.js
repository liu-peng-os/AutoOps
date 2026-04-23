/**
 * 基础API接口管理
 *
 * @author xiaoRui
 */

import request from "@/utils/request"
import systemAPI from './system'
import cmdbAPI from './cmdb'
import dashboardAPI from './dashboard'
import integrationAPI from './integration'
import * as toolAPI from './tool'

export default {
    // 基础API
    captcha() {
        return request({
            url: 'captcha',
            method: 'get'
        })
    },
    login(params) {
        return request({
            url: 'login',
            method: 'post',
            data: params
        })
    },
    // 系统管理API - 扩展到根级别以兼容现有代码
    ...systemAPI,
    // CMDB管理API - 扩展到根级别以兼容现有代码
    ...cmdbAPI,
    // Dashboard管理API - 扩展到根级别以兼容现有代码
    ...dashboardAPI,
    // Integration管理API
    ...integrationAPI,
    // 导航工具API - 扩展到根级别以兼容现有代码
    ...toolAPI,

    // 同时提供命名空间访问
    system: systemAPI,
    cmdb: cmdbAPI,
    dashboard: dashboardAPI,
    integration: integrationAPI,
    tool: toolAPI
}
