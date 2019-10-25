<template>
    <div class="body">
        <el-alert style="margin: 0 0 1rem 0"
                  title="帮助消息"
                  description="系统目前处于研发阶段，功能正在不断完善中!"
                  show-icon>
        </el-alert>
        <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item :to="{ path: '/' }">数据中心</el-breadcrumb-item>
            <el-breadcrumb-item>数据应用</el-breadcrumb-item>
        </el-breadcrumb>
        <div>
            <div class="toolbar">
                <el-button-group>
                    <el-button type="success" icon="el-icon-circle-plus" size="small" @click="add">新增应用</el-button>
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
                            <el-form-item label="App Id：">
                                <span style="color: #888">
                                    <template><el-tag>{{ props.row.id }}</el-tag></template>
                                </span>
                            </el-form-item>
                            <el-form-item label="Secret Key：">
                                <span style="color: #888">
                                    <template><el-tag type="danger">{{ props.row.secret_key }}</el-tag></template>
                                </span>
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
                        label="应用类型"
                        prop="type">
                    <template slot-scope="props">
                        <el-tag type="green" size="mini">{{ props.row.type }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        label="建设单位">
                    <template slot-scope="props">
                        <el-tag type="success" size="mini" effect="plain">{{props.row.section}}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                        label="数据源"
                        prop="sources">
                    <template slot-scope="props">
                        <el-button size="mini" round @click="openSource(props.row.id)">数据源</el-button>
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
                                <el-dropdown-item @click.native="config(props.row.id)">配置</el-dropdown-item>
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
            <el-drawer
                    title="数据源"
                    :visible.sync="drawer"
                    :direction="direction">
                <el-row>
                    <el-table
                            :data="tt"
                            style="width: 100%;height: 100%;">
                        <el-table-column
                                label="数据名称">
                            <template slot-scope="props">
                                <el-tag>{{props.row.name}}</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column
                                width="100"
                                label="操作">
                            <template slot-scope="props">
                                <el-button type="danger" size="mini" round
                                           @click="removeSource(props.row.id,props.row.appId)">移除
                                </el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-row>
            </el-drawer>

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
        name: "application",
        data() {
            return {
                tableData: [],
                total: 0,
                c: 1,
                drawer: false,
                direction: 'rtl',
                tt: []
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
            config(id){
                let _this = this;
                _this.$router.push({path: '/1-2/edit/' + id})
            },
            openSource(id) {
                let _this = this;
                _this.$data.tt = [];
                axios.get(window.uris.server + window.uris.application.source, {
                    params: {
                        id: id,
                    }
                }).then(function (res) {
                    if (res.data.data == null) {
                        _this.$data.drawer = true;
                        return;
                    }
                    for (let i = 0; i < res.data.data.length; i++) {
                        let dd = res.data.data[i];
                        dd.appId = id;
                        _this.$data.tt.push(dd);
                    }
                    _this.$data.drawer = true
                })
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
            removeSource(val, app) {
                let _this = this;
                this.$confirm('此操作将移除此项数据服务, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'info'
                }).then(() => {
                    axios.get(window.uris.server + window.uris.application.source_delete, {
                        params: {
                            id: val,
                            app: app,
                        }
                    }).then(function (res) {
                        _this.$message({
                            type: 'success',
                            message: '服务成功移除!'
                        });
                        _this.openSource(app);
                    }).catch(function (error) {
                        _this.$message({
                            type: 'error',
                            message: '服务移除失败！'
                        });
                    });
                }).catch(() => {
                    _this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });
                });
            },
            remove(val) {
                let _this = this;
                this.$confirm('此操作将永久删除该项数据应用, 是否继续?', '提示', {
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