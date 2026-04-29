<template>
    <div class="login-container">
        <div class="login-card">
            <div class="card-header">
                <h2 class="card-title">AutoOps</h2>
                <span class="card-subtitle">运维管理系统</span>
            </div>

            <!-- 表单 -->
            <el-form ref="loginFormRef" :rules="rules" :model="loginForm">
                <el-form-item prop="username">
                    <el-input
                        v-model="loginForm.username"
                        placeholder="请输入账号"
                        clearable
                        class="dark-input"
                    />
                </el-form-item>

                <el-form-item prop="password">
                    <el-input
                        v-model="loginForm.password"
                        placeholder="请输入密码"
                        type="password"
                        show-password
                        clearable
                        class="dark-input"
                    />
                </el-form-item>

                <el-form-item prop="image">
                    <div class="captcha-row">
                        <el-input
                            v-model="loginForm.image"
                            placeholder="请输入验证码"
                            maxlength="6"
                            clearable
                            class="dark-input"
                        />
                        <div class="captcha-box" @click="getCaptcha">
                            <el-image :src="image" class="captcha-img" />
                        </div>
                    </div>
                </el-form-item>

                <el-form-item>
                    <el-button class="login-btn" type="primary" @click="loginBtn">登 录</el-button>
                    <el-button class="reset-btn" @click="resetLoginForm">重 置</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
export default {
    name: "Login",
    data() {
        return {
            image: '',
            rules: {
                username: [{ required: true, message: "请输入账号", trigger: "blur" }],
                password: [{ required: true, message: "请输入密码", trigger: "blur" }],
                image:    [{ required: true, message: "请输入验证码", trigger: "blur" }]
            },
            loginForm: {
                username: '',
                password: '',
                image: '',
                idKey: ''
            }
        }
    },
    methods: {
        async getCaptcha() {
            const { data: res } = await this.$api.captcha()
            if (res.code !== 200) {
                this.$message.error(res.message)
            } else {
                this.image = res.data.image
                this.loginForm.idKey = res.data.idKey
            }
        },
        loginBtn() {
            this.$refs.loginFormRef.validate(async valid => {
                if (valid) {
                    const { data: res } = await this.$api.login(this.loginForm)
                    if (res.code !== 200) {
                        this.$message.error(res.message)
                    } else {
                        this.$message.success("登录成功")
                        this.$store.commit('saveSysAdmin', res.data.sysAdmin)
                        this.$store.commit('saveToken', res.data.token)
                        this.$store.commit('saveLeftMenuList', res.data.leftMenuList)
                        this.$store.commit('savePermissionList', res.data.permissionList)
                        await this.$router.push("/home")
                    }
                } else {
                    return false
                }
            })
        },
        resetLoginForm() {
            this.$refs.loginFormRef.resetFields()
        }
    },
    created() {
        this.getCaptcha()
    }
}
</script>

<style lang="less" scoped>
.login-container {
    width: 100%;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: url('../assets/image/背景.jpg') center / cover no-repeat;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', system-ui, sans-serif;
}

.login-card {
    position: relative;
    z-index: 1;
    width: 380px;
    padding: 36px 32px 28px;
    background: oklch(14% 0.04 265);
    border: 1px solid oklch(100% 0 0 / 0.08);
    border-radius: var(--radius-xl);
    box-shadow: 0 20px 48px oklch(0% 0 0 / 0.55);
    text-align: center;
}

.card-header {
    display: flex;
    align-items: baseline;
    justify-content: center;
    gap: 8px;
    margin-bottom: 28px;
}

.card-title {
    font-size: 20px;
    font-weight: 700;
    color: oklch(94% 0.01 265);
    margin: 0;
    letter-spacing: 0.5px;
}

.card-subtitle {
    font-size: 14px;
    font-weight: 400;
    color: oklch(62% 0.02 265);
    letter-spacing: 0;
}

.dark-input {
    :deep(.el-input__wrapper) {
        background: oklch(20% 0.04 265) !important;
        border: 1px solid oklch(100% 0 0 / 0.10) !important;
        border-radius: var(--radius-md) !important;
        box-shadow: none !important;
        transition: border-color var(--duration-fast) var(--ease-out) !important;

        &:hover {
            border-color: oklch(100% 0 0 / 0.18) !important;
        }

        &.is-focus {
            border-color: var(--color-accent) !important;
            box-shadow: 0 0 0 2px oklch(58% 0.18 265 / 0.20) !important;
        }
    }

    :deep(.el-input__inner) {
        color: oklch(90% 0.01 265) !important;
        font-size: 14px !important;
        height: 40px !important;
        line-height: 40px !important;

        &::placeholder {
            color: oklch(50% 0.02 265) !important;
        }
    }

    :deep(.el-input__prefix-inner),
    :deep(.el-input__suffix-inner) {
        color: oklch(50% 0.02 265);
    }
}

.captcha-row {
    display: flex;
    gap: 10px;

    .dark-input { flex: 1; }
}

.captcha-box {
    flex-shrink: 0;
    width: 108px;
    height: 40px;
    border-radius: var(--radius-md);
    overflow: hidden;
    cursor: pointer;
    border: 1px solid oklch(100% 0 0 / 0.10);
    transition: border-color var(--duration-fast) var(--ease-out);

    &:hover { border-color: var(--color-accent); }
}

.captcha-img {
    width: 100%;
    height: 100%;
    display: block;
}

:deep(.el-form-item) {
    margin-bottom: 16px;

    &:last-child { margin-bottom: 0; }

    .el-form-item__error {
        font-size: 11px;
        color: oklch(65% 0.18 25);
        padding-top: 3px;
    }
}

.login-btn {
    width: calc(100% - 96px);
    height: 40px;
    border: none;
    border-radius: var(--radius-md);
    background: var(--color-accent) !important;
    color: oklch(100% 0 0) !important;
    font-size: 14px !important;
    font-weight: 600 !important;
    letter-spacing: 2px;
    box-shadow: 0 2px 8px oklch(58% 0.18 265 / 0.30) !important;
    transition: background var(--duration-fast) var(--ease-out),
                box-shadow var(--duration-fast) var(--ease-out) !important;

    &:hover {
        background: var(--color-accent-hover) !important;
        box-shadow: 0 4px 16px oklch(58% 0.18 265 / 0.45) !important;
    }
}

.reset-btn {
    width: 82px;
    height: 40px;
    margin-left: 10px !important;
    border-radius: var(--radius-md);
    border: 1px solid oklch(100% 0 0 / 0.10) !important;
    background: transparent !important;
    color: oklch(60% 0.02 265) !important;
    font-size: 14px !important;
    letter-spacing: 0.5px;
    transition: border-color var(--duration-fast),
                color var(--duration-fast),
                background var(--duration-fast) !important;

    &:hover {
        border-color: oklch(100% 0 0 / 0.20) !important;
        color: oklch(80% 0.01 265) !important;
        background: oklch(100% 0 0 / 0.05) !important;
    }
}
</style>
