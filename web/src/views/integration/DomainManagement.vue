<template>
  <div class="domain-page">
    <section class="domain-hero">
      <div>
        <p class="eyebrow">DNSMGR LIVE API</p>
        <h1>域名管理</h1>
        <p class="subtitle">
          直接调用 dnsmgr API，不在 AutoOps 本地同步或保存域名数据。所有变更都会实时提交到 dnsmgr。
        </p>
      </div>
      <div class="actions">
        <el-button :loading="healthLoading" @click="loadHealth">健康检查</el-button>
        <el-button type="primary" :loading="domainLoading" @click="loadDomains">刷新域名</el-button>
      </div>
    </section>

    <section class="status-grid">
      <article class="status-card">
        <span>dnsmgr</span>
        <strong :class="health.healthy ? 'ok' : 'warn'">{{ healthText }}</strong>
        <small>{{ health.message || '尚未检查' }}</small>
      </article>
      <article class="status-card">
        <span>当前域名</span>
        <strong>{{ selectedDomainName || '未选择' }}</strong>
        <small>点击域名行查看解析记录</small>
      </article>
      <article class="status-card">
        <span>数据模式</span>
        <strong>实时 API</strong>
        <small>不落本地库，不做同步任务</small>
      </article>
    </section>

    <section class="table-panel">
      <div class="table-head">
        <div>
          <h2>域名列表</h2>
          <p>来自 dnsmgr 的实时域名数据。</p>
        </div>
        <el-input
          v-model="domainKeyword"
          clearable
          placeholder="搜索域名"
          class="search"
          @keyup.enter="loadDomains"
          @clear="loadDomains"
        />
      </div>
      <el-table v-loading="domainLoading" :data="domains" border highlight-current-row @row-click="selectDomain">
        <el-table-column prop="domain" label="域名" min-width="220">
          <template #default="{ row }">{{ pick(row, ['domain', 'Domain', 'name', 'Name', 'title', 'Title']) }}</template>
        </el-table-column>
        <el-table-column prop="typename" label="DNS 服务商" min-width="140">
          <template #default="{ row }">{{ pick(row, ['typename', 'TypeName', 'type', 'Type', 'provider', 'Provider']) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">{{ pick(row, ['status', 'Status', 'checkstatus', 'CheckStatus', 'state', 'State']) }}</template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="180">
          <template #default="{ row }">{{ pick(row, ['remark', 'Remark', 'note', 'Note']) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="showDomainDetail(row)">详情</el-button>
            <el-button size="small" type="primary" @click.stop="openLoginUrl(row)">登录链接</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="domainTotal"
        :current-page="domainPage"
        :page-size="domainPageSize"
        @current-change="changeDomainPage"
      />
    </section>

    <section class="table-panel">
      <div class="table-head">
        <div>
          <h2>解析记录</h2>
          <p>{{ selectedDomainName ? `当前域名：${selectedDomainName}` : '请先点击一个域名。' }}</p>
        </div>
        <div class="record-actions">
          <el-input
            v-model="recordKeyword"
            clearable
            placeholder="搜索主机记录或值"
            class="search"
            @keyup.enter="loadRecords"
            @clear="loadRecords"
          />
          <el-button type="primary" :disabled="!selectedDomainId" @click="openRecordDialog()">新增解析</el-button>
          <el-dropdown :disabled="!selectedRecordIds.length" @command="batchRecords">
            <el-button>
              批量操作
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="open">批量启用</el-dropdown-item>
                <el-dropdown-item command="pause">批量暂停</el-dropdown-item>
                <el-dropdown-item command="remark">批量改备注</el-dropdown-item>
                <el-dropdown-item command="delete">批量删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      <el-table
        v-loading="recordLoading"
        :data="records"
        border
        @selection-change="records => selectedRecords = records"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column label="主机记录" min-width="150">
          <template #default="{ row }">{{ pick(row, ['name', 'Name', 'host', 'Host', 'rr', 'RR']) }}</template>
        </el-table-column>
        <el-table-column label="类型" width="90">
          <template #default="{ row }">{{ pick(row, ['type', 'Type', 'record_type']) }}</template>
        </el-table-column>
        <el-table-column label="记录值" min-width="260" show-overflow-tooltip>
          <template #default="{ row }">{{ pick(row, ['value', 'Value', 'record', 'Record', 'content', 'Content']) }}</template>
        </el-table-column>
        <el-table-column label="线路" width="120">
          <template #default="{ row }">{{ pick(row, ['linename', 'LineName', 'line', 'Line']) }}</template>
        </el-table-column>
        <el-table-column label="TTL" width="90">
          <template #default="{ row }">{{ pick(row, ['ttl', 'TTL']) }}</template>
        </el-table-column>
        <el-table-column label="备注" min-width="160">
          <template #default="{ row }">{{ pick(row, ['remark', 'Remark', 'note', 'Note']) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">{{ pick(row, ['status', 'Status', 'state', 'State']) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="330" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openRecordDialog(row)">编辑</el-button>
            <el-button size="small" @click="changeRemark(row)">备注</el-button>
            <el-button size="small" type="warning" @click="toggleRecord(row)">启停</el-button>
            <el-button size="small" type="danger" @click="deleteRecord(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="recordTotal"
        :current-page="recordPage"
        :page-size="recordPageSize"
        @current-change="changeRecordPage"
      />
    </section>

    <el-dialog v-model="recordDialogVisible" :title="editingRecordId ? '编辑解析' : '新增解析'" width="560px">
      <el-form :model="recordForm" label-width="90px">
        <el-form-item label="主机记录">
          <el-input v-model="recordForm.name" placeholder="@ 或 www" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="recordForm.type" placeholder="请选择">
            <el-option label="A" value="A" />
            <el-option label="AAAA" value="AAAA" />
            <el-option label="CNAME" value="CNAME" />
            <el-option label="MX" value="MX" />
            <el-option label="TXT" value="TXT" />
            <el-option label="NS" value="NS" />
          </el-select>
        </el-form-item>
        <el-form-item label="记录值">
          <el-input v-model="recordForm.value" />
        </el-form-item>
        <el-form-item label="线路">
          <el-input v-model="recordForm.line" placeholder="默认线路可留空" />
        </el-form-item>
        <el-form-item label="TTL">
          <el-input v-model="recordForm.ttl" placeholder="例如 600" />
        </el-form-item>
        <el-form-item label="MX优先级">
          <el-input v-model="recordForm.mx" placeholder="仅 MX 记录需要" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="recordForm.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="recordSaving" @click="saveRecord">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'DomainManagement',
  data() {
    return {
      healthLoading: false,
      domainLoading: false,
      recordLoading: false,
      recordSaving: false,
      health: {},
      domains: [],
      records: [],
      selectedDomain: null,
      selectedRecords: [],
      domainKeyword: '',
      recordKeyword: '',
      domainPage: 1,
      domainPageSize: 20,
      domainTotal: 0,
      recordPage: 1,
      recordPageSize: 20,
      recordTotal: 0,
      recordDialogVisible: false,
      editingRecordId: '',
      recordForm: this.emptyRecordForm()
    }
  },
  computed: {
    healthText() {
      if (!this.health.enabled) return 'disabled'
      return this.health.healthy ? 'healthy' : 'unhealthy'
    },
    selectedDomainId() {
      return this.getDomainId(this.selectedDomain)
    },
    selectedDomainName() {
      return this.pick(this.selectedDomain, ['domain', 'Domain', 'name', 'Name', 'title', 'Title'])
    },
    selectedRecordIds() {
      return this.selectedRecords.map(row => this.getRecordId(row)).filter(Boolean)
    }
  },
  mounted() {
    this.loadHealth()
    this.loadDomains()
  },
  methods: {
    emptyRecordForm() {
      return {
        name: '',
        type: 'A',
        value: '',
        line: '',
        ttl: '',
        mx: '',
        remark: ''
      }
    },
    pick(row, keys) {
      if (!row) return ''
      for (const key of keys) {
        if (row[key] !== undefined && row[key] !== null && row[key] !== '') {
          return row[key]
        }
      }
      return ''
    },
    getDomainId(row) {
      return this.pick(row, ['id', 'ID', 'domain_id', 'domainId', 'DomainId'])
    },
    getRecordId(row) {
      return this.pick(row, ['recordid', 'record_id', 'recordId', 'RecordId', 'id', 'ID'])
    },
    assertSuccess(res) {
      if (res?.code && res.code !== 200) {
        throw new Error(res.message || 'dnsmgr 操作失败')
      }
      return res?.data || {}
    },
    async loadHealth() {
      this.healthLoading = true
      try {
        const { data: res } = await this.$api.domainHealth()
        this.health = this.assertSuccess(res)
      } finally {
        this.healthLoading = false
      }
    },
    async loadDomains() {
      this.domainLoading = true
      try {
        const { data: res } = await this.$api.domainList({
          page: this.domainPage,
          pageSize: this.domainPageSize,
          keyword: this.domainKeyword
        })
        const payload = this.assertSuccess(res)
        this.domains = payload.list || []
        this.domainTotal = payload.total || this.domains.length
      } finally {
        this.domainLoading = false
      }
    },
    async selectDomain(row) {
      this.selectedDomain = row
      this.recordPage = 1
      await this.loadRecords()
    },
    async showDomainDetail(row) {
      const id = this.getDomainId(row)
      const { data: res } = await this.$api.domainDetail(id)
      const detail = this.assertSuccess(res)
      ElMessageBox.alert(`<pre>${this.escapeHtml(JSON.stringify(detail, null, 2))}</pre>`, '域名详情', {
        dangerouslyUseHTMLString: true,
        customClass: 'domain-json-dialog'
      })
    },
    async openLoginUrl(row) {
      const id = this.getDomainId(row)
      const { data: res } = await this.$api.domainDetail(id, { loginUrl: 1 })
      const detail = this.assertSuccess(res)
      const loginUrl = this.pick(detail, ['loginurl', 'loginUrl', 'LoginUrl', 'url', 'Url'])
      if (!loginUrl) {
        ElMessage.warning('dnsmgr 未返回登录链接')
        return
      }
      window.open(loginUrl, '_blank')
    },
    async loadRecords() {
      if (!this.selectedDomainId) return
      this.recordLoading = true
      try {
        const { data: res } = await this.$api.domainRecords({
          domainId: this.selectedDomainId,
          page: this.recordPage,
          pageSize: this.recordPageSize,
          keyword: this.recordKeyword
        })
        const payload = this.assertSuccess(res)
        this.records = payload.list || []
        this.recordTotal = payload.total || this.records.length
      } finally {
        this.recordLoading = false
      }
    },
    openRecordDialog(row) {
      this.editingRecordId = row ? this.getRecordId(row) : ''
      this.recordForm = row ? {
        name: this.pick(row, ['name', 'Name', 'host', 'Host', 'rr', 'RR']),
        type: this.pick(row, ['type', 'Type', 'record_type']) || 'A',
        value: this.pick(row, ['value', 'Value', 'record', 'Record', 'content', 'Content']),
        line: this.pick(row, ['line', 'Line', 'linename', 'LineName']),
        ttl: `${this.pick(row, ['ttl', 'TTL']) || ''}`,
        mx: `${this.pick(row, ['mx', 'MX', 'priority', 'Priority']) || ''}`,
        remark: this.pick(row, ['remark', 'Remark', 'note', 'Note'])
      } : this.emptyRecordForm()
      this.recordDialogVisible = true
    },
    async saveRecord() {
      if (!this.selectedDomainId) return
      this.recordSaving = true
      try {
        const request = this.editingRecordId
          ? this.$api.domainRecordUpdate(this.selectedDomainId, this.editingRecordId, this.recordForm)
          : this.$api.domainRecordAdd(this.selectedDomainId, this.recordForm)
        const { data: res } = await request
        this.assertSuccess(res)
        ElMessage.success('保存成功')
        this.recordDialogVisible = false
        await this.loadRecords()
      } finally {
        this.recordSaving = false
      }
    },
    async deleteRecord(row) {
      await ElMessageBox.confirm('确认删除这条解析记录吗？', '删除确认', { type: 'warning' })
      const { data: res } = await this.$api.domainRecordDelete(this.selectedDomainId, this.getRecordId(row))
      this.assertSuccess(res)
      ElMessage.success('删除成功')
      await this.loadRecords()
    },
    async toggleRecord(row) {
      const current = `${this.pick(row, ['status', 'Status', 'state', 'State'])}`.toLowerCase()
      const nextStatus = ['1', 'true', 'enable', 'enabled', 'normal'].includes(current) ? '0' : '1'
      const { data: res } = await this.$api.domainRecordStatus(this.selectedDomainId, this.getRecordId(row), nextStatus)
      this.assertSuccess(res)
      ElMessage.success('状态已更新')
      await this.loadRecords()
    },
    async changeRemark(row) {
      const { value } = await ElMessageBox.prompt('请输入新的备注', '修改备注', {
        inputValue: this.pick(row, ['remark', 'Remark', 'note', 'Note'])
      })
      const { data: res } = await this.$api.domainRecordRemark(this.selectedDomainId, this.getRecordId(row), value)
      this.assertSuccess(res)
      ElMessage.success('备注已更新')
      await this.loadRecords()
    },
    async batchRecords(command) {
      if (!this.selectedDomainId || !this.selectedRecordIds.length) return
      let payload = { action: command, recordInfo: this.selectedRecords }
      if (command === 'remark') {
        const { value } = await ElMessageBox.prompt('请输入批量备注', '批量改备注')
        payload.remark = value
      } else {
        await ElMessageBox.confirm(`确认对 ${this.selectedRecordIds.length} 条记录执行 ${command} 吗？`, '批量确认', { type: 'warning' })
      }
      const { data: res } = await this.$api.domainRecordBatch(this.selectedDomainId, payload)
      this.assertSuccess(res)
      ElMessage.success('批量操作已提交')
      await this.loadRecords()
    },
    changeDomainPage(page) {
      this.domainPage = page
      this.loadDomains()
    },
    changeRecordPage(page) {
      this.recordPage = page
      this.loadRecords()
    },
    escapeHtml(text) {
      return text.replace(/[&<>"']/g, char => ({
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;'
      }[char]))
    }
  }
}
</script>

<style scoped>
.domain-page {
  display: grid;
  gap: 20px;
  padding: 24px;
}

.domain-hero,
.status-card,
.table-panel {
  background: #fff;
  border: 1px solid #e6edf5;
  border-radius: 18px;
  box-shadow: 0 14px 34px rgba(23, 38, 72, 0.08);
}

.domain-hero {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  padding: 28px;
  background: linear-gradient(135deg, #f6fbff 0%, #fff 52%, #f0f8ec 100%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f6f73;
  font-size: 12px;
  letter-spacing: 0.12em;
}

h1,
h2 {
  margin: 0;
  color: #1d2a44;
}

.subtitle,
.table-head p,
.status-card small {
  color: #64748b;
}

.subtitle {
  max-width: 780px;
  line-height: 1.7;
}

.actions,
.record-actions {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.status-card {
  display: grid;
  gap: 8px;
  padding: 18px;
}

.status-card span {
  color: #64748b;
}

.status-card strong {
  font-size: 24px;
  color: #1d2a44;
}

.status-card strong.ok {
  color: #16794c;
}

.status-card strong.warn {
  color: #b45309;
}

.table-panel {
  padding: 20px;
}

.table-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 16px;
}

.search {
  max-width: 260px;
}

.el-pagination {
  margin-top: 16px;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .domain-page {
    padding: 16px;
  }

  .domain-hero,
  .table-head,
  .record-actions {
    flex-direction: column;
  }

  .actions,
  .search {
    width: 100%;
  }
}
</style>

<style>
.domain-json-dialog pre {
  max-height: 520px;
  overflow: auto;
  text-align: left;
  white-space: pre-wrap;
}
</style>
