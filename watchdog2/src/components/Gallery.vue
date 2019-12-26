<template>
  <div class="container-fluid" :key="componentKey">
    <div v-for="(n, i) in snapshots.length / columnSize" :key="n" class="row">
      <div v-for="(m, j) in columnSize" :key="m" class="col-sm">
        <ImageItem
          v-if="snapshots != null && snapshots.length > 1"
          :source="snapshots[i * columnSize + j].photosmallPath"
        />
      </div>
    </div>
    <div class="row">
      <div class="col-sm"></div>
      <div class="col-sm">
        <button v-if="cursor > 0" v-on:click="loadSnapshots(-1)" type="button" class="btn btn-primary px-3">&lt;</button>
        <button v-if="cursor * maxSnapshots + snapshots.length < countSnapshots" v-on:click="loadSnapshots(1)" type="button" class="btn btn-primary px-3">&gt;</button>
      </div>
      <div class="col-sm"></div>
    </div>
  </div>
</template>

<script>
import snapshotAPI from "@/services/SnapshotAPI";
import ImageItem from "./ImageItem";
import axios from "axios";

export default {
  name: "Gallery",
  data() {
    return {
      componentKey: 0,
      maxSnapshots: 0,
      cursor: 0,
      countSnapshots: 0,
      snapshots: [],
      columnSize: 5
    };
  },
  async mounted() {
    try {
      snapshotAPI
        .countSnapshotAllCam("2019-11-01T05:41:00")
        .then(response => {
          this.countSnapshots = response;
        })
        .catch(error => {
          console.log(error);
        });

      snapshotAPI
        .getSnapshotLimit()
        .then(response => {
          this.maxSnapshots = response;
        })
        .catch(error => {
          console.log(error);
        });

        this.loadSnapshots(0);
    } catch (error) {
      throw error;
    }
  },
  methods: {
    loadSnapshots: function(cursor) {
      this.cursor += cursor;
      snapshotAPI
        .fetchSnapshotAllCam("2019-11-01T05:41:00", this.cursor)
        .then(response => {
          response.forEach(function(snapshot) {
            snapshot.photosmallPath =
              axios.defaults.baseURL + snapshot.photosmallPath;
          });
          this.snapshots = response;
          this.componentKey = (this.componentKey + 1) % 2; // To force the reload of the component
        })
        .catch(error => {
          console.log(error);
        });
    }
  },
  components: {
    ImageItem
  }
};
</script>