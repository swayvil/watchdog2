import axios from 'axios'

export default {

  fetchSnapshotAllCam (fromTimestamp, cursor) {
    return axios.get('/snapshots-all-cams/' + fromTimestamp + '/' + cursor)
    .then(response => {
      return response.data;
    })
  },
  countSnapshotAllCam (fromTimestamp) {
    return axios.get('/count-snapshots-all-cams/' + fromTimestamp)
    .then(response => {
      return response.data;
    })
  },
  getSnapshotLimit () {
    return axios.get('/snapshots-limit')
    .then(response => {
      return response.data;
    })
  }
}