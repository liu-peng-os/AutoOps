<template>
  <div class="integration-landing">
    <section class="hero-card">
      <div class="hero-copy">
        <p class="eyebrow">阶段 1 占位入口</p>
        <h1>{{ title }}</h1>
        <p class="subtitle">{{ subtitle }}</p>
      </div>
      <div class="status-chip">可插拔接入优先</div>
    </section>

    <section class="grid">
      <article class="panel">
        <h2>当前模块定位</h2>
        <ul>
          <li>先建立稳定入口，不急着把业务细节一次定死。</li>
          <li>外部系统作为可替换依赖处理，平台内部保留通用模型。</li>
          <li>后续实现会优先复用同步、回写、审计和适配器骨架。</li>
        </ul>
      </article>

      <article class="panel">
        <h2>阶段状态</h2>
        <div class="meta-row">
          <span>当前状态</span>
          <strong>{{ moduleInfo.statusLabel }}</strong>
        </div>
        <div class="meta-row">
          <span>数据来源</span>
          <strong>{{ moduleInfo.sourceLabel }}</strong>
        </div>
        <div class="meta-row">
          <span>支持回写</span>
          <strong>{{ moduleInfo.writebackLabel }}</strong>
        </div>
        <p class="notes">{{ moduleInfo.notes }}</p>
      </article>
    </section>

    <section class="panel">
      <h2>平台原则</h2>
      <ul>
        <li v-for="item in principles" :key="item">{{ item }}</li>
      </ul>
    </section>

    <section class="panel">
      <div class="section-head">
        <h2>已配置的外部系统</h2>
        <button class="refresh-btn" @click="loadCatalog">刷新</button>
      </div>
      <p v-if="loading" class="muted">正在读取接入目录...</p>
      <p v-else-if="systems.length === 0" class="muted">
        当前还没有在配置文件里声明外部系统，这不影响先做骨架。
      </p>
      <div v-else class="system-list">
        <article v-for="system in systems" :key="system.key" class="system-card">
          <div class="system-head">
            <h3>{{ system.displayName || system.key }}</h3>
            <span :class="['badge', system.enabled ? 'enabled' : 'disabled']">
              {{ system.enabled ? 'enabled' : 'disabled' }}
            </span>
          </div>
          <p class="muted">{{ system.category || 'external-system' }} / {{ system.provider || 'provider-pending' }}</p>
          <p class="url">{{ system.baseUrl || '未配置 Base URL' }}</p>
          <div class="capabilities">
            <span v-for="capability in system.capabilities || []" :key="capability" class="capability">
              {{ capability }}
            </span>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script>
const fallbackModules = {
  'domain-management': {
    sourceLabel: '外部 DNS 系统统一入口',
    statusLabel: 'planned',
    writebackLabel: '后续按接入能力决定',
    notes: '短期先接入成熟外部系统，后续如有必要再增强平台内的汇总与联动能力。'
  },
  'operations-workorder': {
    sourceLabel: '外部表格主数据源',
    statusLabel: 'planned',
    writebackLabel: '支持',
    notes: '会优先围绕同步、展示、基础编辑和回写闭环来设计。'
  },
  'site-management': {
    sourceLabel: '外部表格主数据源',
    statusLabel: 'planned',
    writebackLabel: '支持',
    notes: '挂在资产管理下，先做统一入口和数据同步回写闭环。'
  }
}

export default {
  name: 'ExternalModuleLanding',
  props: {
    pageKey: {
      type: String,
      required: true
    },
    title: {
      type: String,
      required: true
    },
    subtitle: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      loading: false,
      systems: [],
      principles: [
        '所有外部系统都作为可替换依赖处理。',
        '平台内部优先保留通用业务模型，而不是第三方专有字段。',
        '模块优先复用统一适配层、同步任务和审计能力。'
      ],
      moduleInfo: fallbackModules[this.pageKey] || fallbackModules['domain-management']
    }
  },
  mounted() {
    this.loadCatalog()
  },
  methods: {
    async loadCatalog() {
      this.loading = true
      try {
        const { data: res } = await this.$api.integrationCatalog()
        const payload = res?.data || {}
        const modules = payload.modules || []
        const matchedModule = modules.find(item => item.key === this.pageKey)
        if (matchedModule) {
          this.moduleInfo = {
            sourceLabel: matchedModule.source || this.moduleInfo.sourceLabel,
            statusLabel: matchedModule.status || this.moduleInfo.statusLabel,
            writebackLabel: matchedModule.supportsWriteback ? '支持' : '暂未启用',
            notes: matchedModule.notes || this.moduleInfo.notes
          }
        }
        this.principles = payload.principles || this.principles
        this.systems = payload.systems || []
      } catch (error) {
        console.warn('读取外部接入目录失败，继续使用本地占位信息', error)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.integration-landing {
  padding: 24px;
  display: grid;
  gap: 20px;
}

.hero-card,
.panel,
.system-card {
  background: #ffffff;
  border: 1px solid #e8edf5;
  border-radius: 18px;
  box-shadow: 0 10px 30px rgba(30, 42, 90, 0.08);
}

.hero-card {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  padding: 28px;
  background: linear-gradient(135deg, #f4f8ff 0%, #ffffff 55%, #eef6f2 100%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #5271a5;
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

h1,
h2,
h3 {
  margin: 0;
  color: #20314f;
}

.subtitle,
.muted,
.notes,
.url {
  color: #63748a;
}

.subtitle {
  margin: 10px 0 0;
  max-width: 680px;
  line-height: 1.7;
}

.status-chip {
  height: fit-content;
  padding: 10px 14px;
  background: #20314f;
  color: #fff;
  border-radius: 999px;
  font-size: 12px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.panel {
  padding: 22px;
}

.panel ul {
  margin: 14px 0 0;
  padding-left: 18px;
  color: #334155;
  line-height: 1.7;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 10px 0;
  border-bottom: 1px solid #edf2f7;
  color: #475569;
}

.meta-row:last-of-type {
  border-bottom: none;
}

.notes {
  margin: 14px 0 0;
  line-height: 1.7;
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.refresh-btn {
  border: 1px solid #d6dfeb;
  background: #fff;
  color: #20314f;
  padding: 8px 14px;
  border-radius: 999px;
  cursor: pointer;
}

.system-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.system-card {
  padding: 18px;
}

.system-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.badge {
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.enabled {
  background: #e7f8ee;
  color: #1f7a46;
}

.disabled {
  background: #f4f4f5;
  color: #667085;
}

.capabilities {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.capability {
  background: #eef4ff;
  color: #35599a;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
}

@media (max-width: 768px) {
  .integration-landing {
    padding: 16px;
  }

  .hero-card {
    flex-direction: column;
    padding: 22px;
  }
}
</style>
