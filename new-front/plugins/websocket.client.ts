import {defineNuxtPlugin} from "nuxt/app";

export default defineNuxtPlugin(() => {
    let socket = new WebSocket(`ws://localhost:2137/api/ws/get-streams`);
        
    return {
        provide: {
            socket,
        },
    }
})