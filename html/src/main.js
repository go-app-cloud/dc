import Vue from 'vue'
import App from './App.vue'
import Main from './Main.vue'
import VueRouter from 'vue-router'
import IndexPage from './pages/index'
import LoginPage from './pages/login'
import SourcePage from './pages/source'
import SourceEditPage from './pages/source_edit'
import ApplicationPage from './pages/application'
import ApplicationEditorPage from './pages/application_edit'
import DockerPage from './pages/docker'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

window.axios = require('axios');
const qs = require('qs');
window.axios.interceptors.request.use(
    config => {
        if (config.method === 'post') {
            config.data = qs.stringify(config.data)
        }
        return config
    },
    error => {
        console.log(error);
        Promise.reject(error)
    }
);

window.uris = {
    server: 'http://localhost:8000'
};
window.uris.login = '/login.cgi';
window.uris.source = {
    add: '/source/add.cgi',
    remove: '/source/delete.cgi',
    modify: '/source/modify.cgi',
    select: '/source/select.cgi',
    get: '/source/get.cgi'
};
window.uris.application = {
    add: '/application/add.cgi',
    remove: '/application/delete.cgi',
    modify: '/application/modify.cgi',
    select: '/application/select.cgi',
    source: '/application/source.cgi',
    get: '/application/get.cgi',
    source_delete: '/application/source_delete.cgi'
};
Vue.use(VueRouter).use(ElementUI);

const router = new VueRouter({
    routes: [
        {path: '/login', component: LoginPage},
        {
            path: '/', component: Main, children: [
                {path: '/', component: IndexPage},
                {path: '/1-1', component: SourcePage},
                {path: '/1-1/edit', component: SourceEditPage},
                {path: '/1-1/edit/:id', component: SourceEditPage},
                {path: '/1-2', component: ApplicationPage},
                {path: '/1-2/edit', component: ApplicationEditorPage},
                {path: '/1-2/edit/:id', component: ApplicationEditorPage},
                {path: '/1-3', component: DockerPage}
            ]
        },
    ],
});

Vue.config.productionTip = false;

new Vue({
    router,
    render: h => h(App),
}).$mount('#app');