import axios from 'axios'

export default {

  fetchSnapshotAllCam (fromTimestamp, offset) {
    console.log(fromTimestamp)
    return axios.get('/snapshots-all-cams/' + fromTimestamp + '/' + offset)
    .then(response => {
      return response.data;
    })
  }
}