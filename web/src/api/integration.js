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
    domainList(params) {
        return request({
            url: '/domain/domains',
            method: 'get',
            params
        })
    },
    domainDetail(id, params) {
        return request({
            url: `/domain/domains/${id}`,
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
    },
    domainRecordAdd(domainId, data) {
        return request({
            url: `/domain/domains/${domainId}/records`,
            method: 'post',
            data
        })
    },
    domainRecordUpdate(domainId, recordId, data) {
        return request({
            url: `/domain/domains/${domainId}/records/${recordId}`,
            method: 'put',
            data
        })
    },
    domainRecordDelete(domainId, recordId) {
        return request({
            url: `/domain/domains/${domainId}/records/${recordId}`,
            method: 'delete'
        })
    },
    domainRecordStatus(domainId, recordId, status) {
        return request({
            url: `/domain/domains/${domainId}/records/${recordId}/status`,
            method: 'put',
            data: { status }
        })
    },
    domainRecordRemark(domainId, recordId, remark) {
        return request({
            url: `/domain/domains/${domainId}/records/${recordId}/remark`,
            method: 'put',
            data: { remark }
        })
    },
    domainRecordBatch(domainId, data) {
        return request({
            url: `/domain/domains/${domainId}/records/batch`,
            method: 'post',
            data
        })
    }
}
