import axios from "axios";
import {apiUrl} from "../../constants/config";
import {toRole} from "../../utils/auth.roles";

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
                let url = apiUrl() + "/api/users/current";
                let response = await axios.get(url, {withCredentials: true});
                response.data.role = toRole(response.data.role)
                commit("setUser", response.data)
            } catch (e) {
                if (e.response.status !== 401) {
                    console.error(e)
                }
            }
        },
        signOut: async ({commit}) => {
            try {
                let response = await axios.get(apiUrl() + "/auth/sign-out", {withCredentials: true});
            } catch (e) {
                console.error(e)
            }
            commit("setLogout")
        }
    }
}
