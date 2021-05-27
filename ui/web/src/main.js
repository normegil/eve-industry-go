import Vue from 'vue'
import App from './App'
// BootstrapVue add
import BootstrapVue from 'bootstrap-vue'
// Router & Store add
import loadRoutes from './router'
import store from './store'
// Multi Language Add
import en from './locales/en.json'
import VueI18n from 'vue-i18n'
// Notification Component Add
import Notifications from './components/Common/Notification'
// Breadcrumb Component Add
import Breadcrumb from './components/Common/Breadcrumb'
// RefreshButton Component Add
import RefreshButton from './components/Common/RefreshButton'
// Colxx Component Add
import Colxx from './components/Common/Colxx'
// Perfect Scrollbar Add
import vuePerfectScrollbar from 'vue-perfect-scrollbar'
import contentmenu from 'v-contextmenu'
import VueLineClamp from 'vue-line-clamp'
import VueScrollTo from 'vue-scrollto'
import 'firebase/auth'
import {getCurrentLanguage} from './utils'

Vue.use(BootstrapVue);
Vue.use(VueI18n);
const messages = {en: en};
const locale = getCurrentLanguage();
const i18n = new VueI18n({
    locale: locale,
    fallbackLocale: 'en',
    messages
});
Vue.use(Notifications);
Vue.use(require('vue-shortkey'));
Vue.use(contentmenu);
Vue.use(VueScrollTo);
Vue.use(VueLineClamp, {
    importCss: true
});

Vue.component('piaf-breadcrumb', Breadcrumb);
Vue.component('b-refresh-button', RefreshButton);
Vue.component('b-colxx', Colxx);
Vue.component('vue-perfect-scrollbar', vuePerfectScrollbar);

let router = loadRoutes(store);
Vue.config.productionTip = false

store.dispatch("app/init")
    .then(() => {
        new Vue({
            i18n,
            router,
            store,
            render: h => h(App)
        }).$mount('#app')
    })
    .catch((e) => {
        console.error(e)
    })