<template>
    <div class="body">
        <el-alert style="margin: 0 0 1rem 0"
                  title="帮助消息"
                  type="warning"
                  description="系统目前处于研发阶段，功能正在不断完善中!"
                  show-icon>
        </el-alert>
        <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item :to="{ path: '/' }">数据中心</el-breadcrumb-item>
            <el-breadcrumb-item>容器服务</el-breadcrumb-item>
        </el-breadcrumb>
        <div>
            <div class="toolbar">
                <el-button-group>
                    <el-button type="success" icon="el-icon-circle-plus" size="small" @click="add">Dockerfile</el-button>
                    <el-button type="success" icon="el-icon-circle-plus" size="small" @click="add">容器应用</el-button>
                    <el-button type="success" icon="el-icon-refresh" size="small" @click="fresh">刷新</el-button>
                    <el-button type="success" icon="el-icon-question" size="small" @click="help">帮助</el-button>
                </el-button-group>
            </div>
            <el-table
                    :data="tableData"
                    size="small"
                    :row-style="tableRowStyle"
                    :header-cell-style="tableHeaderColor"
                    style="width: 100%;margin:16px 0 0 0;">
                <el-table-column type="expand">
                    <template slot-scope="props">
                        <el-form label-position="left" inline class="demo-table-expand">
                            <el-form-item label="描述信息：">
                                <span style="color: #888">{{ props.row.description }}</span>
                            </el-form-item>
                        </el-form>
                    </template>
                </el-table-column>
                <el-table-column
                        width="200"
                        label="ID(应用名称)"
                        prop="name">
                </el-table-column>
                <el-table-column
                        label="运行状态"
                        prop="type">
                    <template slot-scope="props">
                        <el-tag type="green" size="mini">{{ props.row.type }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        label="容器镜像"
                        prop="sources">
                    <template slot-scope="props">
                        <el-tag type="danger" size="mini">{{ props.row.sources.split('|').length }} 套数据源</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        width="120"
                        label="映射">
                    <template slot-scope="props">
                        <el-tag type="success" size="mini" effect="plain">{{props.row.section}}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        label="操作"
                        prop="id">
                    <template slot-scope="props">
                        <el-dropdown>
                            <el-button type="success" size="mini" round>
                                操作<i class="el-icon-arrow-down el-icon--right"></i>
                            </el-button>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item>配置</el-dropdown-item>
                                <el-dropdown-item>授权</el-dropdown-item>
                                <el-dropdown-item @click.native="remove(props.row.id)" divided>删除
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination
                    @current-change="current"
                    style="margin: 0.6rem 0 0 0"
                    background
                    layout="prev, pager, next"
                    :total="total">
            </el-pagination>
        </div>
    </div>
</template>

<script>
    const page = 10;

    function resource(uri, index, page, callback) {
        axios.get(uri, {
            params: {
                page: page,
                index: index,
            }
        }).then(function (res) {
            callback(res.data)
        })
    }

    export default {
        name: "docker",
        data() {
            return {
                tableData: [],
                total: 0,
                c: 1
            }
        },
        mounted() {
            let _this = this;
            resource(window.uris.server + window.uris.application.select, _this.$data.c, page, function (res) {
                _this.$data.tableData = res.data.items;
                _this.$data.total = res.data.total;
            });
        },
        methods: {
            tableRowStyle({row, rowIndex}) {
                return 'font-size:0.8rem;'
            },
            tableHeaderColor({row, column, rowIndex, columnIndex}) {
                if (rowIndex === 0) {
                    return 'font-size:0.8rem;background-color: #f0f9eb;color:#67C23A;font-weight: bold;'
                }
            },
            add() {
                this.$router.push({path: '/1-2/edit'})
            },
            fresh() {
                let _this = this;
                resource(window.uris.server + window.uris.application.select, _this.$data.c, page, function (res) {
                    _this.$data.tableData = res.data.items;
                    _this.$data.total = res.data.total;
                });
            },
            help() {

            },
            open(uri) {
                window.open(uri, '_blank')
            },
            remove(val) {
                let _this = this;
                this.$confirm('此操作将永久删除该项数据服务, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    axios.get(window.uris.server + window.uris.application.remove, {
                        params: {
                            id: val,
                        }
                    }).then(function (res) {
                        _this.$message({
                            type: 'success',
                            message: '删除成功!'
                        });
                        resource(window.uris.server + window.uris.application.select, 1, page, function (res) {
                            _this.$data.tableData = res.data.items;
                            _this.$data.total = res.data.total;
                        });
                    }).catch(function (error) {
                        _this.$message({
                            type: 'error',
                            message: '删除失败！'
                        });
                    });
                }).catch(() => {
                    _this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });
                });
            },
            current(val) {
                let _this = this;
                _this.$data.c = val;
                resource(window.uris.server + window.uris.application.select, val, page, function (res) {
                    _this.$data.tableData = res.data.items;
                    _this.$data.total = res.data.total;
                });
            }
        }
    }
</script>

<style scoped>
    .demo-table-expand {
        font-size: 0;
    }

    .demo-table-expand label {
        width: 90px;
        color: #99a9bf;
    }

    .demo-table-expand .el-form-item {
        margin-right: 0;
        margin-bottom: 0;
        width: 100%;
    }

    .toolbar {
        margin: 1rem 0 0 0;
    }
</style>