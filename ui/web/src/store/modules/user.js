import axios from "axios";

export default {
    namespaced: true,
    state: {
        currentUser: null
    },
    getters: {
        currentUser: state => state.currentUser,
    },
    mutations: {
        setUser(state, payload) {
            state.currentUser = payload
        },
        setLogout(state) {
            state.currentUser = null
        },
    },
    actions: {
        loadCurrent: async ({commit}) => {
            try {
                let url = window.location.origin + "/api/users/current";
                let response = await axios.get(url);
                commit("setUser", response.data)
            } catch (e) {
                if (e.response.status !== 401) {
                    console.error(e)
                }
            }
        },
        signOut: async ({commit}) => {
            try {
                let response = await axios.get(window.location.origin + "/auth/sign-out");
            } catch (e) {
                console.error(e)
            }
            commit("setLogout")
        }
    }
}
