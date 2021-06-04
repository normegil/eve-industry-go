import Vue from "vue";
import VueRouter from "vue-router";
import loadAuthGuard from "./utils/auth.guard";
import {adminRoot} from "./constants/config";
import {UserRole} from "./utils/auth.roles";

Vue.use(VueRouter);

export default function loadRoutes(store) {
    const routes = [
        {
            path: adminRoot,
            component: () => import(/* webpackChunkName: "app" */ "./views/app"),
            redirect: `${adminRoot}/manufacturing`,
            meta: {loginRequired: true},
            children: [
                {
                    path: "manufacturing",
                    component: () =>
                        import(/* webpackChunkName: "manufacturing" */ "./views/app/manufacturing"),
                    redirect: `${adminRoot}/manufacturing/myblueprints`,
                    children: [
                        {
                            path: 'myblueprints',
                            component: () => import(/* webpackChunkName: "manufacturing" */ './views/app/manufacturing/MyBlueprints.vue'),
                            meta: {roles: [UserRole.User]},
                        }
                    ]
                },
                {
                    path: "research",
                    component: () =>
                        import(/* webpackChunkName: "research" */ "./views/app/research"),
                    redirect: `${adminRoot}/research/invention`,
                    children: [
                        {
                            path: 'invention',
                            component: () => import(/* webpackChunkName: "manufacturing" */ './views/app/research/Invention'),
                            meta: {roles: [UserRole.User]}
                        }
                    ]
                },
            ]
        },
        {
            path: "/error",
            component: () => import(/* webpackChunkName: "error" */ "./views/Error")
        },
        {
            path: "/user",
            component: () => import(/* webpackChunkName: "user" */ "./views/user"),
            redirect: "/user/login",
            children: [
                {
                    path: "login",
                    component: () =>
                        import(/* webpackChunkName: "user" */ "./views/user/Login")
                },
            ]
        },
        {
            path: "*",
            component: () => import(/* webpackChunkName: "error" */ "./views/Error")
        }
    ];

    const router = new VueRouter({
        linkActiveClass: "active",
        routes,
        mode: "history",
    });
    let authGuard = loadAuthGuard(store)
    router.beforeEach(authGuard);
    return router;
}