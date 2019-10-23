<template>
    <div style="padding: 10rem 0 0 0;background: #ccc;height: 100%;width: 100%;box-sizing: border-box;">
        <el-row>
            <el-col :span="8">
                <div class="grid-content "></div>
            </el-col>
            <el-col :span="8">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span><i class="el-icon-eleme"></i> 数据融合平台</span>
                        <el-tag style="float: right;" size="small">极光云</el-tag>
                    </div>
                    <div class="text item">
                        <el-form label-position="right" label-width="80px" :model="form">
                            <el-form-item label="用户名">
                                <el-input v-model="form.username"></el-input>
                            </el-form-item>
                            <el-form-item label="密码">
                                <el-input v-model="form.pass" type="password" autocomplete="off"></el-input>
                            </el-form-item>
                            <el-divider content-position="left">极光云</el-divider>
                            <el-form-item>
                                <el-button type="primary" @click="onSubmit">立即登录</el-button>
                                <el-button>取消</el-button>
                            </el-form-item>
                        </el-form>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <div class="grid-content"></div>
            </el-col>
        </el-row>
    </div>
</template>

<script>
    import {Loading} from 'element-ui';

    export default {
        name: "login",
        data() {
            return {
                form: {
                    username: '',
                    pass: ''
                }
            };
        },
        methods: {
            onSubmit() {
                let _this = this;
                const loading = Loading.service({
                    lock: true,
                    text: '正在提交数据...',
                });
                setTimeout(function () {
                    axios.post(window.uris.server + window.uris.login, _this.$data.form).then(function (response) {
                        loading.close();
                        if (response.data.code === 0) {
                            window.localStorage.setItem('token', response.data.data.token);
                            _this.$router.push({path: '/'})
                        } else {
                            _this.$notify.error({
                                title: '认证失败',
                                message: '登录认证失败，请检测用户名和密码.'
                            });
                        }
                    }).catch(function (error) {
                        loading.close();
                        this.$message.error(error);
                    });
                }, 1000);
            }
        }
    }
</script>

<style scoped>

    .grid-content {
        border-radius: 4px;
        min-height: 36px;
    }
</style>