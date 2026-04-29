<template>
  <div class="tags-bar">
    <div
      v-for="(item, index) in tags"
      :key="item.path"
      :class="['tag-item', { 'is-active': item.title === $route.meta.tTitle }]"
      @click="goTo(item.path)"
    >
      <span class="tag-dot" v-show="item.title === $route.meta.tTitle"></span>
      <span class="tag-label">{{ item.title }}</span>
      <button
        v-if="index > 0"
        class="tag-close"
        @click.stop="close(index)"
        aria-label="关闭"
      >
        <svg viewBox="0 0 12 12" fill="currentColor" width="10" height="10">
          <path d="M1 1l10 10M11 1L1 11" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script>
import index from "vuex";

export default {
  name: "Tags",
  computed: {
    index() {
      return index
    }
  },
  data() {
    return {
      tags: [{
        path: "/dashboard",
        title: "仪表盘",
      }]
    }
  },
  watch: {
    $route: {
      immediate: true,
      handler(val) {
        const exists = this.tags.find(item => val.path === item.path)
        if (!exists) {
          this.tags.push({
            title: val.meta.tTitle,
            path: val.path
          })
        }
      }
    }
  },
  methods: {
    goTo(path) {
      this.$router.push(path)
    },
    close(i) {
      this.tags.splice(i, 1)
    }
  }
}
</script>

<style scoped>
.tags-bar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 16px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border-subtle);
  overflow-x: auto;
  scrollbar-width: none;
}

.tags-bar::-webkit-scrollbar {
  display: none;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 3px 10px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  color: var(--color-text-secondary);
  background: transparent;
  border: 1px solid transparent;
  cursor: pointer;
  white-space: nowrap;
  transition: background var(--duration-fast) var(--ease-out),
              color var(--duration-fast) var(--ease-out),
              border-color var(--duration-fast) var(--ease-out);
  user-select: none;
}

.tag-item:hover {
  background: var(--color-accent-subtle);
  color: var(--color-accent);
  border-color: var(--color-accent-muted);
}

.tag-item.is-active {
  background: var(--color-accent-subtle);
  color: var(--color-accent);
  border-color: var(--color-accent-muted);
  font-weight: 500;
}

.tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-accent);
  flex-shrink: 0;
}

.tag-label {
  line-height: 1;
}

.tag-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  border: none;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  border-radius: 2px;
  padding: 0;
  transition: background var(--duration-fast), color var(--duration-fast);
  flex-shrink: 0;
}

.tag-close:hover {
  background: var(--color-accent-muted);
  color: var(--color-accent);
}
</style>
