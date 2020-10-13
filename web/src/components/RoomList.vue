<template>
  <div>

    <v-card
    v-for="room in rooms"
    :key="room.room_id"
    class="mx-auto"
    max-width="344"
    >
      <v-card-title>
        {{room.creator}}
      </v-card-title>
      <v-card-subtitle>
        {{ room.room_id}}
      </v-card-subtitle>

      <v-card-text>
        虚拟人数
        {{ room.connected + room.wait_connect}}
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
export default {
name: "RoomList",
  data: function (){
  return {
    rooms: []
  }
  },
  mounted() {
  this.$axios({
    url: "/api/room/status",
    method: "get",
  }).then(res => {
    console.log(res)
    this.rooms = res.data
    // this.rooms = JSON.parse(res.data)
  }).catch(res => {
    console.log(res)
  })

  }
}
</script>

<style scoped>

</style>