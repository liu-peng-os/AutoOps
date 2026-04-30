<template>
  <div class="domain-page">
    <section class="domain-hero">
      <div>
        <p class="eyebrow">DNSMGR ADAPTER</p>
        <h1>Domain Monitor</h1>
        <p class="subtitle">
          Read-only domain inventory from dnsmgr. AutoOps stores normalized zones and DNS records for monitoring first,
          and keeps write-back disabled until the adapter is proven safe.
        </p>
      </div>
      <div class="actions">
        <el-button :loading="healthLoading" @click="loadHealth">Health Check</el-button>
        <el-button type="primary" :loading="syncLoading" @click="syncDomain">Sync Now</el-button>
      </div>
    </section>

    <section class="status-grid">
      <article class="status-card">
        <span>dnsmgr</span>
        <strong :class="health.healthy ? 'ok' : 'warn'">{{ healthText }}</strong>
        <small>{{ health.message || 'No health check yet' }}</small>
      </article>
      <article class="status-card">
        <span>Last sync</span>
        <strong>{{ lastSync.status || 'none' }}</strong>
        <small>{{ lastSync.message || 'No sync run recorded' }}</small>
      </article>
      <article class="status-card">
        <span>Synced objects</span>
        <strong>{{ lastSync.zonesSynced || 0 }} / {{ lastSync.recordsSynced || 0 }}</strong>
        <small>zones / records</small>
      </article>
    </section>

    <section class="table-panel">
      <div class="table-head">
        <div>
          <h2>Zones</h2>
          <p>Normalized domain zones from dnsmgr.</p>
        </div>
        <el-input
          v-model="zoneKeyword"
          clearable
          placeholder="Search domain"
          class="search"
          @keyup.enter="loadZones"
          @clear="loadZones"
        />
      </div>
      <el-table v-loading="zoneLoading" :data="zones" border @row-click="selectZone">
        <el-table-column prop="name" label="Domain" min-width="220" />
        <el-table-column prop="provider" label="Provider" min-width="120" />
        <el-table-column prop="status" label="Status" width="120" />
        <el-table-column prop="lastSyncedAt" label="Last Synced" min-width="180" />
      </el-table>
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="zoneTotal"
        :current-page="zonePage"
        :page-size="zonePageSize"
        @current-change="changeZonePage"
      />
    </section>

    <section class="table-panel">
      <div class="table-head">
        <div>
          <h2>DNS Records</h2>
          <p>{{ selectedZone ? `Records for ${selectedZone.name}` : 'Click a zone to filter records.' }}</p>
        </div>
        <el-input
          v-model="recordKeyword"
          clearable
          placeholder="Search record or value"
          class="search"
          @keyup.enter="loadRecords"
          @clear="loadRecords"
        />
      </div>
      <el-table v-loading="recordLoading" :data="records" border>
        <el-table-column prop="zoneName" label="Zone" min-width="180" />
        <el-table-column prop="name" label="Name" min-width="160" />
        <el-table-column prop="type" label="Type" width="90" />
        <el-table-column prop="value" label="Value" min-width="260" show-overflow-tooltip />
        <el-table-column prop="line" label="Line" width="120" />
        <el-table-column prop="ttl" label="TTL" width="90" />
        <el-table-column prop="status" label="Status" width="120" />
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
  </div>
</template>

<script>
export default {
  name: 'DomainManagement',
  data() {
    return {
      healthLoading: false,
      syncLoading: false,
      zoneLoading: false,
      recordLoading: false,
      health: {},
      lastSync: {},
      zones: [],
      records: [],
      selectedZone: null,
      zoneKeyword: '',
      recordKeyword: '',
      zonePage: 1,
      zonePageSize: 20,
      zoneTotal: 0,
      recordPage: 1,
      recordPageSize: 20,
      recordTotal: 0
    }
  },
  computed: {
    healthText() {
      if (!this.health.enabled) return 'disabled'
      return this.health.healthy ? 'healthy' : 'unhealthy'
    }
  },
  mounted() {
    this.loadHealth()
    this.loadLastSync()
    this.loadZones()
    this.loadRecords()
  },
  methods: {
    async loadHealth() {
      this.healthLoading = true
      try {
        const { data: res } = await this.$api.domainHealth()
        this.health = res?.data || {}
      } finally {
        this.healthLoading = false
      }
    },
    async loadLastSync() {
      const { data: res } = await this.$api.domainLastSync()
      this.lastSync = res?.data || {}
    },
    async syncDomain() {
      this.syncLoading = true
      try {
        const { data: res } = await this.$api.domainSync()
        this.lastSync = res?.data || {}
        await this.loadZones()
        await this.loadRecords()
      } finally {
        this.syncLoading = false
      }
    },
    async loadZones() {
      this.zoneLoading = true
      try {
        const { data: res } = await this.$api.domainZones({
          page: this.zonePage,
          pageSize: this.zonePageSize,
          keyword: this.zoneKeyword
        })
        const payload = res?.data || {}
        this.zones = payload.list || []
        this.zoneTotal = payload.total || 0
      } finally {
        this.zoneLoading = false
      }
    },
    async loadRecords() {
      this.recordLoading = true
      try {
        const { data: res } = await this.$api.domainRecords({
          page: this.recordPage,
          pageSize: this.recordPageSize,
          keyword: this.recordKeyword,
          zoneId: this.selectedZone?.id || undefined
        })
        const payload = res?.data || {}
        this.records = payload.list || []
        this.recordTotal = payload.total || 0
      } finally {
        this.recordLoading = false
      }
    },
    selectZone(row) {
      this.selectedZone = row
      this.recordPage = 1
      this.loadRecords()
    },
    changeZonePage(page) {
      this.zonePage = page
      this.loadZones()
    },
    changeRecordPage(page) {
      this.recordPage = page
      this.loadRecords()
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
  background: linear-gradient(135deg, #f5fbff 0%, #fff 52%, #edf8f2 100%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #3a6ea5;
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

.actions {
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
  .table-head {
    flex-direction: column;
  }

  .actions,
  .search {
    width: 100%;
  }
}
</style>
