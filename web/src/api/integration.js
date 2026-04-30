import request from "@/utils/request"

export default {
    integrationCatalog() {
        return request({
            url: '/integration/catalog',
            method: 'get'
        })
    },
    domainHealth() {
        return request({
            url: '/domain/health',
            method: 'get'
        })
    },
    domainSync() {
        return request({
            url: '/domain/sync',
            method: 'post'
        })
    },
    domainLastSync() {
        return request({
            url: '/domain/sync/last',
            method: 'get'
        })
    },
    domainZones(params) {
        return request({
            url: '/domain/zones',
            method: 'get',
            params
        })
    },
    domainRecords(params) {
        return request({
            url: '/domain/records',
            method: 'get',
            params
        })
    }
}
