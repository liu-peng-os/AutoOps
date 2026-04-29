<template>
  <el-container class="home-container">
    <el-aside :width="isCollapse ? '64px' : '200px'">
      <div class="logo">
        <img src="../assets/image/DevOps平台.svg" class="siderbar-logo">
        <h2 v-show="!isCollapse">devops系统</h2>
      </div>
      <el-menu background-color="transparent" text-color="rgba(255,255,255,0.9)" active-text-color="#ffffff" router :default-active="$route.path"
               :collapse="isCollapse" :collapse-transition="false" class="modern-menu">
        <!--无子集菜单-->
        <el-menu-item :index="'/' + item.url" v-for="item in noChildren" :key="item.menuName" @click="saveNavState('/' + item.url)">
          <el-icon><component :is="item.icon" /></el-icon>
          <template v-slot:title>
            <span>{{ item.menuName }}</span>
          </template>
        </el-menu-item>
        <!--有子集菜单-->
        <el-sub-menu :index="item.id + ''" v-for="item in hasChildren" :key="item.id">
          <template #title>
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.menuName }}</span>
          </template>
          <el-menu-item :index="'/' + subItem.url" v-for="subItem in item.menuSvoList" :key="subItem.id"
                        @click="saveNavState('/' + subItem.url)">
            <el-icon><component :is="subItem.icon" /></el-icon>
            <template #title>
              <span>{{ subItem.menuName }}</span>
            </template>
          </el-menu-item>
        </el-sub-menu>

      </el-menu>
    </el-aside>

    <!-- 主体内容 -->
    <el-container>
      <el-header height="50px">
        <!-- 顶部导航栏,折叠图标 -->
        <div class="fold-btn">
          <el-button type="text" @click="toggleCollapse" class="collapse-btn">
            <el-icon size="24"><component :is="collapseBtnClass" /></el-icon>
          </el-button>
        </div>
        <div class="bread-btn">
          <!-- 面包屑 -->
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">仪表盘</el-breadcrumb-item>

            <!-- 显示二级标题 -->
            <el-breadcrumb-item v-if="$route.meta && $route.meta.sTitle">
              {{ $route.meta.sTitle }}
            </el-breadcrumb-item>

            <!-- 显示三级标题 -->
            <el-breadcrumb-item v-if="$route.meta && $route.meta.tTitle">
              {{ $route.meta.tTitle }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <HeadImage />
      </el-header>
      <Tags />
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script>


import storage from "@/utils/storage";
import HeadImage   from "@/components/HeadImage.vue";
import Tags from "@/components/Tags.vue";

export default {
  // eslint-disable-next-line vue/multi-word-component-names
  name: "Home",
  components: { HeadImage, Tags },
  data() {
    return {
      leftMenuList: null, // 初始化为null，在mounted中设置
      activePath: '',
      collapseBtnClass: "Fold",
      isCollapse: false,
    }
  },
  computed: {
    // 无子集
    noChildren() {
      return (this.leftMenuList || []).filter(item => !item.menuSvoList) // 过滤无子集
    },
    // 有子集
    hasChildren() {
      return (this.leftMenuList || []).filter(item => item.menuSvoList)
    }
  },
  methods: {
    ensureMenuGroup(menuList, group) {
      let target = menuList.find(item => item.menuName === group.menuName)
      if (!target) {
        target = {
          id: group.id,
          menuName: group.menuName,
          icon: group.icon,
          menuSvoList: []
        }
        menuList.push(target)
      } else if (!Array.isArray(target.menuSvoList)) {
        target.menuSvoList = []
      }
      return target
    },
    ensureChildMenu(group, child) {
      const exists = group.menuSvoList.some(item => item.url === child.url)
      if (!exists) {
        group.menuSvoList.push(child)
      }
    },
    bootstrapPlannedMenus(menuList) {
      const domainGroup = this.ensureMenuGroup(menuList, {
        id: 91000,
        menuName: '域名管理',
        icon: 'Link'
      })
      this.ensureChildMenu(domainGroup, {
        id: 91001,
        menuName: '统一入口',
        url: 'integration/domain',
        icon: 'Link'
      })

      const workorderGroup = this.ensureMenuGroup(menuList, {
        id: 92000,
        menuName: '运营工单',
        icon: 'Tickets'
      })
      this.ensureChildMenu(workorderGroup, {
        id: 92001,
        menuName: '工单中心',
        url: 'integration/workorder',
        icon: 'Tickets'
      })

      const assetGroup = this.ensureMenuGroup(menuList, {
        id: 93000,
        menuName: '资产管理',
        icon: 'Box'
      })
      this.ensureChildMenu(assetGroup, {
        id: 93001,
        menuName: '站点管理',
        url: 'cmdb/site',
        icon: 'OfficeBuilding'
      })
    },
    // 初始化菜单数据
    initMenuData() {
      try {
        const menuData = storage.getItem("leftMenuList");
        console.log('初始化菜单数据:', menuData);

        // 确保数据是数组格式
        if (Array.isArray(menuData)) {
          this.leftMenuList = menuData;
          this.bootstrapPlannedMenus(this.leftMenuList)
          
          // 手动添加配置管理菜单到任务中心
          const taskMenu = this.leftMenuList.find(item => item.menuName === '任务中心')
          if (taskMenu && taskMenu.menuSvoList) {
            const configExists = taskMenu.menuSvoList.some(sub => sub.url === 'task/config')
            if (!configExists) {
              taskMenu.menuSvoList.push({
                id: 99999,
                menuName: '配置管理',
                url: 'task/config',
                icon: 'Setting'
              })
            }
          }
        } else if (menuData) {
          // 如果数据存在但不是数组，尝试解析
          console.warn('菜单数据格式异常，尝试修复:', menuData);
          this.leftMenuList = [];
          this.bootstrapPlannedMenus(this.leftMenuList)
        } else {
          // 如果没有数据，设为空数组
          console.warn('未找到菜单数据，使用空数组');
          this.leftMenuList = [];
          this.bootstrapPlannedMenus(this.leftMenuList)
        }

        // 强制触发视图更新
        this.$forceUpdate();
      } catch (error) {
        console.error('初始化菜单数据失败:', error);
        this.leftMenuList = [];
      }
    },
    // 点击实现跳转
    saveNavState(activePath) {
      storage.setItem('activePath', activePath)
      this.activePath = activePath
    },
    // 张开和折叠
    toggleCollapse() {
      this.isCollapse = !this.isCollapse
      if (this.isCollapse) {
        this.collapseBtnClass = 'Fold'  // 折叠
      } else {
        this.collapseBtnClass = 'Expand'  // 展开
      }
    },
    // 移除菜单项的focus效果
    removeFocusOutline() {
      this.$nextTick(() => {
        const menuItems = document.querySelectorAll('.el-menu-item, .el-sub-menu__title')
        menuItems.forEach(item => {
          item.style.outline = 'none'
          item.style.border = 'none'
          item.addEventListener('focus', (e) => {
            e.target.style.outline = 'none'
            e.target.style.border = 'none'
            e.target.blur()
          })
          item.addEventListener('click', (e) => {
            e.target.style.outline = 'none'
            e.target.style.border = 'none'
            setTimeout(() => {
              e.target.blur()
            }, 0)
          })
        })
      })
    }
  },
  mounted() {
    // 确保在mounted阶段重新获取菜单数据，解决浏览器兼容性问题
    this.initMenuData();
    this.removeFocusOutline();
  }
}
</script>

<style lang="less" scoped>
.home-container {
  height: 100%;

  .el-aside {
    background: linear-gradient(135deg, var(--color-sidebar-bg) 0%, var(--color-sidebar-deep) 100%);
    border-right: 1px solid var(--color-sidebar-hover);
    box-shadow: 2px 0 10px oklch(0% 0 0 / 0.08);

    .logo {
      margin-top: 8px;
      display: flex;
      align-items: center;
      font-size: 14px;
      height: 50px;
      color: var(--color-sidebar-text);
      font-weight: 500;
      padding: 8px 14px;
      white-space: nowrap;
      overflow: hidden;

      .siderbar-logo {
        width: 40px;
        height: 32px;
        margin-right: 10px;
        flex-shrink: 0;
        border-radius: 4px;
      }

      h2 {
        margin: 0;
        font-weight: 700;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
        letter-spacing: 0.5px;
        font-size: 17px;
        color: var(--color-sidebar-text);
        flex-shrink: 0;
        min-width: 0;
      }
    }

    .modern-menu {
      border-right: none;
      background: transparent !important;
    }
  }

  .modern-menu {
    .el-menu-item {
      transition: background var(--duration-fast) var(--ease-out),
                  color var(--duration-fast) var(--ease-out);
      border-radius: var(--radius-md);
      margin: 2px 8px;
      color: var(--color-sidebar-text) !important;
      outline: none !important;

      &:hover {
        background: var(--color-sidebar-hover) !important;
        color: oklch(100% 0 0) !important;
      }

      &:focus-visible {
        outline: 2px solid var(--color-accent-muted) !important;
        outline-offset: -2px;
      }
    }

    > .el-menu-item.is-active {
      background: linear-gradient(135deg, var(--color-sidebar-active-from) 0%, var(--color-sidebar-active-to) 100%) !important;
      color: oklch(100% 0 0) !important;
      box-shadow: 0 2px 8px oklch(58% 0.18 265 / 0.35);
    }

    .el-sub-menu {
      .el-sub-menu__title {
        transition: background var(--duration-fast) var(--ease-out);
        border-radius: var(--radius-md);
        margin: 2px 8px;
        color: var(--color-sidebar-text) !important;
        outline: none !important;

        &:hover {
          background: var(--color-sidebar-hover) !important;
          color: oklch(100% 0 0) !important;
        }
      }

      .el-menu-item.is-active {
        background: oklch(58% 0.18 265 / 0.85) !important;
        color: oklch(100% 0 0) !important;
        border-radius: var(--radius-md);
        margin: 2px 12px 2px 20px;
        width: calc(100% - 32px);
        box-shadow: 0 2px 6px oklch(58% 0.18 265 / 0.25);
      }
    }

    > .el-sub-menu.is-active,
    > .el-sub-menu.is-opened {
      background-color: transparent !important;

      .el-sub-menu__title {
        background: oklch(100% 0 0 / 0.05) !important;
        border-radius: var(--radius-md);
      }
    }
  }

  .el-header {
    background: var(--color-header-bg);
    border-bottom: 1px solid var(--color-header-border);
    box-shadow: var(--shadow-xs);
    align-items: center;
    justify-content: space-between;
    display: flex;

    .fold-btn {
      font-size: 23px;
      cursor: pointer;

      .collapse-btn {
        padding: 6px;
        border-radius: var(--radius-md);
        transition: background var(--duration-fast) var(--ease-out);
        color: var(--color-accent);

        &:hover {
          background: var(--color-accent-subtle);
        }
      }
    }

    .bread-btn {
      position: fixed;
      margin-left: 40px;

      .el-breadcrumb {
        .el-breadcrumb__item {
          .el-breadcrumb__inner {
            color: var(--color-accent);
            font-weight: 500;
            transition: color var(--duration-fast) var(--ease-out);

            &:hover {
              color: var(--color-accent-hover);
            }
          }

          &:last-child .el-breadcrumb__inner {
            color: var(--color-text-secondary);
          }
        }
      }
    }
  }

  .el-main {
    background: var(--color-bg);
  }
}
</style>

<style lang="less">
.el-menu--popup-bottom-start,
.el-menu--popup {
  background: linear-gradient(160deg, var(--color-sidebar-bg) 0%, var(--color-sidebar-deep) 100%) !important;
  border: 1px solid oklch(100% 0 0 / 0.08) !important;
  box-shadow: var(--shadow-lg) !important;
  border-radius: var(--radius-lg) !important;

  .el-menu-item {
    color: var(--color-sidebar-text) !important;
    background: transparent !important;
    transition: background var(--duration-fast) !important;
    margin: 2px 6px !important;
    border-radius: var(--radius-md) !important;

    &:hover {
      background: oklch(100% 0 0 / 0.10) !important;
      color: oklch(100% 0 0) !important;
    }

    &.is-active {
      background: linear-gradient(135deg, var(--color-sidebar-active-from) 0%, var(--color-sidebar-active-to) 100%) !important;
      color: oklch(100% 0 0) !important;
    }

    .el-icon, span {
      color: inherit !important;
    }
  }
}
</style>
