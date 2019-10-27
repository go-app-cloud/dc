<template>
    <div class="body">
        <el-page-header @back="goBack" content="编辑数据源">
        </el-page-header>
        <div>
            <el-form ref="form" class="form" labelPosition="right" :model="form" label-width="130px">
                <el-form-item label="服务名称:" :rules="[{ required: true, message: '服务名称不能为空'}]">
                    <el-input v-model="form.name" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="描述:" :rules="[{ required: true, message: '描述不能为空'}]">
                    <el-input type="textarea" v-model="form.description" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="连通方式:" :rules="[{ required: true, message: '连通方式不能为空'}]">
                    <el-select v-model="form.type" placeholder="请选择连通方式">
                        <el-option
                                v-for="item in types"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="所属部门:" :rules="[{ required: true, message: '所属部门不能为空'}]">
                    <el-select v-model="form.section" placeholder="请选择所属部门">
                        <el-option
                                v-for="item in sections"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="上报接口:" :rules="[{ required: true, message: '数据上报接口不能为空'}]">
                    <el-input v-model="form.service"></el-input>
                </el-form-item>
                <el-form-item label="查询接口:" :rules="[{ required: true, message: '数据查询接口不能为空'}]">
                    <el-input v-model="form.search"></el-input>
                </el-form-item>
                <el-form-item label="文档:" :rules="[{ required: true, message: '文档不能为空'}]">
                    <el-input v-model="form.api_doc"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="success" @click="onSubmit">{{button_text}}</el-button>
                    <el-button type="danger" plain>取消</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
    let submitType = 0;
    import {Loading} from 'element-ui';

    export default {
        name: "source_edit",
        data() {
            return {
                types: [],
                sections: [],
                button_text: '立即创建',
                form: {
                    name: '',
                    service: '',
                    section: '',
                    type: '',
                    description: '',
                    api_doc: '',
                    search: ''
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
                if (submitType === 0) {
                    setTimeout(function () {
                        axios.post(window.uris.server + window.uris.source.add, _this.$data.form).then(function (response) {
                            loading.close();
                        }).catch(function (error) {
                            loading.close();
                            this.$message.error(error);
                        });
                    }, 1000);
                } else {
                    setTimeout(function () {
                        axios.post(window.uris.server + window.uris.source.modify, _this.$data.form).then(function (response) {
                            loading.close();
                        }).catch(function (error) {
                            loading.close();
                            this.$message.error(error);
                        });
                    }, 1000);
                }

            },
            goBack() {
                this.$router.push({path: '/1-1'});
            }
        },
        mounted() {
            let _this = this;
            if (typeof this.$route.params.id != 'undefined') {
                submitType = 1;
                _this.$data.button_text = '更新数据';
                axios.post(window.uris.server + window.uris.source.get, {
                    id: this.$route.params.id,
                }).then(function (response) {
                    console.log(response.data);
                    _this.$data.form = response.data.data
                }).catch(function (error) {
                });
            }
            axios.get(window.uris.server + '/type.json')
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
        }
    }
</script>

<style scoped>
    .form {
        margin: 2rem 1rem 1rem 0;
    }
</style>