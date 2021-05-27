import {isAuthGuardActive} from '../constants/config'

export default function loadAuthGuard(store) {
    return (to, from, next) => {
        if (to.matched.some(record => record.meta.loginRequired)) {
            if (isAuthGuardActive) {
                const user = store.getters["user/currentUser"];
                if (user) {
                    const roleArrayHierarchic = to.matched.filter(x => x.meta.roles).map(x => x.meta.roles);
                    if (roleArrayHierarchic.every(x => x.includes(user.role))) {
                        next();
                    } else {
                        next('/unauthorized')
                    }
                } else {
                    store.commit("user/setUser", null)
                    next('/user/login')
                }
            } else {
                next();
            }
        } else {
            next()
        }
    }
}
