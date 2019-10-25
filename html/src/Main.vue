<template>
    <el-container class="body">
        <el-aside width="auto">
            <el-menu :collapse="isCollapse"
                     default-active="/"
                     class="el-menu-vertical"
                     background-color="#141f29"
                     text-color="#fff"
                     :default-openeds="['/1','/2','/3']"
                     router
                     active-text-color="#67C23A">
                <div class="logo">

                </div>
                <el-menu-item index="/">
                    <i class="el-icon-menu"></i>
                    <span slot="title">总览</span>
                </el-menu-item>
                <el-submenu index="/1">
                    <template slot="title">
                        <i class="el-icon-s-platform"></i>
                        <span slot="title">数据中心</span>
                    </template>
                    <el-menu-item index="/1-1">
                        <i class="el-icon-s-flag"></i>数据来源
                    </el-menu-item>
                    <el-menu-item index="/1-2">
                        <i class="el-icon-s-shop"></i>数据应用
                    </el-menu-item>
                    <el-menu-item index="/1-3">
                        <i class="el-icon-s-grid"></i>容器服务
                    </el-menu-item>
                </el-submenu>
                <el-submenu index="/2">
                    <template slot="title">
                        <i class="el-icon-coin"></i>
                        <span slot="title">集群管理</span>
                    </template>
                    <el-menu-item index="/2-1">
                        <i class="el-icon-mobile"></i>裸机管理
                    </el-menu-item>
                    <el-menu-item index="/2-2">
                        <i class="el-icon-success"></i>存储管理
                    </el-menu-item>
                    <el-menu-item index="/2-3">
                        <i class="el-icon-s-finance"></i>镜像管理
                    </el-menu-item>
                    <el-menu-item index="/2-4">
                        <i class="el-icon-document-delete"></i>网卡管理
                    </el-menu-item>
                </el-submenu>
                <el-submenu index="/3">
                    <template slot="title">
                        <i class="el-icon-eleme"></i>
                        <span slot="title">云主机</span>
                    </template>
                    <el-menu-item index="/2-1">
                        <i class="el-icon-folder-add"></i>云主机
                    </el-menu-item>
                    <el-menu-item index="/2-2">
                        <i class="el-icon-sort"></i>防火墙
                    </el-menu-item>
                    <el-menu-item index="/2-3">
                        <i class="el-icon-document-copy"></i>磁盘管理
                    </el-menu-item>
                </el-submenu>
                <div :class="navBottom">
                    <el-button size="mini" type="success" :icon="elIcon" round @click="switchMenu"></el-button>
                </div>
            </el-menu>
        </el-aside>
        <el-container>
            <el-header class="header" height="48px">
                <el-dropdown size="small" split-button type="success">
                    功能
                    <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item>消息</el-dropdown-item>
                        <el-dropdown-item>工单</el-dropdown-item>
                        <el-dropdown-item @click.native="exit">退出</el-dropdown-item>
                    </el-dropdown-menu>
                </el-dropdown>
            </el-header>
            <el-main>
                <router-view></router-view>
            </el-main>
        </el-container>
    </el-container>
</template>

<script>
    export default {
        name: 'app',
        components: {},
        data() {
            return {
                isCollapse: false,
                navBottom: 'nav-bottom-close',
                elIcon: 'el-icon-back',
            }
        },
        methods: {
            openMenu() {
                this.$data.isCollapse = true;
                this.$data.navBottom = 'nav-bottom-open';
                this.$data.elIcon = 'el-icon-right';
            },
            hiddenMenu() {
                this.$data.isCollapse = false;
                this.$data.navBottom = 'nav-bottom-close';
                this.$data.elIcon = 'el-icon-back';
            },
            switchMenu() {
                if (this.$data.isCollapse) {
                    this.hiddenMenu();
                } else {
                    this.openMenu();
                }
            },
            tableRowStyle({row, rowIndex}) {
                return 'font-size:0.8rem;'
            },
            tableHeaderColor({row, column, rowIndex, columnIndex}) {
                if (rowIndex === 0) {
                    return 'font-size:0.8rem;background-color: #f0f9eb;color:#67C23A;font-weight: bold;'
                }
            },
            exit() {
                window.localStorage.removeItem('token');
                this.$router.push({path: '/login'});
            }
        }
    }
</script>

<style>
    .body {
        height: 100%;
        width: 100%;
    }

    .el-menu-vertical {
        height: 100%;
    }

    .el-menu-vertical:not(.el-menu--collapse) {
        width: 200px;
        height: 100%;
    }

    .nav-bottom-open {
        position: fixed;
        bottom: 0.5rem;
        left: 0.6rem;
    }

    .nav-bottom-close {
        position: fixed;
        bottom: 0.5rem;
        left: 9rem;
    }

    .header {
        box-shadow: 1px 0 5px #ccc;
        line-height: 46px;
    }
</style>
