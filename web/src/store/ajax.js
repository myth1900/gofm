import axios from 'axios'
export function GetRoomsStatus() {
    axios.get("/api/rooms/status")
        .then(res =>{
            return res.data
        })
}

export function UpdateRoomAudiences(roomID,nums) {
    let url = "/api/room/" + roomID + "/audience"
    axios.put(url, {},{
        params: {nums:nums}
    }).then(res => {
        console.log(res.data)
    })
}