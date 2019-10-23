<template>
    <div class="body">
        <el-page-header @back="goBack" content="编辑数据应用">
        </el-page-header>
        <div>
            <el-form ref="form" class="form" labelPosition="right" :model="form" label-width="100px">
                <el-form-item label="应用名称:" :rules="[{ required: true, message: '服务名称不能为空'}]">
                    <el-input v-model="form.name" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="应用类型:" :rules="[{ required: true, message: '连通方式不能为空'}]">
                    <el-select v-model="form.type" placeholder="请选择连通方式">
                        <el-option
                                v-for="item in types"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="建设单位:" :rules="[{ required: true, message: '所属部门不能为空'}]">
                    <el-select v-model="form.section" placeholder="请选择所属部门">
                        <el-option
                                v-for="item in sections"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="数据服务:" :rules="[{ required: true, message: '描述不能为空'}]">
                    <el-transfer
                            :titles="['数据源','已选择']"
                            filterable
                            :filter-method="filterMethod"
                            filter-placeholder="请输入数据源名称"
                            :button-texts="['到左边', '到右边']"
                            v-model="value"
                            :data="data">
                    </el-transfer>
                </el-form-item>
                <el-form-item label="描述:" :rules="[{ required: true, message: '描述不能为空'}]">
                    <el-input type="textarea" v-model="form.description" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="success" @click="onSubmit">立即创建</el-button>
                    <el-button type="danger" plain>取消</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
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

    import {Loading} from 'element-ui';

    const pinyin = require("pinyin");

    export default {
        name: "application_edit",
        data() {
            return {
                data: [],
                value: [],
                types: [],
                sections: [],
                form: {
                    name: '',
                    uri: '',
                    check: '0',
                    section: '',
                    type: '',
                    description: '',
                    api_doc: ''
                }
            }
        },
        methods: {
            onSubmit() {
                let _this = this;
                const loading = Loading.service({
                    lock: true,
                    text: '正在提交数据保存...',
                });
                setTimeout(function () {
                    axios.post(window.uris.server + window.uris.source.add, _this.$data.form).then(function (response) {
                        loading.close();
                    }).catch(function (error) {
                        loading.close();
                        this.$message.error(error);
                    });
                }, 1000);
            },
            goBack() {
                this.$router.push({path: '/1-2'});
            },
            filterMethod(query, item) {
                return item.pinyin.indexOf(query) > -1;
            }
        },
        mounted() {
            let _this = this;
            axios.get(window.uris.server + '/app.json')
                .then(function (res) {
                    _this.$data.types = res.data
                })
                .catch(function (error) {
                    console.log(error);
                });
            axios.get(window.uris.server + '/section.json')
                .then(function (res) {
                    _this.$data.sections = res.data
                })
                .catch(function (error) {
                    console.log(error);
                });
            resource(window.uris.server + window.uris.source.select, 0, 9999999, function (res) {
                let items = res.data.items;
                for (let i = 0; i < items.length; i++) {
                    _this.$data.data.push({
                        label: items[i].name,
                        key: i,
                        pinyin: pinyin(items[i].name, {style: pinyin.STYLE_NORMAL}).join('')
                    });
                }
                _this.$data.total = res.data.total;
            });
        }
    }
</script>

<style scoped>
    .form {
        margin: 2rem 1rem 1rem 0;
    }
</style>