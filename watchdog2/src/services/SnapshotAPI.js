import axios from 'axios'

export default {

  getSnapshots (fromTimestamp, cameras, cursor) {
    return axios.get('/snapshots/' + fromTimestamp + '/' + cameras + '/' + cursor)
    .then(response => {
      return response.data;
    })
  },
  countSnapshots (fromTimestamp, cameras) {
    return axios.get('/count-snapshots/' + fromTimestamp + '/' + cameras)
    .then(response => {
      return response.data;
    })
  },
  getSnapshotsLimit () {
    return axios.get('/snapshots-limit')
    .then(response => {
      return response.data;
    })
  },
  getCameras () {
    return axios.get('/cameras')
    .then(response => {
      return response.data;
    })
  },
  getFirtSnapshotDate () {
    return axios.get('/first-snapshot-date')
    .then(response => {
      return response.data;
    })
  },
  getLastSnapshotDate () {
    return axios.get('/last-snapshot-date')
    .then(response => {
      return response.data;
    })
  }
}