<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-sm">
                <ImageItem
                    v-if="snapshots != null && snapshots.length > 1"
                    :source="snapshots[1].photosmallPath"
                />
            </div>
            <div class="col-sm">
                <img alt="Vue logo" src="../assets/logo.png"/>
            </div>
            <div class="col-sm">
                <img alt="Vue logo" src="../assets/logo.png"/>
            </div>
        </div>
    </div>
</template>

<script>
import snapshotAPI from '@/services/SnapshotAPI'
import ImageItem from "./ImageItem";
import axios from 'axios'

export default {
    name: "Gallery",
    data () {
        return {
        snapshots: []
        }
    },
    async mounted () {
        try {
            snapshotAPI.fetchSnapshotAllCam('2019-11-01T05:41:00', 0)
            .then(response => {
                this.snapshots = response
                
                this.snapshots.forEach(function(snapshot) {
                    snapshot.photosmallPath = axios.defaults.baseURL + snapshot.photosmallPath
                });
                console.log(this.snapshots)
            })
            .catch(error => {
                console.log(error)
            })
        } catch (error) {
            throw(error)
        }
    },
    methods: {
    },
    components: {
        ImageItem
    }
}
</script>