<template>
    <h3 class="mb-4 ml-4 text-sm font-medium text-bodydark2">Streams</h3>
    <ul class="mb-6 flex flex-col gap-1.5">
        <li>
            <div class="streams" id="streams">
                <div v-for="(m, idx) in messages">
                    <div
                        class="rounded-sm border border-stroke bg-white py-6 px-7.5 shadow-default dark:border-strokedark dark:bg-boxdark">
                        <div class="mt-0 flex items-end justify-between">
                            <div>
                                <h5 class="text-title-sm font-bold text-black dark:text-white">
                                    {{ $route.params.name }}
                                    {{ m }}
                                </h5>
                                <span class="text-sm font-medium">{{ idx }}</span>
                            </div>

                            <span class="flex items-center gap-1 text-sm font-medium text-meta-3">
                                URL
                                <svg class="fill-meta-3" width="10" height="11" viewBox="0 0 10 11" fill="none"
                                    xmlns="http://www.w3.org/2000/svg">
                                    <path
                                        d="M4.35716 2.47737L0.908974 5.82987L5.0443e-07 4.94612L5 0.0848689L10 4.94612L9.09103 5.82987L5.64284 2.47737L5.64284 10.0849L4.35716 10.0849L4.35716 2.47737Z"
                                        fill="" />
                                </svg>
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </li>
    </ul>
</template>

<script>
import axios from 'axios'

export default {
    data: function () {
        return {
            messages: [],
            connection: null,
        }
    },

    mounted: function () {

    },

    created: function () {
        this.init()
    },

    // deactivated: function () {
    //     this.connection.close()
    // },

    destroyed: function () {

    },

    methods: {
        init() {
            console.log("Starting websocket connection to backend");
            this.connection = new WebSocket("ws://localhost:2137/api/ws");

            let actions = { "room": "007854ce7b93476487c7ca8826d17eba", "info": "1121212" };

            // this.connection.onopen = () => this.connection.send(JSON.stringify(actions));
            this.connection.onmessage = this.onSocketMessage
        },
        onSockOpen() {
            console.log("connected")
        },
        onSocketMessage(evt) {
            console.log(evt)
            this.messages.push(evt.data)
        },
        getStreams: function () {
            axios.get(
                "http://localhost:2137" + "/api/streams-by-service",
                {
                    "name": "fig",
                    "port": 1111,
                    "offset": 0,
                    "limit": 20,
                }
            ).then(result => {
                this.streams = result.data
                console.log(result, result.data)
            }).catch(error => {
                console.error(error)
            })
        }
    },
}

</script>