export default {
    namespaced: true,
    actions: {
        init: async (context) => {
            await context.dispatch("user/loadCurrent", {}, {root: true})
        }
    }
}