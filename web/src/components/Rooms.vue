<template>
  <v-form>
    <v-container>
      <v-row>
        <v-col
        cols="12"
        md="5"
        >
          <v-text-field
              v-model="room_id"
              label="房间ID （在浏览器打开直播间地址栏后的那串数字）"
              required
          >
          </v-text-field>
        </v-col>

        <v-col
            cols="12"
            md="5"
        >
          <v-text-field
              v-model="nums"
              label="人数（0-100，其余无效）"
              required
          >
          </v-text-field>
        </v-col>


        <v-col
            cols="12"
            md="2"
        >
          <v-btn class="mx-2" fab dark color="indigo"
                 v-on:click="submitUpdate"
          >
            <v-icon dark>mdi-check-bold</v-icon>
          </v-btn>
        </v-col>


      </v-row>
      <v-row>
        <v-col
            v-for="room in rooms"
            :key="room.room_id"
            cols="12"
            md="3"
        >
          <a :href="'https://fm.missevan.com/live/' + room.room_id"
          target="_blank">
            <v-card>
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
          </a>


        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
export default {
  name: "UpdateForm",
  data: function () {
    return {
      rooms: [],
      room_id: '',
      nums: '',
    }
  },
  methods: {
    submitUpdate() {
      let url = "/api" + "/room/" + this.room_id + "/audience"
      this.$axios({
        url: url,
        method: "put",
        params: {
          nums: this.nums,
        }
      }).then(res => {
        console.log(res.data)
      }).catch(res => {
        console.log(res.data)
      }).finally(
          this.refreshRooms
      )
    },
    refreshRooms: function () {
      this.$axios({
        url: "/api/room/status",
        method: "get",
      }).then(res => {
        console.log(res)
        this.rooms = res.data
      }).catch(res => {
        console.log(res)
      })
    },
  },
  mounted() {
    this.refreshRooms()
  }
}
</script>

<style scoped>
a {
  text-decoration: none;
}

</style>