import Vue from 'vue'
import Vuex from 'vuex'
import {REFRESH_ROOMS_STATUS} from "@/store/mutation-types";
import {GetRoomsStatus} from "@/store/ajax";
Vue.use(Vuex)
export default new Vuex.Store({
    state:{
        rooms: []
    },
    getters:{
        getRooms: state => {
            return state.rooms
        }
    },
    mutations: {
        [REFRESH_ROOMS_STATUS](state){
            state.rooms = GetRoomsStatus()
        }
    }
})