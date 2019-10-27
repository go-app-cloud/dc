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
            <el-breadcrumb-item>数据来源</el-breadcrumb-item>
        </el-breadcrumb>
        <div>
            <div class="toolbar">
                <el-button-group>
                    <el-button type="success" icon="el-icon-circle-plus" size="small" @click="add">新增源</el-button>
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
                            <el-form-item label="上报接口：">
                                <el-tag type="warning" size="small">{{ props.row.service }}</el-tag>
                            </el-form-item>
                            <el-form-item label="查询接口：">
                                <el-tag type="warning" size="small">{{ props.row.search }}</el-tag>
                            </el-form-item>
                            <el-form-item label="App Id：">
                                <el-tag>{{ props.row.id}}</el-tag>
                            </el-form-item>
                            <el-form-item label="Secret Key：">
                                <el-tag>{{ props.row.secret}}</el-tag>
                            </el-form-item>
                            <el-form-item label="当前状态：">
                                <span style="color: #888">
                                    <template v-if="props.row.check === 0"><el-tag
                                            type="danger" effect="dark" size="mini">服务离线</el-tag></template>
                                    <template v-else><el-tag type="success" size="mini" effect="dark">服务在线</el-tag></template>
                                </span>
                            </el-form-item>
                        </el-form>
                    </template>
                </el-table-column>
                <el-table-column
                        width="200"
                        label="ID(名称)"
                        prop="name">
                </el-table-column>
                <el-table-column
                        label="连通方式"
                        prop="type">
                    <template slot-scope="props">
                        <el-tag type="green" size="mini">{{ props.row.type }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        label="数据提供"
                        prop="section">
                    <template slot-scope="props">
                        <el-tag type="danger" size="mini">{{ props.row.section }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        width="120"
                        label="所属者">
                    <template slot-scope="props" v-if="props.row.owner === ''">
                        <el-tag type="success" size="mini" effect="plain">系统应用</el-tag>
                    </template>
                    <template v-else>
                        <el-tag type="primary" size="mini">{{props.row.owner}}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        width="100"
                        label="工具">
                    <el-button type="danger" size="mini" round>测试工具</el-button>
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
                                <el-dropdown-item @click.native="config(props.row.id)">配置</el-dropdown-item>
                                <el-dropdown-item @click.native="open(props.row.api_doc)">文档</el-dropdown-item>
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
        name: "sources",
        data() {
            return {
                tableData: [],
                total: 0,
                c: 1
            }
        },
        mounted() {
            let _this = this;
            resource(window.uris.server + window.uris.source.select, _this.$data.c, page, function (res) {
                console.log(res.data.items);
                _this.$data.tableData = res.data.items;
                _this.$data.total = res.data.total;
            });
        },
        methods: {
            config(id) {
                let _this = this;
                _this.$router.push({path: '/1-1/edit/' + id})
            },
            tableRowStyle({row, rowIndex}) {
                return 'font-size:0.8rem;'
            },
            tableHeaderColor({row, column, rowIndex, columnIndex}) {
                if (rowIndex === 0) {
                    return 'font-size:0.8rem;background-color: #f0f9eb;color:#67C23A;font-weight: bold;'
                }
            },
            add() {
                this.$router.push({path: '/1-1/edit'})
            },
            fresh() {
                let _this = this;
                resource(window.uris.server + window.uris.source.select, _this.$data.c, page, function (res) {
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
                    axios.get(window.uris.server + window.uris.source.remove, {
                        params: {
                            id: val,
                        }
                    }).then(function (res) {
                        _this.$message({
                            type: 'success',
                            message: '删除成功!'
                        });
                        resource(window.uris.server + window.uris.source.select, 1, page, function (res) {
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
                resource(window.uris.server + window.uris.source.select, val, page, function (res) {
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