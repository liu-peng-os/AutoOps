import request from "@/utils/request"

export default {
    integrationCatalog() {
        return request({
            url: '/integration/catalog',
            method: 'get'
        })
    }
}
