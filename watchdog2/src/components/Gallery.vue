<template>
  <div class="container-fluid" :key="galleryKey">
    <div class="row">
      <div class="col-sm">
        <h3>Date from:</h3>
        <date-pick
          v-model="dateFrom"
          :isDateDisabled="isDisabledDate"
        ></date-pick>
      </div>
      <div class="col-sm">
        <h3>Cameras:</h3>
        <div class="custom-control custom-checkbox text-left">
          <div v-for="camera in cameras" v-bind:key="camera.id">
            <input
              type="checkbox"
              :id="camera.id"
              :value="camera.name"
              v-model="selectedCameras"
              class="custom-control-input"
            />
            <label :for="camera.id" class="custom-control-label">{{
              camera.name
            }}</label>
          </div>
        </div>
      </div>
      <div class="col-sm"></div>
    </div>
    <div
      v-for="(n, i) in Math.ceil(snapshots.length / columnSize)"
      :key="n"
      class="row"
    >
      <div v-for="(m, j) in columnSize" :key="m" class="col-sm">
        <a
          v-if="i * columnSize + j < snapshots.length"
          v-bind:href="snapshots[i * columnSize + j].photoPath"
          target="_blank"
        >
          <ImageItem :source="snapshots[i * columnSize + j].photosmallPath" />
        </a>
        <p v-if="i * columnSize + j < snapshots.length">
          {{ formatTimestamp(snapshots[i * columnSize + j].timestamp) }}
        </p>
      </div>
    </div>
    <div class="row">
      <div class="col-sm"></div>
      <div class="col-sm">
        <button
          v-if="cursor > 0"
          v-on:click="loadSnapshots(-1)"
          type="button"
          class="btn btn-primary px-3"
        >
          &lt;
        </button>
        <button
          v-if="cursor * maxSnapshots + snapshots.length < countSnapshots"
          v-on:click="loadSnapshots(1)"
          type="button"
          class="btn btn-primary px-3"
        >
          &gt;
        </button>
      </div>
      <div class="col-sm"></div>
    </div>
  </div>
</template>

<script>
import snapshotAPI from "@/services/SnapshotAPI";
import ImageItem from "./ImageItem";
import axios from "axios";
import DatePick from "vue-date-pick";
import "vue-date-pick/dist/vueDatePick.css";

export default {
  name: "Gallery",
  data() {
    return {
      dateFrom: this.today(),
      galleryKey: 0,
      maxSnapshots: 0,
      cursor: 0,
      countSnapshots: 0,
      snapshots: [],
      columnSize: 5,
      cameras: [],
      selectedCameras: [],
      firstSnapshotDate: Date.now(), // Timestamp
      lastSnapshotDate: Date.now() // Timestamp
    };
  },
  async mounted() {
    // Get max number of snapshots to display per page
    snapshotAPI
      .getSnapshotsLimit()
      .then((response) => {
        this.maxSnapshots = response;
      })
      .catch((error) => {
        console.log(error);
      });

    // Get cameras list
    snapshotAPI
      .getCameras()
      .then((response) => {
        response.forEach(
          function(item, index) {
            var camera = { name: item, id: index };
            this.cameras.push(camera);
            this.selectedCameras.push(item);
          }.bind(this)
        );
        this.loadSnapshots(0);
      })
      .catch((error) => {
        console.log(error);
      });

    // Set calendar starting date
    snapshotAPI
      .getFirtSnapshotDate()
      .then((response) => {
        this.firstSnapshotDate = new Date(response).getTime();
      })
      .catch((error) => {
        console.log(error);
      });
  },
  methods: {
    // Set calendar default selected date
    refreshLastSnapshotDate: function() {
      snapshotAPI
        .getLastSnapshotDate()
        .then((response) => {
          console.log(response);
          this.lastSnapshotDate = new Date(response).getTime();
        })
        .catch((error) => {
          console.log(error);
        })
    },
    loadSnapshots: function(cursor) {
      this.refreshLastSnapshotDate(); // Refresh calendar ending date

      // Get total number of snapshots
      snapshotAPI
        .countSnapshots(this.dateFrom, this.formatSelectedCameras())
        .then((response) => {
          this.countSnapshots = response;
        })
        .catch((error) => {
          console.log(error);
        });

      // Get the snapshots
      this.cursor += cursor;
      snapshotAPI
        .getSnapshots(this.dateFrom, this.formatSelectedCameras(), this.cursor)
        .then((response) => {
          response.forEach(function(snapshot) {
            snapshot.photoPath = axios.defaults.baseURL + snapshot.photoPath;
            snapshot.photosmallPath =
              axios.defaults.baseURL + snapshot.photosmallPath;
          });
          this.snapshots = response;
          this.galleryKey = (this.galleryKey + 1) % 2; // To force the reload of the component
          window.scrollTo(0, 0); // Scroll to top
        })
        .catch((error) => {
          console.log(error);
        });
    },
    // 2020-04-23T04:10:01 => 23-04-2020 04:10:01
    formatTimestamp: function(ts) {
      console.log(ts.length)
      if (ts.length == 19) {
        var year = ts.substr(0, 4)
        var month = ts.substr(5, 2)
        var day = ts.substr(8, 2)
        var time = ts.substr(11, 8)
        return day + "-" + month + "-" + year + " " + time
      }
      return ts
    },
    isDisabledDate(date) {
      return date.getTime() < this.firstSnapshotDate ||
              date.getTime() > this.lastSnapshotDate;
    },
    today() {
      var d = new Date(),
        month = "" + (d.getMonth() + 1),
        day = "" + d.getDate(),
        year = d.getFullYear();

      if (month.length < 2) month = "0" + month;
      if (day.length < 2) day = "0" + day;

      return [year, month, day].join("-");
    },
    formatSelectedCameras() {
      var camerasStr = "";
      for (var i = 0; i < this.selectedCameras.length; i++) {
        camerasStr +=
          this.selectedCameras[i].trim() +
          (i < this.selectedCameras.length - 1 ? "," : "");
      }
      return camerasStr;
    },
  },
  watch: {
    dateFrom: function() {
      this.cursor = 0;
      this.loadSnapshots(0);
    },
    selectedCameras: function() {
      if (this.selectedCameras.length == 0) {
        this.selectedCameras.push(this.cameras[0].name);
      }
      this.cursor = 0;
      this.loadSnapshots(0);
    },
  },
  components: {
    ImageItem,
    DatePick,
  },
};
</script>
